// Copyright 2018-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/api/events"
	"github.com/containerd/containerd/api/services/tasks/v1"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/pkg/ttrpcutil"
	"github.com/containerd/containerd/runtime"
	"github.com/containerd/typeurl"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/process"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/firecracker-microvm/firecracker-containerd/firecracker-control"
	fcClient "github.com/firecracker-microvm/firecracker-containerd/firecracker-control/client"
	"github.com/firecracker-microvm/firecracker-containerd/internal"
	"github.com/firecracker-microvm/firecracker-containerd/internal/vm"
	"github.com/firecracker-microvm/firecracker-containerd/proto"
	fccontrol "github.com/firecracker-microvm/firecracker-containerd/proto/service/fccontrol/ttrpc"
	"github.com/firecracker-microvm/firecracker-containerd/runtime/firecrackeroci"
)

const (
	defaultNamespace = namespaces.Default

	containerdSockPath = "/run/containerd/containerd.sock"

	firecrackerRuntime = "aws.firecracker"
	shimProcessName    = "containerd-shim-aws-firecracker"

	defaultVMRootfsPath = "/var/lib/firecracker-containerd/runtime/default-rootfs.img"
	defaultVMNetDevName = "eth0"
	varRunDir           = "/run/firecracker-containerd"
)

// Images are presumed by the isolated tests to have already been pulled
// into the content store. This will just unpack the layers into an
// image with the provided snapshotter.
func unpackImage(ctx context.Context, client *containerd.Client, snapshotterName string, imageRef string) (containerd.Image, error) {
	img, err := client.GetImage(ctx, imageRef)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get image")
	}

	err = img.Unpack(ctx, snapshotterName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack image")
	}

	return img, nil
}

func alpineImage(ctx context.Context, client *containerd.Client, snapshotterName string) (containerd.Image, error) {
	return unpackImage(ctx, client, snapshotterName, "docker.io/library/alpine:3.10.1")
}

func iperf3Image(ctx context.Context, client *containerd.Client, snapshotterName string) (containerd.Image, error) {
	return unpackImage(ctx, client, snapshotterName, "docker.io/mlabbe/iperf3:3.6-r0")
}

func TestShimExitsUponContainerDelete_Isolated(t *testing.T) {
	prepareIntegTest(t)

	ctx := namespaces.WithNamespace(context.Background(), defaultNamespace)

	client, err := containerd.New(containerdSockPath)
	require.NoError(t, err, "unable to create client to containerd service at %s, is containerd running?", containerdSockPath)
	defer client.Close()

	image, err := alpineImage(ctx, client, defaultSnapshotterName())
	require.NoError(t, err, "failed to get alpine image")

	testTimeout := 60 * time.Second
	testCtx, testCancel := context.WithTimeout(ctx, testTimeout)
	defer testCancel()

	containerName := fmt.Sprintf("%s-%d", t.Name(), time.Now().UnixNano())
	snapshotName := fmt.Sprintf("%s-snapshot", containerName)
	container, err := client.NewContainer(testCtx,
		containerName,
		containerd.WithRuntime(firecrackerRuntime, nil),
		containerd.WithSnapshotter(defaultSnapshotterName()),
		containerd.WithNewSnapshot(snapshotName, image),
		containerd.WithNewSpec(
			oci.WithProcessArgs("sleep", fmt.Sprintf("%d", testTimeout/time.Second)),
			oci.WithDefaultPathEnv,
		),
	)
	require.NoError(t, err, "failed to create container %s", containerName)

	_, err = client.NewContainer(testCtx,
		fmt.Sprintf("should-fail-%s-%d", t.Name(), time.Now().UnixNano()),
		containerd.WithRuntime(firecrackerRuntime, nil),
		containerd.WithSnapshotter(defaultSnapshotterName()),
		containerd.WithNewSnapshot(snapshotName, image),
		containerd.WithNewSpec(
			oci.WithProcessArgs("sleep", fmt.Sprintf("%d", testTimeout/time.Second)),
			oci.WithDefaultPathEnv,
		),
	)
	require.Error(t, err, "should not be able to create additional container when no drives are available")

	task, err := container.NewTask(testCtx, cio.NewCreator(cio.WithStdio))
	require.NoError(t, err, "failed to create task for container %s", containerName)

	exitEventCh, exitEventErrCh := client.Subscribe(testCtx, fmt.Sprintf(`topic=="%s"`, runtime.TaskExitEventTopic))

	err = task.Start(testCtx)
	require.NoError(t, err, "failed to start task for container %s", containerName)

	shimProcesses, err := internal.WaitForProcessToExist(testCtx, time.Second,
		func(ctx context.Context, p *process.Process) (bool, error) {
			processExecutable, err := p.ExeWithContext(ctx)
			if err != nil {
				return false, err
			}

			return filepath.Base(processExecutable) == shimProcessName, nil
		},
	)
	require.NoError(t, err, "failed waiting for expected shim process %q to come up", shimProcessName)
	require.Len(t, shimProcesses, 1, "expected only one shim process to exist")
	shimProcess := shimProcesses[0]

	err = task.Kill(testCtx, syscall.SIGKILL)
	require.NoError(t, err, "failed to SIGKILL containerd task %s", containerName)

	_, err = task.Delete(testCtx)
	require.NoError(t, err, "failed to Delete containerd task %s", containerName)

	select {
	case envelope := <-exitEventCh:
		unmarshaledEvent, err := typeurl.UnmarshalAny(envelope.Event)
		require.NoError(t, err, "failed to unmarshal event")

		switch event := unmarshaledEvent.(type) {
		case *events.TaskExit:
			require.Equal(t, container.ID(), event.ContainerID, "received exit event from expected container %s", container.ID())
		default:
			require.Fail(t, "unexpected event type", "received unexpected non-exit event type on topic: %s", envelope.Topic)
		}

		err = internal.WaitForPidToExit(testCtx, time.Second, shimProcess.Pid)
		require.NoError(t, err, "failed waiting for shim process \"%s\" to exit", shimProcessName)

		namespaceVarRunDir := filepath.Join(varRunDir, namespaces.Default)
		varRunFCContents, err := ioutil.ReadDir(namespaceVarRunDir)
		require.NoError(t, err, `failed to list directory "%s"`, namespaceVarRunDir)
		require.Len(t, varRunFCContents, 0, "expect %s to be cleared after shims shutdown", namespaceVarRunDir)
	case err = <-exitEventErrCh:
		require.Fail(t, "unexpected error", "unexpectedly received on task exit error channel: %s", err.Error())
	case <-testCtx.Done():
		require.Fail(t, "context canceled", "context canceled while waiting for container \"%s\" exit: %s", containerName, testCtx.Err())
	}
}

