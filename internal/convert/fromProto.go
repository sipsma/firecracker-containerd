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

package convert

import (
	"net"
	"time"

	"github.com/pkg/errors"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"

	"github.com/firecracker-microvm/firecracker-containerd/proto"
)

func MachineConfigurationFromProto(req *proto.FirecrackerMachineConfiguration) models.MachineConfiguration {
	config := models.MachineConfiguration{}
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

// NetworkConfigFromProto creates a firecracker NetworkInterface object from
// the protobuf FirecrackerNetworkInterface message.
func NetworkConfigFromProto(nwIface *proto.FirecrackerNetworkInterface, vmID string) (*firecracker.NetworkInterface, error) {
	result := &firecracker.NetworkInterface{
		AllowMMDS: nwIface.AllowMMDS,
	}

	if nwIface.InRateLimiter != nil {
		result.InRateLimiter = RateLimiterFromProto(nwIface.InRateLimiter)
	}

	if nwIface.OutRateLimiter != nil {
		result.OutRateLimiter = RateLimiterFromProto(nwIface.OutRateLimiter)
	}

	if staticConf := nwIface.StaticConfig; staticConf != nil {
		result.StaticConfiguration = &firecracker.StaticNetworkConfiguration{
			HostDevName: staticConf.HostDevName,
			MacAddress:  staticConf.MacAddress,
		}

		if ipConf := staticConf.IPConfig; ipConf != nil {
			ip, ipNet, err := net.ParseCIDR(ipConf.PrimaryAddr)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse CIDR from %q", ipConf.PrimaryAddr)
			}

			result.StaticConfiguration.IPConfiguration = &firecracker.IPConfiguration{
				IPAddr: net.IPNet{
					IP:   ip,
					Mask: ipNet.Mask,
				},
				Gateway:     net.ParseIP(ipConf.GatewayAddr),
				Nameservers: ipConf.Nameservers,
			}
		}
	}

	if cniConf := nwIface.CNIConfig; cniConf != nil {
		result.CNIConfiguration = &firecracker.CNIConfiguration{
			NetworkName: cniConf.NetworkName,
			IfName:      cniConf.InterfaceName,
			BinPath:     cniConf.BinPath,
			ConfDir:     cniConf.ConfDir,
			CacheDir:    cniConf.CacheDir,
		}

		for _, cniArg := range cniConf.Args {
			var kv [2]string
			kv[0] = cniArg.Key
			kv[1] = cniArg.Value
			result.CNIConfiguration.Args = append(result.CNIConfiguration.Args, kv)
		}
	}

	return result, nil
}

// RateLimiterFromProto creates a firecracker RateLimiter object from the
// protobuf message.
func RateLimiterFromProto(rl *proto.FirecrackerRateLimiter) *models.RateLimiter {
	result := models.RateLimiter{}
	if rl.Bandwidth != nil {
		result.Bandwidth = TokenBucketFromProto(rl.Bandwidth)
	}

	if rl.Ops != nil {
		result.Ops = TokenBucketFromProto(rl.Ops)
	}

	return &result
}

func WithRateLimiterFromProto(rl *proto.FirecrackerRateLimiter) firecracker.DriveOpt {
	if rl == nil {
		return func(d *models.Drive) {
			// no-op
		}
	}
	return firecracker.WithRateLimiter(*RateLimiterFromProto(rl))
}

// TokenBucketFromProto creates a firecracker TokenBucket object from the
// protobuf message.
func TokenBucketFromProto(bucket *proto.FirecrackerTokenBucket) *models.TokenBucket {
	builder := firecracker.TokenBucketBuilder{}
	if bucket.OneTimeBurst > 0 {
		builder = builder.WithInitialSize(bucket.OneTimeBurst)
	}

	if bucket.RefillTime > 0 {
		builder = builder.WithRefillDuration(time.Duration(bucket.RefillTime) * time.Millisecond)
	}

	if bucket.Capacity > 0 {
		builder = builder.WithBucketSize(bucket.Capacity)
	}

	res := builder.Build()
	return &res
}