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

package stubdrive

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"

	"github.com/firecracker-microvm/firecracker-containerd/internal"
	"github.com/firecracker-microvm/firecracker-containerd/internal/convert"
	"github.com/firecracker-microvm/firecracker-containerd/proto"
	drivemount "github.com/firecracker-microvm/firecracker-containerd/proto/service/drivemount/ttrpc"
	"github.com/firecracker-microvm/firecracker-containerd/runtime/jailer"
)

const (
	// fcSectorSize is the sector size of Firecracker drives
	fcSectorSize = 512
)

var (
	// ErrDrivesExhausted occurs when there are no more drives left to use. This
	// can happen by calling PatchStubDrive greater than the number of drives.
	ErrDrivesExhausted = fmt.Errorf("There are no remaining drives to be used")
)

func CreateContainerStubs(
	machineCfg *firecracker.Config,
	jail jailer.Jailer,
	containerCount int,
	logger *logrus.Entry,
) (*StubDriveHandler, error) {
	var containerStubs []*stubDrive
	for i := 0; i < containerCount; i++ {
		isReadOnly := false
		var rateLimiter *proto.FirecrackerRateLimiter
		stubFileName := fmt.Sprintf("ctrstub%d", i)

		// TODO(sipsma) put in subdir?
		stubDrive, err := newStubDrive(
			filepath.Join(jail.JailPath().RootPath(), stubFileName),
			jail, isReadOnly, rateLimiter, logger)

		if err != nil {
			return nil, errors.Wrap(err, "failed to create container stub drive")
		}

		machineCfg.Drives = append(machineCfg.Drives, models.Drive{
			DriveID:      firecracker.String(stubDrive.driveID),
			PathOnHost:   firecracker.String(stubDrive.stubPath),
			IsReadOnly:   firecracker.Bool(isReadOnly),
			RateLimiter:  convert.RateLimiterFromProto(rateLimiter),
			IsRootDevice: firecracker.Bool(false),
		})
		containerStubs = append(containerStubs, stubDrive)
	}

	return &StubDriveHandler{
		freeDrives: containerStubs,
		usedDrives: make(map[string]*stubDrive),
	}, nil
}

type StubDriveHandler struct {
	freeDrives []*stubDrive
	// map of id -> stub drive being used by that task
	usedDrives map[string]*stubDrive
	mu         sync.Mutex
}

func (h *StubDriveHandler) Reserve(
	id string,
	hostPath string,
	vmPath string,
	filesystemType string,
	options []string,
) (MountableStubDrive, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.freeDrives) == 0 {
		return nil, ErrDrivesExhausted
	}

	freeDrive := h.freeDrives[0]
	h.freeDrives = h.freeDrives[1:]
	h.usedDrives[id] = freeDrive

	return freeDrive.withMountConfig(
		hostPath,
		vmPath,
		filesystemType,
		options,
	), nil
}

func CreateDriveMountStubs(
	machineCfg *firecracker.Config,
	jail jailer.Jailer,
	driveMounts []*proto.FirecrackerDriveMount,
	logger *logrus.Entry,
) ([]MountableStubDrive, error) {
	var containerStubs []MountableStubDrive
	for i, driveMount := range driveMounts {
		isReadOnly := driveMount.IsReadOnly
		rateLimiter := driveMount.RateLimiter
		stubFileName := fmt.Sprintf("drivemntstub%d", i)

		// TODO(sipsma) put in subdir?
		stubDrive, err := newStubDrive(
			filepath.Join(jail.JailPath().RootPath(), stubFileName),
			jail, isReadOnly, rateLimiter, logger)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create drive mount stub drive")
		}

		machineCfg.Drives = append(machineCfg.Drives, models.Drive{
			DriveID:      firecracker.String(stubDrive.driveID),
			PathOnHost:   firecracker.String(stubDrive.stubPath),
			IsReadOnly:   firecracker.Bool(isReadOnly),
			RateLimiter:  convert.RateLimiterFromProto(rateLimiter),
			IsRootDevice: firecracker.Bool(false),
		})
		containerStubs = append(containerStubs, stubDrive.withMountConfig(
			driveMount.HostPath,
			driveMount.VMPath,
			driveMount.FilesystemType,
			driveMount.Options))
	}

	return containerStubs, nil
}