// vmIDtoMacAddr converts a provided VMID to a unique Mac Address. This is a convenient way of providing the VMID to clients within
// the VM without the extra complication of alternative approaches like MMDS.
func vmIDtoMacAddr(vmID uint) string {
	var addrParts []string

	// mac addresses have 6 hex components separate by ":", i.e. "11:22:33:44:55:66"
	numMacAddrComponents := uint(6)

	for n := uint(0); n < numMacAddrComponents; n++ {
		// To isolate the value of the nth component, right bit shift the vmID by 8*n (there are 8 bits per component) and
		// mask out any upper bits leftover (bitwise AND with 255)
		addrComponent := (vmID >> (8 * n)) & 255

		// format the component as a two-digit hex string
		addrParts = append(addrParts, fmt.Sprintf("%02x", addrComponent))
	}

	return strings.Join(addrParts, ":")
}

func createTapDevice(ctx context.Context, tapName string) error {
	err := exec.CommandContext(ctx, "ip", "tuntap", "add", tapName, "mode", "tap").Run()
	if err != nil {
		return errors.Wrapf(err, "failed to create tap device %s", tapName)
	}

	err = exec.CommandContext(ctx, "ip", "link", "set", tapName, "up").Run()
	if err != nil {
		return errors.Wrapf(err, "failed to up tap device %s", tapName)
	}

	return nil
}

