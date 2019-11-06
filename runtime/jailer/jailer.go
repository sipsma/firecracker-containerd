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

package jailer

import (
	"context"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/sirupsen/logrus"

	"github.com/firecracker-microvm/firecracker-containerd/internal/vm"
	"github.com/firecracker-microvm/firecracker-containerd/proto"
	"github.com/firecracker-microvm/firecracker-containerd/runtime/stubdrive"
)

const (
	kernelImageFileName   = "kernel-image"
	jailerHandlerName     = "firecracker-containerd-jail-handler"
	jailerFifoHandlerName = "firecracker-containerd-jail-fifo-handler"
	rootfsFolder          = "rootfs"

	// TODO eventually we can get rid of this when we add usernamespaces to
	// jailing.
	jailerUID = 300000
	jailerGID = 300000
)

var (
	runcConfigPath = "/etc/containerd/firecracker-runc-config.json"
)

// Jailer will allow modification and provide options to the the Firecracker VM
// to allow for jailing. In addition, this will allow for given files to be exposed
// to the jailed filesystem.
type Jailer interface {
	// BuildJailedMachine will modify the firecracker.Config and provide
	// firecracker.Opt to be passed into firecracker.NewMachine which will allow
	// for the VM to be jailed.
	BuildJailedMachine(machineCfg *firecracker.Config, firecrackerBinaryPath string) ([]firecracker.Opt, error)
	// ExposeFileToJail will expose the given file to the jailed filesystem, including
	// regular files and block devices. An error is returned if provided a path to a file
	// with type that is not supported.
	ExposeFileToJail(path string) error
	// JailPath is used to return the directory we are supposed to be working in.
	JailPath() vm.Dir
	// StubDrivesOptions will return a set of options used to create a new stub
	// drive handler.
	StubDrivesOptions() []stubdrive.StubDrivesOpt
}

// NewJailer is used to construct a Jailer from the CreateVM request. If no
// request or jailer config was provided, then the noopJailer will be returned.
func NewJailer(
	ctx context.Context,
	logger *logrus.Entry,
	shimDir vm.Dir,
	runcBinaryPath string,
	request *proto.CreateVMRequest,
) (Jailer, error) {
	if request == nil || request.JailerConfig == nil {
		l := logger.WithField("jailer", "noop")
		return NewNoopJailer(ctx, l, shimDir), nil
	}

	l := logger.WithField("jailer", "runc")
	return NewRuncJailer(ctx, l, shimDir.RootPath(), runcBinaryPath, jailerUID, jailerGID)
}