type MountableStubDrive interface {
	PatchAndMount(
		requestCtx context.Context,
		machine firecracker.MachineIface,
		driveMounter drivemount.DriveMounterService,
	) error
}

func newStubDrive(
	stubPath string,
	jail jailer.Jailer,
	isReadOnly bool,
	rateLimiter *proto.FirecrackerRateLimiter,
	logger *logrus.Entry,
) (*stubDrive, error) {
	// use the stubPath as the drive ID since it needs to be unique per-stubdrive anyways
	driveID := base64.RawURLEncoding.EncodeToString([]byte(stubPath))

	f, err := os.OpenFile(stubPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.WithError(err).Errorf("unexpected error during %v close", f.Name())
		}
	}()

	stubContent, err := internal.GenerateStubContent(driveID)
	if err != nil {
		return nil, err
	}

	if _, err := f.WriteString(stubContent); err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := info.Size()
	sectorCount := fileSize / fcSectorSize
	driveSize := fcSectorSize * sectorCount

	remainingBytes := fileSize % fcSectorSize
	if remainingBytes != 0 {
		// If there are any residual bytes, this means we've need to fill the
		// appropriate sector size to ensure that the data is visible to
		// Firecracker.
		driveSize += fcSectorSize
	}

	// Firecracker will not show any drives smaller than 512 bytes. In
	// addition, the drive is read in chunks of 512 bytes; if the drive size is
	// not a multiple of 512 bytes, then the remainder will not be visible to
	// Firecracker. So we adjust to the appropriate size based on the residual
	// bytes remaining.
	if err := os.Truncate(stubPath, driveSize); err != nil {
		return nil, err
	}

	for _, opt := range jail.StubDrivesOptions() {
		err := opt(f)
		if err != nil {
			return nil, err
		}
	}

	return &stubDrive{
		stubPath: stubPath,
		jail:     jail,
		driveID:  driveID,
		driveMount: &proto.FirecrackerDriveMount{
			IsReadOnly:  isReadOnly,
			RateLimiter: rateLimiter,
		},
	}, nil
}

type stubDrive struct {
	stubPath   string
	jail       jailer.Jailer
	driveID    string
	driveMount *proto.FirecrackerDriveMount
}

func (sd stubDrive) withMountConfig(
	hostPath string,
	vmPath string,
	filesystemType string,
	options []string,
) stubDrive {
	sd.driveMount = &proto.FirecrackerDriveMount{
		HostPath:       hostPath,
		VMPath:         vmPath,
		FilesystemType: filesystemType,
		Options:        options,
		IsReadOnly:     sd.driveMount.IsReadOnly,
		RateLimiter:    sd.driveMount.RateLimiter,
	}
	return sd
}

func (sd stubDrive) PatchAndMount(
	requestCtx context.Context,
	machine firecracker.MachineIface,
	driveMounter drivemount.DriveMounterService,
) error {
	// TODO wrap errors
	err := sd.jail.ExposeFileToJail(sd.driveMount.HostPath)
	if err != nil {
		return err
	}

	err = machine.UpdateGuestDrive(requestCtx, sd.driveID, sd.driveMount.HostPath)
	if err != nil {
		return err
	}

	_, err = driveMounter.MountDrive(requestCtx, &drivemount.MountDriveRequest{
		DriveID:         sd.driveID,
		DestinationPath: sd.driveMount.VMPath,
		FilesytemType:   sd.driveMount.FilesystemType,
		Options:         sd.driveMount.Options,
	})
	if err != nil {
		return err
	}

	return nil
}