func TestMultipleVMs_Isolated(t *testing.T) {
	prepareIntegTest(t, withJailer())

	cases := []struct {
		MaxContainers int32
		JailerConfig  *proto.JailerConfig
	}{
		{
			MaxContainers: 5,
		},
		{
			MaxContainers: 5,
		},
		{
			MaxContainers: 5,
		},
		{
			MaxContainers: 3,
			JailerConfig:  &proto.JailerConfig{},
		},
		{
			MaxContainers: 3,
			JailerConfig:  &proto.JailerConfig{},
		},
	}

	testTimeout := 600 * time.Second
	ctx, cancel := context.WithTimeout(namespaces.WithNamespace(context.Background(), defaultNamespace), testTimeout)
	defer cancel()

	client, err := containerd.New(containerdSockPath, containerd.WithDefaultRuntime(firecrackerRuntime))
	require.NoError(t, err, "unable to create client to containerd service at %s, is containerd running?", containerdSockPath)
	defer client.Close()

	image, err := alpineImage(ctx, client, defaultSnapshotterName())
	require.NoError(t, err, "failed to get alpine image")

	pluginClient, err := ttrpcutil.NewClient(containerdSockPath + ".ttrpc")
	require.NoError(t, err, "failed to create ttrpc client")

	// This test spawns separate VMs in parallel and ensures containers are spawned within each expected VM. It asserts each
	// container ends up in the right VM by assigning each VM a network device with a unique mac address and having each container
	// print the mac address it sees inside its VM.
	var vmWg sync.WaitGroup
	for vmID, c := range cases {
		vmWg.Add(1)
		go func(vmID int, containerCount int32, jailerConfig *proto.JailerConfig) {
			defer vmWg.Done()

			tapName := fmt.Sprintf("tap%d", vmID)
			err = createTapDevice(ctx, tapName)
			require.NoError(t, err, "failed to create tap device for vm %d", vmID)

			rootfsPath := defaultVMRootfsPath

			fcClient := fccontrol.NewFirecrackerClient(pluginClient.Client())
			req := &proto.CreateVMRequest{
				VMID: strconv.Itoa(vmID),
				MachineCfg: &proto.FirecrackerMachineConfiguration{
					MemSizeMib: 512,
				},
				RootDrive: &proto.FirecrackerRootDrive{
					HostPath: rootfsPath,
				},
				NetworkInterfaces: []*proto.FirecrackerNetworkInterface{
					{
						AllowMMDS: true,
						StaticConfig: &proto.StaticNetworkConfiguration{
							HostDevName: tapName,
							MacAddress:  vmIDtoMacAddr(uint(vmID)),
						},
					},
				},
				ContainerCount: containerCount,
				JailerConfig:   jailerConfig,
			}

			if jailerConfig != nil {
				req.NetworkInterfaces = nil
			}

			_, err = fcClient.CreateVM(ctx, req)
			require.NoError(t, err, "failed to create vm")

			var containerWg sync.WaitGroup
			for containerID := 0; containerID < int(containerCount); containerID++ {
				containerWg.Add(1)
				go func(containerID int) {
					defer containerWg.Done()
					containerName := fmt.Sprintf("container-%d-%d", vmID, containerID)
					snapshotName := fmt.Sprintf("snapshot-%d-%d", vmID, containerID)
					processArgs := oci.WithProcessArgs("/bin/sh", "-c", strings.Join([]string{
						fmt.Sprintf("/bin/cat /sys/class/net/%s/address", defaultVMNetDevName),
						"/usr/bin/readlink /proc/self/ns/mnt",
						fmt.Sprintf("/bin/sleep %d", testTimeout/time.Second),
					}, " && "))

					if jailerConfig != nil {
						// TODO: this if statement block can go away once we add netns
						processArgs = oci.WithProcessArgs("/bin/sh", "-c", strings.Join([]string{
							fmt.Sprintf("/bin/sleep %d", testTimeout/time.Second),
						}, " && "))
					}

					// spawn a container that just prints the VM's eth0 mac address (which we have set uniquely per VM)
					newContainer, err := client.NewContainer(ctx,
						containerName,
						containerd.WithSnapshotter(defaultSnapshotterName()),
						containerd.WithNewSnapshot(snapshotName, image),
						containerd.WithNewSpec(
							processArgs,
							oci.WithHostNamespace(specs.NetworkNamespace),
							firecrackeroci.WithVMID(strconv.Itoa(vmID)),
						),
					)
					require.NoError(t, err, "failed to create container %s", containerName)

					var taskStdout bytes.Buffer
					var taskStderr bytes.Buffer

					newTask, err := newContainer.NewTask(ctx,
						cio.NewCreator(cio.WithStreams(nil, &taskStdout, &taskStderr)))
					require.NoError(t, err, "failed to create task for container %s", containerName)

					taskExitCh, err := newTask.Wait(ctx)
					require.NoError(t, err, "failed to wait on task for container %s", containerName)

					err = newTask.Start(ctx)
					require.NoError(t, err, "failed to start task for container %s", containerName)

					// Create a few execs for the task, including one with the same ID as the taskID (to provide
					// regression coverage for a bug related to using the same task and exec ID).
					//
					// Save each of their stdout buffers, which will later be compared to ensure they each have
					// the same output.
					//
					// The output of the exec is the mount namespace in which it found itself executing. This
					// will be compared with the mount namespace the task is executing to ensure they are the same.
					// This is a rudimentary way of asserting that each exec was created in the expected task.
					execIDs := []string{fmt.Sprintf("exec-%d-%d", vmID, containerID), containerName}
					execStdouts := make(chan string, len(execIDs))
					var execWg sync.WaitGroup
					for _, execID := range execIDs {
						execWg.Add(1)
						go func(execID string) {
							defer execWg.Done()
							var execStdout bytes.Buffer
							var execStderr bytes.Buffer

							newExec, err := newTask.Exec(ctx, execID, &specs.Process{
								Args: []string{"/usr/bin/readlink", "/proc/self/ns/mnt"},
								Cwd:  "/",
							}, cio.NewCreator(cio.WithStreams(nil, &execStdout, &execStderr)))
							require.NoError(t, err, "failed to exec %s", execID)

							execExitCh, err := newExec.Wait(ctx)
							require.NoError(t, err, "failed to wait on exec %s", execID)

							err = newExec.Start(ctx)
							require.NoError(t, err, "failed to start exec %s", execID)

							select {
							case exitStatus := <-execExitCh:
								_, err = client.TaskService().DeleteProcess(ctx, &tasks.DeleteProcessRequest{
									ContainerID: containerName,
									ExecID:      execID,
								})
								require.NoError(t, err, "failed to delete exec %q", execID)

								// if there was anything on stderr, print it to assist debugging
								stderrOutput := execStderr.String()
								if len(stderrOutput) != 0 {
									fmt.Printf("stderr output from exec %q: %q", execID, stderrOutput)
								}

								mntNS := strings.TrimSpace(execStdout.String())
								require.NotEmptyf(t, mntNS, "no stdout output for task %q exec %q", containerName, execID)
								execStdouts <- mntNS

								require.Equal(t, uint32(0), exitStatus.ExitCode())
							case <-ctx.Done():
								require.Fail(t, "context cancelled",
									"context cancelled while waiting for exec %s to exit, err: %v", execID, ctx.Err())
							}
						}(execID)
					}
					execWg.Wait()
					close(execStdouts)

					if jailerConfig != nil {
						shimDir, err := vm.ShimDir("default", strconv.Itoa(vmID))
						require.NoError(t, err, "failed to get shim dir")

						_, err = os.Stat(filepath.Join(shimDir.RootPath(), "config.json"))
						require.NoError(t, err, "failed to stat jailer's runc config file")
					}

					// Verify each exec had the same stdout and use that value as the mount namespace that will be compared
					// against that of the task below.
					var execMntNS string
					for execStdout := range execStdouts {
						if execMntNS == "" {
							// This is the first iteration of loop; we do a check that execStdout is not "" via require.NotEmptyf
							// in the execID loop above.
							execMntNS = execStdout
						}

						require.Equal(t, execMntNS, execStdout, "execs in same task unexpectedly have different outputs")
					}

					// Now kill the task and verify it was in the right VM and has the same mnt namespace as its execs
					err = newTask.Kill(ctx, syscall.SIGKILL)
					require.NoError(t, err, "failed to kill task %q", containerName)

					select {
					case <-taskExitCh:
						_, err = client.TaskService().DeleteProcess(ctx, &tasks.DeleteProcessRequest{
							ContainerID: containerName,
						})
						require.NoError(t, err, "failed to delete task %q", containerName)

						// if there was anything on stderr, print it to assist debugging
						stderrOutput := taskStderr.String()
						if len(stderrOutput) != 0 {
							fmt.Printf("stderr output from task %q: %q", containerName, stderrOutput)
						}

						stdoutLines := strings.Split(strings.TrimSpace(taskStdout.String()), "\n")
						lines := 2
						if jailerConfig != nil {
							lines = 1
						}
						require.Len(t, stdoutLines, lines)

						printedVMID := strings.TrimSpace(stdoutLines[0])
						// TODO: Remove this if statement once we can add a netns which
						// will allow firecracker to have visibility of the tap devices.
						if jailerConfig == nil {
							require.Equal(t, vmIDtoMacAddr(uint(vmID)), printedVMID, "unexpected VMID output from container %q", containerName)

							taskMntNS := strings.TrimSpace(stdoutLines[1])
							require.Equal(t, execMntNS, taskMntNS, "unexpected mnt NS output from container %q", containerName)
						}

					case <-ctx.Done():
						require.Fail(t, "context cancelled",
							"context cancelled while waiting for container %s to exit, err: %v", containerName, ctx.Err())
					}
				}(containerID)
			}

			// verify duplicate CreateVM call fails with right error
			_, err = fcClient.CreateVM(ctx, &proto.CreateVMRequest{VMID: strconv.Itoa(vmID)})
			require.Error(t, err, "did not receive expected error for duplicate CreateVM call")

			// verify GetVMInfo returns expected data
			vmInfoResp, err := fcClient.GetVMInfo(ctx, &proto.GetVMInfoRequest{VMID: strconv.Itoa(vmID)})
			require.NoError(t, err, "failed to get VM Info for VM %d", vmID)
			require.Equal(t, vmInfoResp.VMID, strconv.Itoa(vmID))
			require.Equal(t, vmInfoResp.SocketPath, filepath.Join(varRunDir, defaultNamespace, strconv.Itoa(vmID), "firecracker.sock"))
			require.Equal(t, vmInfoResp.LogFifoPath, filepath.Join(varRunDir, defaultNamespace, strconv.Itoa(vmID), "fc-logs.fifo"))
			require.Equal(t, vmInfoResp.MetricsFifoPath, filepath.Join(varRunDir, defaultNamespace, strconv.Itoa(vmID), "fc-metrics.fifo"))

			// just verify that updating the metadata doesn't return an error, a separate test case is needed
			// to very the MMDS update propagates to the container correctly
			_, err = fcClient.SetVMMetadata(ctx, &proto.SetVMMetadataRequest{
				VMID:     strconv.Itoa(vmID),
				Metadata: "{}",
			})
			require.NoError(t, err, "failed to set VM Metadata for VM %d", vmID)

			containerWg.Wait()

			_, err = fcClient.StopVM(ctx, &proto.StopVMRequest{VMID: strconv.Itoa(vmID), TimeoutSeconds: 5})
			require.NoError(t, err, "failed to stop VM %d", vmID)
		}(vmID, c.MaxContainers, c.JailerConfig)
	}

	vmWg.Wait()
}

