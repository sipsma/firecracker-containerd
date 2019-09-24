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
	"log"
	"os"
	"text/template"
)

const installDependencies = `cd ~

# Install git, Go 1.11, make, curl
sudo mkdir -p /etc/apt/sources.list.d
echo "deb http://ftp.debian.org/debian stretch-backports main" | \
     sudo tee /etc/apt/sources.list.d/stretch-backports.list
sudo DEBIAN_FRONTEND=noninteractive apt-get update
sudo DEBIAN_FRONTEND=noninteractive apt-get \
  --target-release stretch-backports \
  install --yes \
  golang-go \
  gnupg2 \
  lsb-release \
  apt-transport-https \
  make \
  git \
  curl \
  e2fsprogs \
  dmsetup \
  util-linux

cd ~

# Install Docker CE
# Docker CE includes containerd, but we need a separate containerd binary, built
# in a later step
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
apt-key finger docker@docker.com | grep '9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88' || echo '**Cannot find Docker key**'
echo "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | \
     sudo tee /etc/apt/sources.list.d/docker.list
sudo DEBIAN_FRONTEND=noninteractive apt-get update
sudo DEBIAN_FRONTEND=noninteractive apt-get \
     install --yes \
     docker-ce aufs-tools-
sudo usermod -aG docker $(whoami)

cd ~

# Download kernel
curl -fsSL -o hello-vmlinux.bin https://s3.amazonaws.com/spec.ccfc.min/img/hello/kernel/hello-vmlinux.bin`

const checkoutFCCD = `# Check out firecracker-containerd
cd ~
git clone https://github.com/firecracker-microvm/firecracker-containerd.git`

const setupFCCD = `# Build firecracker-containerd.  This includes:
# * block-device snapshotter gRPC proxy plugins
# * firecracker-containerd runtime, a containerd v2 runtime
# * firecracker-containerd agent, an inside-VM component
# * runc, to run containers inside the VM
# * a Debian-based root filesystem configured as read-only with a read-write
#   overlay
# * firecracker-containerd, an alternative containerd binary that includes the
#   firecracker VM lifecycle plugin and API
cd ~/firecracker-containerd
sg docker -c 'make all image firecracker'
sudo make install install-firecracker demo-network

# Configure our firecracker-containerd binary to use our new snapshotter and
# separate storage from the default containerd binary
sudo mkdir -p /etc/firecracker-containerd
sudo mkdir -p /var/lib/firecracker-containerd/containerd
sudo mkdir -p /run/firecracker-containerd
sudo tee /etc/firecracker-containerd/config.toml <<EOF
disabled_plugins = ["cri"]
root = "/var/lib/firecracker-containerd/containerd"
state = "/run/firecracker-containerd"
[grpc]
  address = "/run/firecracker-containerd/containerd.sock"
[proxy_plugins]
  [proxy_plugins.firecracker-naive]
    type = "snapshot"
    address = "/var/run/firecracker-containerd/naive-snapshotter.sock"

[debug]
  level = "debug"
EOF

cd ~

# Configure the aws.firecracker runtime
# The long kernel command-line configures systemd inside the Debian-based image
# and uses a special init process to create a read-write overlay on top of the
# read-only image.
sudo mkdir -p /var/lib/firecracker-containerd/runtime
sudo cp ~/firecracker-containerd/tools/image-builder/rootfs.img /var/lib/firecracker-containerd/runtime/default-rootfs.img
sudo cp ~/hello-vmlinux.bin /var/lib/firecracker-containerd/runtime/default-vmlinux.bin
sudo mkdir -p /etc/containerd
sudo tee /etc/containerd/firecracker-runtime.json <<EOF
{
  "firecracker_binary_path": "/usr/local/bin/firecracker",
  "cpu_template": "T2",
  "log_fifo": "fc-logs.fifo",
  "log_level": "Debug",
  "metrics_fifo": "fc-metrics.fifo",
  "kernel_args": "console=ttyS0 noapic reboot=k panic=1 pci=off nomodules ro systemd.journald.forward_to_console systemd.unit=firecracker.target init=/sbin/overlay-init",
  "default_network_interfaces": [{
    "CNIConfig": {
      "NetworkName": "fcnet",
      "InterfaceName": "veth0"
    }
  }]
}
EOF`

const startSnapshotter = `sudo mkdir -p /var/run/firecracker-containerd /var/lib/firecracker-containerd/naive
sudo naive_snapshotter \
     -address /var/run/firecracker-containerd/naive-snapshotter.sock \
     -path /var/lib/firecracker-containerd/naive \
     -debug`

const startContainerd = `sudo firecracker-containerd --config /etc/firecracker-containerd/config.toml`

const pullImage = `sudo firecracker-ctr --address /run/firecracker-containerd/containerd.sock \
     image pull \
     --snapshotter firecracker-naive \
     docker.io/library/debian:latest`

const startContainer = `sudo firecracker-ctr --address /run/firecracker-containerd/containerd.sock \
     run \
     --snapshotter firecracker-naive \
     --runtime aws.firecracker \
     --tty \
     docker.io/library/debian:latest \
     test`

func main() {
	if len(os.Args) != 2 {
		log.Fatal(`unexpected number of args, usage: "<prog> <path to template>"`)
	}
	templatePath := os.Args[1]

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("failed to open template %q: %v", templatePath, err)
	}

	err = tmpl.Execute(os.Stdout, struct {
		InstallDependencies string
		CheckoutFCCD        string
		SetupFCCD           string
		StartSnapshotter    string
		StartContainerd     string
		PullImage           string
		StartContainer      string
	}{
		installDependencies,
		checkoutFCCD,
		setupFCCD,
		startSnapshotter,
		startContainerd,
		pullImage,
		startContainer,
	})

	if err != nil {
		log.Fatalf("failed to generate template: %v", err)
	}
}
