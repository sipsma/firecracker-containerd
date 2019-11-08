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
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/firecracker-microvm/firecracker-containerd/internal/vm"
	"github.com/firecracker-microvm/firecracker-containerd/proto"
	"github.com/firecracker-microvm/firecracker-containerd/runtime/jailer"
	"github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mac         = "AA:FC:00:00:00:01"
	hostDevName = "tap0"
)

func TestBuildVMConfiguration(t *testing.T) {
	namespace := "TestBuildVMConfiguration"
	testcases := []struct {
		name                   string
		request                *proto.CreateVMRequest
		config                 *Config
		expectedCfg            *firecracker.Config
		expectedStubDriveCount int
	}{
		{
			name:    "Only Config",
			request: &proto.CreateVMRequest{},
			config: &Config{
				KernelArgs:      "KERNEL ARGS",
				KernelImagePath: "KERNEL IMAGE",
				RootDrive:       "ROOT DRIVE",
				CPUTemplate:     "C3",
				CPUCount:        2,
			},
			expectedCfg: &firecracker.Config{
				KernelArgs:      "KERNEL ARGS",
				KernelImagePath: "KERNEL IMAGE",
				Drives: []models.Drive{
					{
						DriveID:      firecracker.String("root_drive"),
						PathOnHost:   firecracker.String("ROOT DRIVE"),
						IsReadOnly:   firecracker.Bool(true),
						IsRootDevice: firecracker.Bool(true),
					},
				},
				MachineCfg: models.MachineConfiguration{
					CPUTemplate: models.CPUTemplateC3,
					VcpuCount:   firecracker.Int64(2),
					MemSizeMib:  firecracker.Int64(defaultMemSizeMb),
					HtEnabled:   firecracker.Bool(false),
				},
			},
			expectedStubDriveCount: 1,
		},
		{
			name: "Only Request",
			request: &proto.CreateVMRequest{
				KernelArgs:      "REQUEST KERNEL ARGS",
				KernelImagePath: "REQUEST KERNEL IMAGE",
				RootDrive: &proto.FirecrackerRootDrive{
					HostPath:   "REQUEST ROOT DRIVE",
					IsWritable: true,
				},
				MachineCfg: &proto.FirecrackerMachineConfiguration{
					CPUTemplate: "C3",
					VcpuCount:   2,
				},
			},
			config: &Config{},
			expectedCfg: &firecracker.Config{
				KernelArgs:      "REQUEST KERNEL ARGS",
				KernelImagePath: "REQUEST KERNEL IMAGE",
				Drives: []models.Drive{
					{
						DriveID:      firecracker.String("root_drive"),
						PathOnHost:   firecracker.String("REQUEST ROOT DRIVE"),
						IsReadOnly:   firecracker.Bool(false),
						IsRootDevice: firecracker.Bool(true),
					},
				},
				MachineCfg: models.MachineConfiguration{
					CPUTemplate: models.CPUTemplateC3,
					VcpuCount:   firecracker.Int64(2),
					MemSizeMib:  firecracker.Int64(defaultMemSizeMb),
					HtEnabled:   firecracker.Bool(false),
				},
			},
			expectedStubDriveCount: 1,
		},
		{
			name: "Request is prioritized over Config",
			request: &proto.CreateVMRequest{
				KernelArgs:      "REQUEST KERNEL ARGS",
				KernelImagePath: "REQUEST KERNEL IMAGE",
				RootDrive: &proto.FirecrackerRootDrive{
					HostPath:   "REQUEST ROOT DRIVE",
					IsWritable: true,
				},
				MachineCfg: &proto.FirecrackerMachineConfiguration{
					CPUTemplate: "T2",
					VcpuCount:   3,
				},
			},
			config: &Config{
				KernelArgs:      "KERNEL ARGS",
				KernelImagePath: "KERNEL IMAGE",
				CPUTemplate:     "C3",
				CPUCount:        2,
			},
			expectedCfg: &firecracker.Config{
				KernelArgs:      "REQUEST KERNEL ARGS",
				KernelImagePath: "REQUEST KERNEL IMAGE",
				Drives: []models.Drive{
					{
						DriveID:      firecracker.String("root_drive"),
						PathOnHost:   firecracker.String("REQUEST ROOT DRIVE"),
						IsReadOnly:   firecracker.Bool(false),
						IsRootDevice: firecracker.Bool(true),
					},
				},
				MachineCfg: models.MachineConfiguration{
					CPUTemplate: models.CPUTemplateT2,
					VcpuCount:   firecracker.Int64(3),
					MemSizeMib:  firecracker.Int64(defaultMemSizeMb),
					HtEnabled:   firecracker.Bool(false),
				},
			},
			expectedStubDriveCount: 1,
		},
		{
			name: "Request can omit some fields",
			request: &proto.CreateVMRequest{
				KernelArgs:      "REQUEST KERNEL ARGS",
				KernelImagePath: "REQUEST KERNEL IMAGE",
				RootDrive: &proto.FirecrackerRootDrive{
					HostPath: "REQUEST ROOT DRIVE",
				},
				MachineCfg: &proto.FirecrackerMachineConfiguration{},
			},
			config: &Config{
				KernelArgs:      "KERNEL ARGS",
				KernelImagePath: "KERNEL IMAGE",
				CPUTemplate:     "C3",
				CPUCount:        2,
			},
			expectedCfg: &firecracker.Config{
				KernelArgs:      "REQUEST KERNEL ARGS",
				KernelImagePath: "REQUEST KERNEL IMAGE",
				Drives: []models.Drive{
					{
						DriveID:      firecracker.String("root_drive"),
						PathOnHost:   firecracker.String("REQUEST ROOT DRIVE"),
						IsReadOnly:   firecracker.Bool(true),
						IsRootDevice: firecracker.Bool(true),
					},
				},
				MachineCfg: models.MachineConfiguration{
					CPUTemplate: models.CPUTemplateC3,
					VcpuCount:   firecracker.Int64(2),
					MemSizeMib:  firecracker.Int64(defaultMemSizeMb),
					HtEnabled:   firecracker.Bool(false),
				},
			},
			expectedStubDriveCount: 1,
		},
		{
			name:    "Container Count affects StubDriveCount",
			request: &proto.CreateVMRequest{ContainerCount: 2},
			config: &Config{
				KernelArgs:      "KERNEL ARGS",
				KernelImagePath: "KERNEL IMAGE",
				RootDrive:       "ROOT DRIVE",
				CPUTemplate:     "C3",
				CPUCount:        2,
			},
			expectedCfg: &firecracker.Config{
				KernelArgs:      "KERNEL ARGS",
				KernelImagePath: "KERNEL IMAGE",
				Drives: []models.Drive{
					{
						DriveID:      firecracker.String("root_drive"),
						PathOnHost:   firecracker.String("ROOT DRIVE"),
						IsReadOnly:   firecracker.Bool(true),
						IsRootDevice: firecracker.Bool(true),
					},
				},
				MachineCfg: models.MachineConfiguration{
					CPUTemplate: models.CPUTemplateC3,
					VcpuCount:   firecracker.Int64(2),
					MemSizeMib:  firecracker.Int64(defaultMemSizeMb),
					HtEnabled:   firecracker.Bool(false),
				},
			},
			expectedStubDriveCount: 2,
		},
	}

	for _, tc := range testcases {
		tc := tc // see https://github.com/kyoh86/scopelint/issues/4
		t.Run(tc.name, func(t *testing.T) {
			svc := &service{
				namespace: namespace,
				logger:    logrus.WithField("test", namespace+"/"+tc.name),
				config:    tc.config,
			}

			tempDir, err := ioutil.TempDir(os.TempDir(), namespace)
			assert.NoError(t, err)
			defer os.RemoveAll(tempDir)

			svc.shimDir = vm.Dir(tempDir)
			svc.jailer = jailer.NewNoopJailer(context.Background(), svc.logger, svc.shimDir)

			relSockPath, err := svc.shimDir.FirecrackerSockRelPath()
			require.NoError(t, err, "failed to get firecracker sock rel path")

			relVSockPath, err := svc.shimDir.FirecrackerVSockRelPath()
			require.NoError(t, err, "failed to get firecracker vsock rel path")

			// For values that remain constant between tests, they are written here
			tc.expectedCfg.SocketPath = relSockPath
			tc.expectedCfg.VsockDevices = []firecracker.VsockDevice{{
				Path: relVSockPath,
				ID:   "agent_api",
			}}
			tc.expectedCfg.LogFifo = svc.shimDir.FirecrackerLogFifoPath()
			tc.expectedCfg.MetricsFifo = svc.shimDir.FirecrackerMetricsFifoPath()

			drives := make([]models.Drive, tc.expectedStubDriveCount)
			for i := 0; i < tc.expectedStubDriveCount; i++ {
				drives[i].PathOnHost = firecracker.String(filepath.Join(tempDir, fmt.Sprintf("stub%d", i)))
				drives[i].DriveID = firecracker.String(fmt.Sprintf("stub%d", i))
				drives[i].IsReadOnly = firecracker.Bool(false)
				drives[i].IsRootDevice = firecracker.Bool(false)
			}
			tc.expectedCfg.Drives = append(drives, tc.expectedCfg.Drives...)

			actualCfg, err := svc.buildVMConfiguration(tc.request)
			assert.NoError(t, err)
			require.Equal(t, tc.expectedCfg, actualCfg)

			// TODO require.Equal(t, tc.expectedStubDriveCount, len(svc.stubDriveHandler.drives), "The stub driver only knows stub drives")
		})
	}
}