func TestLongUnixSocketPath_Isolated(t *testing.T) {
	prepareIntegTest(t)

	// Verify that if the absolute path of the Firecracker unix sockets are longer
	// than the max length enforced by the kernel (UNIX_PATH_MAX, usually 108), we
	// don't fail (due to the internal implementation using relative paths).
	// We do this by using the max VMID len (76 chars), which in combination with the
	// default location we store state results in a path like
	// "/run/firecracker-containerd/default/<vmID>" (with len 112).
	const maxUnixSockLen = 108
	vmID := strings.Repeat("x", 76)

	ctx := namespaces.WithNamespace(context.Background(), "default")

	pluginClient, err := ttrpcutil.NewClient(containerdSockPath + ".ttrpc")
	require.NoError(t, err, "failed to create ttrpc client")

	fcClient := fccontrol.NewFirecrackerClient(pluginClient.Client())
	_, err = fcClient.CreateVM(ctx, &proto.CreateVMRequest{
		VMID:              vmID,
		NetworkInterfaces: []*proto.FirecrackerNetworkInterface{},
	})
	require.NoError(t, err, "failed to create VM")

	// double-check that the sockets are at the expected path and that their absolute
	// length exceeds 108 bytes
	shimDir, err := vm.ShimDir("default", vmID)
	require.NoError(t, err, "failed to get shim dir")

	_, err = os.Stat(shimDir.FirecrackerSockPath())
	require.NoError(t, err, "failed to stat firecracker socket path")
	if len(shimDir.FirecrackerSockPath()) <= maxUnixSockLen {
		assert.Failf(t, "firecracker sock absolute path %q is not greater than max unix socket path length", shimDir.FirecrackerSockPath())
	}

	_, err = os.Stat(shimDir.FirecrackerVSockPath())
	require.NoError(t, err, "failed to stat firecracker vsock path")
	if len(shimDir.FirecrackerVSockPath()) <= maxUnixSockLen {
		assert.Failf(t, "firecracker vsock absolute path %q is not greater than max unix socket path length", shimDir.FirecrackerVSockPath())
	}
}

