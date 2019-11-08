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
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/firecracker-microvm/firecracker-containerd/proto"
	"github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

const (
	configPathEnvName  = "FIRECRACKER_CONTAINERD_RUNTIME_CONFIG_PATH"
	defaultConfigPath  = "/etc/containerd/firecracker-runtime.json"
	defaultKernelArgs  = "console=ttyS0 noapic reboot=k panic=1 pci=off nomodules rw"
	defaultFilesPath   = "/var/lib/firecracker-containerd/runtime/"
	defaultKernelPath  = defaultFilesPath + "default-vmlinux.bin"
	defaultRootfsPath  = defaultFilesPath + "default-rootfs.img"
	defaultCPUCount    = 1
	defaultCPUTemplate = models.CPUTemplateT2
	defaultMemSizeMb   = 128
)

// Config represents runtime configuration parameters
type Config struct {
	FirecrackerBinaryPath string `json:"firecracker_binary_path"`
	KernelImagePath       string `json:"kernel_image_path"`
	KernelArgs            string `json:"kernel_args"`
	RootDrive             string `json:"root_drive"`
	CPUCount              int    `json:"cpu_count"`
	CPUTemplate           string `json:"cpu_template"`
	LogLevel              string `json:"log_level"`
	HtEnabled             bool   `json:"ht_enabled"`
	Debug                 bool   `json:"debug"`
	// If a CreateVM call specifies no network interfaces and DefaultNetworkInterfaces is non-empty,
	// the VM will default to using the network interfaces as specified here. This is especially
	// useful when a CNI-based network interface is provided in DefaultNetworkInterfaces.
	DefaultNetworkInterfaces []proto.FirecrackerNetworkInterface `json:"default_network_interfaces"`
	JailerConfig             JailerConfig                        `json:"jailer"`
}

// JailerConfig houses a set of configurable values for jailing
// TODO: Add netns field
type JailerConfig struct {
	RuncBinaryPath string `json:"runc_binary_path"`
}

// LoadConfig loads configuration from JSON file at 'path'
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = os.Getenv(configPathEnvName)
	}

	if path == "" {
		path = defaultConfigPath
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read config from %q", path)
	}

	cfg := &Config{
		KernelArgs:      defaultKernelArgs,
		KernelImagePath: defaultKernelPath,
		RootDrive:       defaultRootfsPath,
		CPUCount:        defaultCPUCount,
		CPUTemplate:     string(defaultCPUTemplate),
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal config from %q", path)
	}

	return cfg, nil
}

func machineConfigurationFromProto(cfg *Config, req *proto.FirecrackerMachineConfiguration) models.MachineConfiguration {
	config := models.MachineConfiguration{
		CPUTemplate: models.CPUTemplate(cfg.CPUTemplate),
		VcpuCount:   firecracker.Int64(int64(cfg.CPUCount)),
		MemSizeMib:  firecracker.Int64(defaultMemSizeMb),
		HtEnabled:   firecracker.Bool(cfg.HtEnabled),
	}

	if req == nil {
		return config
	}

	if name := req.CPUTemplate; name != "" {
		config.CPUTemplate = models.CPUTemplate(name)
	}

	if count := req.VcpuCount; count > 0 {
		config.VcpuCount = firecracker.Int64(int64(count))
	}

	if size := req.MemSizeMib; size > 0 {
		config.MemSizeMib = firecracker.Int64(int64(size))
	}

	config.HtEnabled = firecracker.Bool(req.HtEnabled)

	return config
}
