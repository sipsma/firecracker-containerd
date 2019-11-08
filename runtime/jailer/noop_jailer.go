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

	"github.com/firecracker-microvm/firecracker-containerd/internal/vm"
	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/sirupsen/logrus"
)

// noopJailer is a jailer that returns only successful responses and performs
// no operations during calls
type noopJailer struct {
	logger  *logrus.Entry
	shimDir vm.Dir
	ctx     context.Context
}

func NewNoopJailer(ctx context.Context, logger *logrus.Entry, shimDir vm.Dir) Jailer {
	return noopJailer{
		logger:  logger,
		shimDir: shimDir,
		ctx:     ctx,
	}
}

func (j noopJailer) BuildJailedMachine(machineConfig *firecracker.Config, firecrackerBinaryPath string) ([]firecracker.Opt, error) {
	if len(firecrackerBinaryPath) == 0 {
		return []firecracker.Opt{}, nil
	}

	relSocketPath, err := j.shimDir.FirecrackerSockRelPath()
	if err != nil {
		return nil, err
	}

	cmd := firecracker.VMCommandBuilder{}.
		WithBin(firecrackerBinaryPath).
		WithSocketPath(relSocketPath).
		Build(j.ctx)

	if machineConfig.Debug {
		cmd.Stdout = j.logger.WithField("vmm_stream", "stdout").WriterLevel(logrus.DebugLevel)
		cmd.Stderr = j.logger.WithField("vmm_stream", "stderr").WriterLevel(logrus.DebugLevel)
	}

	j.logger.Debug("noop operation for BuildJailedMachine")
	return []firecracker.Opt{
		firecracker.WithProcessRunner(cmd),
	}, nil
}

func (j noopJailer) JailPath() vm.Dir {
	j.logger.Debug("noop operation returning shim dir for JailPath")
	return j.shimDir
}

func (j noopJailer) ExposeFileToJail(path string) error {
	j.logger.Debug("noop operation for ExposeFileToJail")
	return nil
}

func (j noopJailer) StubDrivesOptions() []FileOpt {
	j.logger.Debug("noop operation for StubDrivesOptions")
	return []FileOpt{}
}