func TestStubBlockDevices_Isolated(t *testing.T) {
	prepareIntegTest(t)

	const vmID = 0

	ctx := namespaces.WithNamespace(context.Background(), "default")

	client, err := containerd.New(containerdSockPath, containerd.WithDefaultRuntime(firecrackerRuntime))
	require.NoError(t, err, "unable to create client to containerd service at %s, is containerd running?", containerdSockPath)
	defer client.Close()

	image, err := alpineImage(ctx, client, defaultSnapshotterName())
	require.NoError(t, err, "failed to get alpine image")

	tapName := fmt.Sprintf("tap%d", vmID)
	err = createTapDevice(ctx, tapName)
	require.NoError(t, err, "failed to create tap device for vm %d", vmID)

	containerName := fmt.Sprintf("%s-%d", t.Name(), time.Now().UnixNano())
	snapshotName := fmt.Sprintf("%s-snapshot", containerName)

	pluginClient, err := ttrpcutil.NewClient(containerdSockPath + ".ttrpc")
	require.NoError(t, err, "failed to create ttrpc client")

	fcClient := fccontrol.NewFirecrackerClient(pluginClient.Client())
	_, err = fcClient.CreateVM(ctx, &proto.CreateVMRequest{
		VMID: strconv.Itoa(vmID),
		NetworkInterfaces: []*proto.FirecrackerNetworkInterface{
			{
				AllowMMDS: true,
				StaticConfig: &proto.StaticNetworkConfiguration{
					HostDevName: tapName,
					MacAddress:  vmIDtoMacAddr(uint(vmID)),
				},
			},
		},
		ContainerCount: 5,
	})
	require.NoError(t, err, "failed to create VM")

	newContainer, err := client.NewContainer(ctx,
		containerName,
		containerd.WithSnapshotter(defaultSnapshotterName()),
		containerd.WithNewSnapshot(snapshotName, image),
		containerd.WithNewSpec(
			firecrackeroci.WithVMID(strconv.Itoa(vmID)),
			oci.WithProcessArgs("/bin/sh", "/var/firecracker-containerd-test/scripts/lsblk.sh"),

			oci.WithMounts([]specs.Mount{
				// Exposes the host kernel's /dev as /dev.
				// By default, runc creates own /dev with a minimal set of pseudo devices such as /dev/null.
				{
					Type:        "bind",
					Options:     []string{"bind"},
					Destination: "/dev",
					Source:      "/dev",
				},

				// Exposes test scripts from the host kernel
				{
					Type:        "bind",
					Options:     []string{"bind"},
					Destination: "/var/firecracker-containerd-test/scripts",
					Source:      "/var/firecracker-containerd-test/scripts",
				},
			}),
			// Make the host kernel's /dev readable
			oci.WithParentCgroupDevices,
		),
	)
	require.NoError(t, err, "failed to create container %s", containerName)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	newTask, err := newContainer.NewTask(ctx,
		cio.NewCreator(cio.WithStreams(nil, &stdout, &stderr)))
	require.NoError(t, err, "failed to create task for container %s", containerName)

	exitCh, err := newTask.Wait(ctx)
	require.NoError(t, err, "failed to wait on task for container %s", containerName)

	err = newTask.Start(ctx)
	require.NoError(t, err, "failed to start task for container %s", containerName)

	const containerID = 0

	select {
	case exitStatus := <-exitCh:
		// if there was anything on stderr, print it to assist debugging
		stderrOutput := stderr.String()
		if len(stderrOutput) != 0 {
			fmt.Printf("stderr output from vm %d, container %d: %s", vmID, containerID, stderrOutput)
		}

		const expectedOutput = `
vdb  254:16   0        0B  0 | 
vdc  254:32   0      512B  0 |  214 244 216 245 215 177 177 177
vdd  254:48   0      512B  0 |  214 244 216 245 215 177 177 177
vde  254:64   0      512B  0 |  214 244 216 245 215 177 177 177
vdf  254:80   0      512B  0 |  214 244 216 245 215 177 177 177`

		parts := strings.Split(stdout.String(), "vdb")
		require.Equal(t, strings.TrimSpace(expectedOutput), strings.TrimSpace("vdb"+parts[1]))
		require.NoError(t, exitStatus.Error(), "failed to retrieve exitStatus")
		require.Equal(t, uint32(0), exitStatus.ExitCode())
	case <-ctx.Done():
		require.Fail(t, "context cancelled",
			"context cancelled while waiting for container %s to exit, err: %v", containerName, ctx.Err())
	}
}

func startAndWaitTask(ctx context.Context, t *testing.T, c containerd.Container) string {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	task, err := c.NewTask(ctx, cio.NewCreator(cio.WithStreams(nil, &stdout, &stderr)))
	require.NoError(t, err, "failed to create task for container %s", c.ID())

	exitCh, err := task.Wait(ctx)
	require.NoError(t, err, "failed to wait on task for container %s", c.ID())

	err = task.Start(ctx)
	require.NoError(t, err, "failed to start task for container %s", c.ID())
	defer func() {
		require.NoError(t, err, "failed to delete task for container %s", c.ID())
	}()

	select {
	case exitStatus := <-exitCh:
		assert.NoError(t, exitStatus.Error(), "failed to retrieve exitStatus")
		assert.Equal(t, uint32(0), exitStatus.ExitCode())

		status, err := task.Delete(ctx)
		assert.NoErrorf(t, err, "failed to delete task %q after exit", c.ID())
		if status != nil {
			assert.NoError(t, status.Error())
		}

		assert.Equal(t, "", stderr.String())
	case <-ctx.Done():
		require.Fail(t, "context cancelled",
			"context cancelled while waiting for container %s to exit, err: %v", c.ID(), ctx.Err())
	}

	return stdout.String()
}