func TestDebugConfig(t *testing.T) {
	cases := []struct {
		name    string
		service *service
	}{
		{
			name: "empty",
			service: &service{
				logger: logrus.NewEntry(logrus.New()),
				config: &Config{},
			},
		},
		{
			name: "LogLevel set",
			service: &service{
				logger: logrus.NewEntry(logrus.New()),
				config: &Config{
					LogLevel: "foo",
				},
			},
		},
		{
			name: "Debug set",
			service: &service{
				logger: logrus.NewEntry(logrus.New()),
				config: &Config{
					Debug: true,
				},
			},
		},
		{
			name: "Both set",
			service: &service{
				logger: logrus.NewEntry(logrus.New()),
				config: &Config{
					LogLevel: "foo",
					Debug:    true,
				},
			},
		},
	}

	cwd, err := os.Getwd()
	require.NoError(t, err, "failed to get working dir")

	path, err := ioutil.TempDir(cwd, "TestDebugConfig")
	assert.NoError(t, err, "failed to create temp directory")

	defer os.RemoveAll(path)

	for i, c := range cases {
		c := c
		stubDrivePath := filepath.Join(path, fmt.Sprintf("%d", i))
		err := os.MkdirAll(stubDrivePath, os.ModePerm)
		assert.NoError(t, err, "failed to create stub drive path")

		c.service.shimDir = vm.Dir(stubDrivePath)
		c.service.jailer = jailer.NewNoopJailer(context.Background(), c.service.logger, c.service.shimDir)

		req := proto.CreateVMRequest{}

		cfg, err := c.service.buildVMConfiguration(&req)
		assert.NoError(t, err, "failed to build firecracker configuration")
		assert.Equal(t, c.service.config.LogLevel, cfg.LogLevel)
		assert.Equal(t, c.service.config.Debug, cfg.Debug)
	}
}
