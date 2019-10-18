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
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/firecracker-microvm/firecracker-containerd/proto"
	"github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/stretchr/testify/assert"
)

const (
	memSize   = 4096
	vcpuCount = 2
)

func TestLoadConfigDefaults(t *testing.T) {
	configContent := `{}`
	configFile, cleanup := createTempConfig(t, configContent)
	defer cleanup()
	cfg, err := LoadConfig(configFile)
	assert.NoError(t, err, "failed to load config")

	assert.Equal(t, defaultKernelArgs, cfg.KernelArgs, "expected default kernel args")
	assert.Equal(t, defaultKernelPath, cfg.KernelImagePath, "expected default kernel path")
	assert.Equal(t, defaultRootfsPath, cfg.RootDrive, "expected default rootfs path")
	assert.Equal(t, defaultCPUCount, cfg.CPUCount, "expected default CPU count")
	assert.Equal(t, string(defaultCPUTemplate), cfg.CPUTemplate, "expected default CPU template")
}

func TestLoadConfigOverrides(t *testing.T) {
	overrideKernelArgs := "OVERRIDE KERNEL ARGS"
	overrideKernelPath := "OVERRIDE KERNEL PATH"
	overrideRootfsPath := "OVERRIDE ROOTFS PATH"
	overrideCPUCount := 42
	overrideCPUTemplate := "OVERRIDE CPU TEMPLATE"
	configContent := fmt.Sprintf(
		`{
			"kernel_args":"%s",
			"kernel_image_path":"%s",
			"root_drive":"%s",
			"cpu_count": %d,
            "cpu_template": "%s"
		}`, overrideKernelArgs, overrideKernelPath, overrideRootfsPath, overrideCPUCount, overrideCPUTemplate)
	configFile, cleanup := createTempConfig(t, configContent)
	defer cleanup()
	cfg, err := LoadConfig(configFile)
	assert.NoError(t, err, "failed to load config")

	assert.Equal(t, overrideKernelArgs, cfg.KernelArgs, "expected overridden kernel args")
	assert.Equal(t, overrideKernelPath, cfg.KernelImagePath, "expected overridden kernel path")
	assert.Equal(t, overrideRootfsPath, cfg.RootDrive, "expected overridden rootfs path")
	assert.Equal(t, overrideCPUCount, cfg.CPUCount, "expected overridden CPU count")
	assert.Equal(t, overrideCPUTemplate, cfg.CPUTemplate, "expected overridden CPU template")
}

func createTempConfig(t *testing.T, contents string) (string, func()) {
	t.Helper()
	configFile, err := ioutil.TempFile("", "config")
	if err != nil {
		t.Fatal(err, "failed to create temp config file")
	}
	err = ioutil.WriteFile(configFile.Name(), []byte(contents), 0644)
	if err != nil {
		os.Remove(configFile.Name())
		t.Fatal(err, "failed to write contents to temp config file")
	}
	return configFile.Name(), func() { os.Remove(configFile.Name()) }
}

func TestMachineConfigurationFromProto(t *testing.T) {
	testcases := []struct {
		name                  string
		config                *Config
		proto                 *proto.FirecrackerMachineConfiguration
		expectedMachineConfig models.MachineConfiguration
	}{
		{
			name:   "ProtoOnly",
			config: &Config{},
			proto: &proto.FirecrackerMachineConfiguration{
				CPUTemplate: string(models.CPUTemplateC3),
				VcpuCount:   vcpuCount,
				MemSizeMib:  memSize,
				HtEnabled:   true,
			},
			expectedMachineConfig: models.MachineConfiguration{
				CPUTemplate: models.CPUTemplateC3,
				VcpuCount:   firecracker.Int64(vcpuCount),
				MemSizeMib:  firecracker.Int64(memSize),
				HtEnabled:   firecracker.Bool(true),
			},
		},
		{
			name: "ConfigOnly",
			config: &Config{
				CPUTemplate: "C3",
				CPUCount:    vcpuCount,
			},
			proto: &proto.FirecrackerMachineConfiguration{},
			expectedMachineConfig: models.MachineConfiguration{
				CPUTemplate: models.CPUTemplateC3,
				VcpuCount:   firecracker.Int64(vcpuCount),
				MemSizeMib:  firecracker.Int64(defaultMemSizeMb),
				HtEnabled:   firecracker.Bool(false),
			},
		},
		{
			name: "NilProto",
			config: &Config{
				CPUTemplate: "C3",
				CPUCount:    vcpuCount,
			},
			expectedMachineConfig: models.MachineConfiguration{
				CPUTemplate: models.CPUTemplateC3,
				VcpuCount:   firecracker.Int64(vcpuCount),
				MemSizeMib:  firecracker.Int64(defaultMemSizeMb),
				HtEnabled:   firecracker.Bool(false),
			},
		},
		{
			name: "Overrides",
			config: &Config{
				CPUTemplate: "T2",
				CPUCount:    vcpuCount + 1,
			},
			proto: &proto.FirecrackerMachineConfiguration{
				CPUTemplate: string(models.CPUTemplateC3),
				VcpuCount:   vcpuCount,
				MemSizeMib:  memSize,
				HtEnabled:   true,
			},
			expectedMachineConfig: models.MachineConfiguration{
				CPUTemplate: models.CPUTemplateC3,
				VcpuCount:   firecracker.Int64(vcpuCount),
				MemSizeMib:  firecracker.Int64(memSize),
				HtEnabled:   firecracker.Bool(true),
			},
		},
	}

	for _, tc := range testcases {
		tc := tc // see https://github.com/kyoh86/scopelint/issues/4
		t.Run(tc.name, func(t *testing.T) {
			machineConfig := machineConfigurationFromProto(tc.config, tc.proto)
			assert.Equal(t, tc.expectedMachineConfig, machineConfig)
		})
	}
}