func testCreateContainerWithSameName(t *testing.T, vmID string) {
	ctx := namespaces.WithNamespace(context.Background(), "default")

	pluginClient, err := ttrpcutil.NewClient(containerdSockPath + ".ttrpc")
	require.NoError(t, err, "failed to create ttrpc client")

	// Explicitly specify Container Count = 2 to workaround #230
	if len(vmID) != 0 {
		fcClient := fccontrol.NewFirecrackerClient(pluginClient.Client())
		_, err = fcClient.CreateVM(ctx, &proto.CreateVMRequest{
			VMID:           vmID,
			ContainerCount: 2,
		})
		require.NoError(t, err)
	}

	withNewSpec := containerd.WithNewSpec(oci.WithProcessArgs("echo", "hello"), firecrackeroci.WithVMID(vmID), oci.WithDefaultPathEnv)

	client, err := containerd.New(containerdSockPath, containerd.WithDefaultRuntime(firecrackerRuntime))
	require.NoError(t, err, "unable to create client to containerd service at %s, is containerd running?", containerdSockPath)
	defer client.Close()

	image, err := alpineImage(ctx, client, defaultSnapshotterName())
	require.NoError(t, err, "failed to get alpine image")

	containerName := fmt.Sprintf("%s-%d", t.Name(), time.Now().UnixNano())
	snapshotName := fmt.Sprintf("%s-snapshot", containerName)

	containerPath := fmt.Sprintf("/run/containerd/io.containerd.runtime.v2.task/default/%s", containerName)

	c1, err := client.NewContainer(ctx,
		containerName,
		containerd.WithSnapshotter(defaultSnapshotterName()),
		containerd.WithNewSnapshot(snapshotName, image),
		withNewSpec,
	)
	require.NoError(t, err, "failed to create container %s", containerName)
	require.Equal(t, "hello\n", startAndWaitTask(ctx, t, c1))

	// All resources regarding the container will be deleted
	err = c1.Delete(ctx, containerd.WithSnapshotCleanup)
	require.NoError(t, err, "failed to delete container %s", containerName)

	_, err = os.Stat(containerPath)
	require.True(t, os.IsNotExist(err))

	if len(vmID) != 0 {
		shimPath := fmt.Sprintf("%s/default/%s/%s", varRunDir, vmID, containerName)
		_, err = os.Stat(shimPath)
		require.True(t, os.IsNotExist(err))
	}

	// So, we can launch a new container with the same name
	c2, err := client.NewContainer(ctx,
		containerName,
		containerd.WithSnapshotter(defaultSnapshotterName()),
		containerd.WithNewSnapshot(snapshotName, image),
		withNewSpec,
	)
	require.NoError(t, err, "failed to create container %s", containerName)
	require.Equal(t, "hello\n", startAndWaitTask(ctx, t, c2))

	err = c2.Delete(ctx, containerd.WithSnapshotCleanup)
	require.NoError(t, err, "failed to delete container %s", containerName)

	_, err = os.Stat(containerPath)
	require.True(t, os.IsNotExist(err))

	if len(vmID) != 0 {
		shimPath := fmt.Sprintf("%s/default/%s/%s", varRunDir, vmID, containerName)
		_, err = os.Stat(shimPath)
		require.True(t, os.IsNotExist(err))
	}
}

func TestCreateContainerWithSameName_Isolated(t *testing.T) {
	prepareIntegTest(t)

	testCreateContainerWithSameName(t, "")

	vmID := fmt.Sprintf("same-vm-%d", time.Now().UnixNano())
	testCreateContainerWithSameName(t, vmID)
}

func TestCreateTooManyContainers_Isolated(t *testing.T) {
	prepareIntegTest(t)

	assert := assert.New(t)

	ctx := namespaces.WithNamespace(context.Background(), "default")

	client, err := containerd.New(containerdSockPath, containerd.WithDefaultRuntime(firecrackerRuntime))
	require.NoError(t, err, "unable to create client to containerd service at %s, is containerd running?", containerdSockPath)
	defer client.Close()

	image, err := alpineImage(ctx, client, defaultSnapshotterName())
	require.NoError(t, err, "failed to get alpine image")

	runEchoHello := containerd.WithNewSpec(oci.WithProcessArgs("echo", "-n", "hello"), firecrackeroci.WithVMID("reuse-same-vm"), oci.WithDefaultPathEnv)

	c1, err := client.NewContainer(ctx,
		"c1",
		containerd.WithSnapshotter(defaultSnapshotterName()),
		containerd.WithNewSnapshot("c1", image),
		runEchoHello,
	)
	assert.Equal("hello", startAndWaitTask(ctx, t, c1))
	require.NoError(t, err, "failed to create a container")

	defer func() {
		err = c1.Delete(ctx, containerd.WithSnapshotCleanup)
		require.NoError(t, err, "failed to delete a container")
	}()

	c2, err := client.NewContainer(ctx,
		"c2",
		containerd.WithSnapshotter(defaultSnapshotterName()),
		containerd.WithNewSnapshot("c2", image),
		runEchoHello,
	)
	require.NoError(t, err, "failed to create a container")

	defer func() {
		err := c2.Delete(ctx, containerd.WithSnapshotCleanup)
		require.NoError(t, err, "failed to delete a container")
	}()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// When we reuse a VM explicitly, we cannot start multiple containers unless we pre-allocate stub drives.
	_, err = c2.NewTask(ctx, cio.NewCreator(cio.WithStreams(nil, &stdout, &stderr)))
	assert.Contains(err.Error(), "There are no remaining stub drives to be used")
	require.Error(t, err)
}

func TestDriveMount_Isolated(t *testing.T) {
	prepareIntegTest(t, func(cfg *Config) {
		cfg.JailerConfig.RuncBinaryPath = "/usr/local/bin/runc"
	})

	testTimeout := 120 * time.Second
	ctx, cancel := context.WithTimeout(namespaces.WithNamespace(context.Background(), defaultNamespace), testTimeout)
	defer cancel()

	ctrdClient, err := containerd.New(containerdSockPath, containerd.WithDefaultRuntime(firecrackerRuntime))
	require.NoError(t, err, "unable to create client to containerd service at %s, is containerd running?", containerdSockPath)

	fcClient, err := fcClient.New(containerdSockPath + ".ttrpc")
	require.NoError(t, err, "failed to create fccontrol client")

	image, err := alpineImage(ctx, ctrdClient, defaultSnapshotterName())
	require.NoError(t, err, "failed to get alpine image")

	vmID := "test-drive-mount"

	vmMounts := []struct {
		VMPath         string
		FilesystemType string
		VMMountOptions []string
		ContainerPath  string
		FSImgFile      internal.FSImgFile
	}{
		{
			// /systemmount meant to make sure logic doesn't ban this just because it begins with /sys
			VMPath:         "/systemmount",
			FilesystemType: "ext4",
			VMMountOptions: []string{"rw", "noatime"},
			ContainerPath:  "/foo",
			FSImgFile: internal.FSImgFile{
				Subpath:  "dir/foo",
				Contents: "foo\n",
			},
		},
		{
			VMPath:         "/mnt",
			FilesystemType: "ext3",
			VMMountOptions: []string{"ro", "relatime"},
			ContainerPath:  "/bar",
			FSImgFile: internal.FSImgFile{
				Subpath:  "dir/bar",
				Contents: "bar\n",
			},
		},
	}

	vmDriveMounts := []*proto.FirecrackerDriveMount{}
	ctrBindMounts := []specs.Mount{}
	ctrCatCommands := []string{}
	for _, vmMount := range vmMounts {
		vmDriveMounts = append(vmDriveMounts, &proto.FirecrackerDriveMount{
			HostPath:       internal.CreateFSImg(ctx, t, vmMount.FilesystemType, vmMount.FSImgFile),
			VMPath:         vmMount.VMPath,
			FilesystemType: vmMount.FilesystemType,
			Options:        vmMount.VMMountOptions,
		})

		ctrBindMounts = append(ctrBindMounts, specs.Mount{
			Source:      vmMount.VMPath,
			Destination: vmMount.ContainerPath,
			Options:     []string{"bind"},
		})

		ctrCatCommands = append(ctrCatCommands, fmt.Sprintf("/bin/cat %s",
			filepath.Join(vmMount.ContainerPath, vmMount.FSImgFile.Subpath),
		))
	}

	_, err = fcClient.CreateVM(ctx, &proto.CreateVMRequest{
		VMID:         vmID,
		DriveMounts:  vmDriveMounts,
		JailerConfig: &proto.JailerConfig{},
	})
	require.NoError(t, err, "failed to create vm")

	containerName := fmt.Sprintf("%s-container", vmID)
	snapshotName := fmt.Sprintf("%s-snapshot", vmID)

	newContainer, err := ctrdClient.NewContainer(ctx,
		containerName,
		containerd.WithSnapshotter(defaultSnapshotterName()),
		containerd.WithNewSnapshot(snapshotName, image),
		containerd.WithNewSpec(
			oci.WithProcessArgs("/bin/sh", "-c", strings.Join(append(ctrCatCommands,
				"/bin/cat /proc/mounts",
			), " && ")),
			oci.WithMounts(ctrBindMounts),
			firecrackeroci.WithVMID(vmID),
		),
	)
	require.NoError(t, err, "failed to create container %s", containerName)

	outputLines := strings.Split(startAndWaitTask(ctx, t, newContainer), "\n")
	if len(outputLines) < len(vmMounts) {
		require.Fail(t, "unexpected ctr output, expected at least %d lines: %+v", len(vmMounts), outputLines)
	}

	mountInfos, err := internal.ParseProcMountLines(outputLines[len(vmMounts):]...)
	require.NoError(t, err, "failed to parse /proc/mount")
	// this is n^2, but it's doubtful the number of mounts will reach a point where that matters...
	for _, vmMount := range vmMounts {
		// Make sure that this vmMount's test file was cat'd by a container previously and output the expected
		// file contents. This ensure the filesystem was successfully mounted in the VM and the container.
		assert.Containsf(t, outputLines[:len(vmMounts)], strings.TrimSpace(vmMount.FSImgFile.Contents),
			"did not find expected test file output for vm mount at %q", vmMount.ContainerPath)

		// iterate over /proc/mounts entries, find this vmMount's entry in there and verify it was mounted
		// with the correct options.
		var foundExpectedMount bool
		for _, actualMountInfo := range mountInfos {
			if actualMountInfo.DestPath == vmMount.ContainerPath {
				foundExpectedMount = true
				assert.Equalf(t, vmMount.FilesystemType, actualMountInfo.Type,
					"vm mount at %q did have expected filesystem type", vmMount.ContainerPath)
				for _, vmMountOption := range vmMount.VMMountOptions {
					assert.Containsf(t, actualMountInfo.Options, vmMountOption,
						"vm mount at %q did not have expected option", vmMount.ContainerPath)
				}
				break
			}
		}
		assert.Truef(t, foundExpectedMount, "did not find expected mount at container path %q", vmMount.ContainerPath)
	}
}

func TestDriveMountFails_Isolated(t *testing.T) {
	prepareIntegTest(t)

	testTimeout := 120 * time.Second
	ctx, cancel := context.WithTimeout(namespaces.WithNamespace(context.Background(), defaultNamespace), testTimeout)
	defer cancel()

	fcClient, err := fcClient.New(containerdSockPath + ".ttrpc")
	require.NoError(t, err, "failed to create fccontrol client")

	testImgHostPath := internal.CreateFSImg(ctx, t, "ext4", internal.FSImgFile{
		Subpath:  "idc",
		Contents: "doesn't matter",
	})

	for _, driveMount := range []*proto.FirecrackerDriveMount{
		{
			HostPath:       testImgHostPath,
			VMPath:         "/proc/foo", // invalid due to being under /proc
			FilesystemType: "ext4",
		},
		{
			HostPath:       testImgHostPath,
			VMPath:         "/dev/foo", // invalid due to being under /dev
			FilesystemType: "ext4",
		},
		{
			HostPath:       testImgHostPath,
			VMPath:         "/sys/foo", // invalid due to being under /sys
			FilesystemType: "ext4",
		},
	} {
		_, err = fcClient.CreateVM(ctx, &proto.CreateVMRequest{
			VMID:        "test-drive-mount-fails",
			DriveMounts: []*proto.FirecrackerDriveMount{driveMount},
		})

		// TODO it would be good to check for more specific error types, see #294 for possible improvements:
		// https://github.com/firecracker-microvm/firecracker-containerd/issues/294
		assert.Error(t, err, "unexpectedly succeeded in creating a VM with a drive mount under banned path")
	}
}

func TestUpdateVMMetadata_Isolated(t *testing.T) {
	prepareIntegTest(t)

	testTimeout := 60 * time.Second
	ctx, cancel := context.WithTimeout(namespaces.WithNamespace(context.Background(), defaultNamespace), testTimeout)
	defer cancel()

	client, err := containerd.New(containerdSockPath, containerd.WithDefaultRuntime(firecrackerRuntime))
	require.NoError(t, err, "unable to create client to containerd service at %s, is containerd running?", containerdSockPath)
	defer client.Close()

	pluginClient, err := ttrpcutil.NewClient(containerdSockPath + ".ttrpc")
	require.NoError(t, err, "failed to create ttrpc client")
	fcClient := fccontrol.NewFirecrackerClient(pluginClient.Client())

	cniNetworkName := "fcnet-test"
	err = writeCNIConf("/etc/cni/conf.d/fcnet-test.conflist",
		"tc-redirect-tap", cniNetworkName, "")
	require.NoError(t, err, "failed to write test cni conf")

	_, err = fcClient.CreateVM(ctx, &proto.CreateVMRequest{
		VMID: "1",
		NetworkInterfaces: []*proto.FirecrackerNetworkInterface{{
			AllowMMDS: true,
			CNIConfig: &proto.CNIConfiguration{
				NetworkName:   cniNetworkName,
				InterfaceName: "veth0",
			},
		}},
		ContainerCount: 2,
	})
	require.NoError(t, err)
	metadata := "{\"thing\":\"42\",\"ThreeThing\":\"wow\"}"
	// Update VMM metadata
	_, err = fcClient.SetVMMetadata(ctx, &proto.SetVMMetadataRequest{
		VMID:     "1",
		Metadata: metadata,
	})
	require.NoError(t, err)
	resp, err := fcClient.GetVMMetadata(ctx, &proto.GetVMMetadataRequest{
		VMID: "1",
	})
	require.NoError(t, err)
	expected := "{\"ThreeThing\":\"wow\",\"thing\":\"42\"}"
	assert.Equal(t, expected, resp.Metadata)
	// Update again to ensure patching works
	_, err = fcClient.UpdateVMMetadata(ctx, &proto.UpdateVMMetadataRequest{
		VMID:     "1",
		Metadata: "{\"TwoThing\":\"6*9\",\"thing\":\"45\"}",
	})
	require.NoError(t, err)

	resp, err = fcClient.GetVMMetadata(ctx, &proto.GetVMMetadataRequest{
		VMID: "1",
	})
	require.NoError(t, err)
	expected = "{\"ThreeThing\":\"wow\",\"TwoThing\":\"6*9\",\"thing\":\"45\"}"
	assert.Equal(t, expected, resp.Metadata)

	// Check inside the vm
	image, err := alpineImage(ctx, client, defaultSnapshotterName())
	require.NoError(t, err, "failed to get alpine image")
	containerName := "mmds-test"

	newContainer, err := client.NewContainer(ctx,
		containerName,
		containerd.WithSnapshotter(defaultSnapshotterName()),
		containerd.WithNewSnapshot("mmds-test-all", image),
		containerd.WithNewSpec(
			oci.WithProcessArgs("/usr/bin/wget",
				"-q",      // don't print to stderr unless an error occurs
				"-O", "-", // write to stdout
				"http://169.254.169.254/"),
			firecrackeroci.WithVMID("1"),
			firecrackeroci.WithVMNetwork,
		),
	)
	require.NoError(t, err, "failed to create container %s", containerName)

	stdout := startAndWaitTask(ctx, t, newContainer)
	t.Logf("stdout output from task %q: %s", containerName, stdout)
	assert.Equalf(t, "ThreeThing\nTwoThing\nthing", stdout, "container %q did not emit expected stdout", containerName)
	// check a single entry
	containerName += "-entry"
	newContainer, err = client.NewContainer(ctx,
		containerName,
		containerd.WithSnapshotter(defaultSnapshotterName()),
		containerd.WithNewSnapshot("mmds-test-entry", image),
		containerd.WithNewSpec(
			oci.WithProcessArgs("/usr/bin/wget",
				"-q",      // don't print to stderr unless an error occurs
				"-O", "-", // write to stdout
				"http://169.254.169.254/thing"),
			firecrackeroci.WithVMID("1"),
			firecrackeroci.WithVMNetwork,
		),
	)
	require.NoError(t, err, "failed to create container %s", containerName)
	stdout = startAndWaitTask(ctx, t, newContainer)
	t.Logf("stdout output from task %q: %s", containerName, stdout)
	assert.Equalf(t, "45", stdout, "container %q did not emit expected stdout", containerName)
}
