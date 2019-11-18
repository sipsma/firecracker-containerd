// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: firecracker.proto

package proto

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// CreateVMRequest specifies creation parameters for a new FC instance
type CreateVMRequest struct {
	// VM identifier to assign
	VMID string `protobuf:"bytes,1,opt,name=VMID,json=vMID,proto3" json:"VMID,omitempty"`
	// Specifies the machine configuration for the VM
	MachineCfg *FirecrackerMachineConfiguration `protobuf:"bytes,2,opt,name=MachineCfg,json=machineCfg,proto3" json:"MachineCfg,omitempty"`
	// Specifies the file path where the kernel image is located
	KernelImagePath string `protobuf:"bytes,3,opt,name=KernelImagePath,json=kernelImagePath,proto3" json:"KernelImagePath,omitempty"`
	// Specifies the commandline arguments that should be passed to the kernel
	KernelArgs string `protobuf:"bytes,4,opt,name=KernelArgs,json=kernelArgs,proto3" json:"KernelArgs,omitempty"`
	// Specifies drive containing the rootfs of the VM
	RootDrive *FirecrackerRootDrive `protobuf:"bytes,5,opt,name=RootDrive,json=rootDrive,proto3" json:"RootDrive,omitempty"`
	// Specifies additional drives whose contents will be mounted inside the VM on boot.
	DriveMounts []*FirecrackerDriveMount `protobuf:"bytes,6,rep,name=DriveMounts,json=driveMounts,proto3" json:"DriveMounts,omitempty"`
	// Specifies the networking configuration for a VM
	NetworkInterfaces []*FirecrackerNetworkInterface `protobuf:"bytes,7,rep,name=NetworkInterfaces,json=networkInterfaces,proto3" json:"NetworkInterfaces,omitempty"`
	// The number of dummy drives to reserve in advance before running FC instance.
	ContainerCount int32 `protobuf:"varint,8,opt,name=ContainerCount,json=containerCount,proto3" json:"ContainerCount,omitempty"`
	// Whether the VM should exit after all tasks running in it have been deleted.
	ExitAfterAllTasksDeleted bool          `protobuf:"varint,9,opt,name=ExitAfterAllTasksDeleted,json=exitAfterAllTasksDeleted,proto3" json:"ExitAfterAllTasksDeleted,omitempty"`
	JailerConfig             *JailerConfig `protobuf:"bytes,10,opt,name=JailerConfig,json=jailerConfig,proto3" json:"JailerConfig,omitempty"`
	XXX_NoUnkeyedLiteral     struct{}      `json:"-"`
	XXX_unrecognized         []byte        `json:"-"`
	XXX_sizecache            int32         `json:"-"`
}

func (m *CreateVMRequest) Reset()         { *m = CreateVMRequest{} }
func (m *CreateVMRequest) String() string { return proto.CompactTextString(m) }
func (*CreateVMRequest) ProtoMessage()    {}
func (*CreateVMRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73317e9fb8da571, []int{0}
}
func (m *CreateVMRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateVMRequest.Unmarshal(m, b)
}
func (m *CreateVMRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateVMRequest.Marshal(b, m, deterministic)
}
func (m *CreateVMRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateVMRequest.Merge(m, src)
}
func (m *CreateVMRequest) XXX_Size() int {
	return xxx_messageInfo_CreateVMRequest.Size(m)
}
func (m *CreateVMRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateVMRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateVMRequest proto.InternalMessageInfo

func (m *CreateVMRequest) GetVMID() string {
	if m != nil {
		return m.VMID
	}
	return ""
}

func (m *CreateVMRequest) GetMachineCfg() *FirecrackerMachineConfiguration {
	if m != nil {
		return m.MachineCfg
	}
	return nil
}

func (m *CreateVMRequest) GetKernelImagePath() string {
	if m != nil {
		return m.KernelImagePath
	}
	return ""
}

func (m *CreateVMRequest) GetKernelArgs() string {
	if m != nil {
		return m.KernelArgs
	}
	return ""
}

func (m *CreateVMRequest) GetRootDrive() *FirecrackerRootDrive {
	if m != nil {
		return m.RootDrive
	}
	return nil
}

func (m *CreateVMRequest) GetDriveMounts() []*FirecrackerDriveMount {
	if m != nil {
		return m.DriveMounts
	}
	return nil
}

func (m *CreateVMRequest) GetNetworkInterfaces() []*FirecrackerNetworkInterface {
	if m != nil {
		return m.NetworkInterfaces
	}
	return nil
}

func (m *CreateVMRequest) GetContainerCount() int32 {
	if m != nil {
		return m.ContainerCount
	}
	return 0
}

func (m *CreateVMRequest) GetExitAfterAllTasksDeleted() bool {
	if m != nil {
		return m.ExitAfterAllTasksDeleted
	}
	return false
}

func (m *CreateVMRequest) GetJailerConfig() *JailerConfig {
	if m != nil {
		return m.JailerConfig
	}
	return nil
}

type CreateVMResponse struct {
	VMID                 string   `protobuf:"bytes,1,opt,name=VMID,json=vMID,proto3" json:"VMID,omitempty"`
	SocketPath           string   `protobuf:"bytes,2,opt,name=SocketPath,json=socketPath,proto3" json:"SocketPath,omitempty"`
	LogFifoPath          string   `protobuf:"bytes,3,opt,name=LogFifoPath,json=logFifoPath,proto3" json:"LogFifoPath,omitempty"`
	MetricsFifoPath      string   `protobuf:"bytes,4,opt,name=MetricsFifoPath,json=metricsFifoPath,proto3" json:"MetricsFifoPath,omitempty"`
	CgroupPath           string   `protobuf:"bytes,5,opt,name=CgroupPath,json=cgroupPath,proto3" json:"CgroupPath,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateVMResponse) Reset()         { *m = CreateVMResponse{} }
func (m *CreateVMResponse) String() string { return proto.CompactTextString(m) }
func (*CreateVMResponse) ProtoMessage()    {}
func (*CreateVMResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73317e9fb8da571, []int{1}
}
func (m *CreateVMResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateVMResponse.Unmarshal(m, b)
}
func (m *CreateVMResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateVMResponse.Marshal(b, m, deterministic)
}
func (m *CreateVMResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateVMResponse.Merge(m, src)
}
func (m *CreateVMResponse) XXX_Size() int {
	return xxx_messageInfo_CreateVMResponse.Size(m)
}
func (m *CreateVMResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateVMResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateVMResponse proto.InternalMessageInfo

func (m *CreateVMResponse) GetVMID() string {
	if m != nil {
		return m.VMID
	}
	return ""
}

func (m *CreateVMResponse) GetSocketPath() string {
	if m != nil {
		return m.SocketPath
	}
	return ""
}

func (m *CreateVMResponse) GetLogFifoPath() string {
	if m != nil {
		return m.LogFifoPath
	}
	return ""
}

func (m *CreateVMResponse) GetMetricsFifoPath() string {
	if m != nil {
		return m.MetricsFifoPath
	}
	return ""
}

func (m *CreateVMResponse) GetCgroupPath() string {
	if m != nil {
		return m.CgroupPath
	}
	return ""
}

type StopVMRequest struct {
	VMID                 string   `protobuf:"bytes,1,opt,name=VMID,json=vMID,proto3" json:"VMID,omitempty"`
	TimeoutSeconds       uint32   `protobuf:"varint,2,opt,name=TimeoutSeconds,json=timeoutSeconds,proto3" json:"TimeoutSeconds,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopVMRequest) Reset()         { *m = StopVMRequest{} }
func (m *StopVMRequest) String() string { return proto.CompactTextString(m) }
func (*StopVMRequest) ProtoMessage()    {}
func (*StopVMRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73317e9fb8da571, []int{2}
}
func (m *StopVMRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopVMRequest.Unmarshal(m, b)
}
func (m *StopVMRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopVMRequest.Marshal(b, m, deterministic)
}
func (m *StopVMRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopVMRequest.Merge(m, src)
}
func (m *StopVMRequest) XXX_Size() int {
	return xxx_messageInfo_StopVMRequest.Size(m)
}
func (m *StopVMRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StopVMRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StopVMRequest proto.InternalMessageInfo

func (m *StopVMRequest) GetVMID() string {
	if m != nil {
		return m.VMID
	}
	return ""
}

func (m *StopVMRequest) GetTimeoutSeconds() uint32 {
	if m != nil {
		return m.TimeoutSeconds
	}
	return 0
}

type GetVMInfoRequest struct {
	VMID                 string   `protobuf:"bytes,1,opt,name=VMID,json=vMID,proto3" json:"VMID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetVMInfoRequest) Reset()         { *m = GetVMInfoRequest{} }
func (m *GetVMInfoRequest) String() string { return proto.CompactTextString(m) }
func (*GetVMInfoRequest) ProtoMessage()    {}
func (*GetVMInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73317e9fb8da571, []int{3}
}
func (m *GetVMInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetVMInfoRequest.Unmarshal(m, b)
}
func (m *GetVMInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetVMInfoRequest.Marshal(b, m, deterministic)
}
func (m *GetVMInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetVMInfoRequest.Merge(m, src)
}
func (m *GetVMInfoRequest) XXX_Size() int {
	return xxx_messageInfo_GetVMInfoRequest.Size(m)
}
func (m *GetVMInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetVMInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetVMInfoRequest proto.InternalMessageInfo

func (m *GetVMInfoRequest) GetVMID() string {
	if m != nil {
		return m.VMID
	}
	return ""
}

type GetVMInfoResponse struct {
	VMID                 string   `protobuf:"bytes,1,opt,name=VMID,json=vMID,proto3" json:"VMID,omitempty"`
	SocketPath           string   `protobuf:"bytes,2,opt,name=SocketPath,json=socketPath,proto3" json:"SocketPath,omitempty"`
	LogFifoPath          string   `protobuf:"bytes,3,opt,name=LogFifoPath,json=logFifoPath,proto3" json:"LogFifoPath,omitempty"`
	MetricsFifoPath      string   `protobuf:"bytes,4,opt,name=MetricsFifoPath,json=metricsFifoPath,proto3" json:"MetricsFifoPath,omitempty"`
	CgroupPath           string   `protobuf:"bytes,5,opt,name=CgroupPath,json=cgroupPath,proto3" json:"CgroupPath,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetVMInfoResponse) Reset()         { *m = GetVMInfoResponse{} }
func (m *GetVMInfoResponse) String() string { return proto.CompactTextString(m) }
func (*GetVMInfoResponse) ProtoMessage()    {}
func (*GetVMInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73317e9fb8da571, []int{4}
}
func (m *GetVMInfoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetVMInfoResponse.Unmarshal(m, b)
}
func (m *GetVMInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetVMInfoResponse.Marshal(b, m, deterministic)
}
func (m *GetVMInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetVMInfoResponse.Merge(m, src)
}
func (m *GetVMInfoResponse) XXX_Size() int {
	return xxx_messageInfo_GetVMInfoResponse.Size(m)
}
func (m *GetVMInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetVMInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetVMInfoResponse proto.InternalMessageInfo

func (m *GetVMInfoResponse) GetVMID() string {
	if m != nil {
		return m.VMID
	}
	return ""
}

func (m *GetVMInfoResponse) GetSocketPath() string {
	if m != nil {
		return m.SocketPath
	}
	return ""
}

func (m *GetVMInfoResponse) GetLogFifoPath() string {
	if m != nil {
		return m.LogFifoPath
	}
	return ""
}

func (m *GetVMInfoResponse) GetMetricsFifoPath() string {
	if m != nil {
		return m.MetricsFifoPath
	}
	return ""
}

func (m *GetVMInfoResponse) GetCgroupPath() string {
	if m != nil {
		return m.CgroupPath
	}
	return ""
}

type SetVMMetadataRequest struct {
	VMID                 string   `protobuf:"bytes,1,opt,name=VMID,json=vMID,proto3" json:"VMID,omitempty"`
	Metadata             string   `protobuf:"bytes,2,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetVMMetadataRequest) Reset()         { *m = SetVMMetadataRequest{} }
func (m *SetVMMetadataRequest) String() string { return proto.CompactTextString(m) }
func (*SetVMMetadataRequest) ProtoMessage()    {}
func (*SetVMMetadataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73317e9fb8da571, []int{5}
}
func (m *SetVMMetadataRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetVMMetadataRequest.Unmarshal(m, b)
}
func (m *SetVMMetadataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetVMMetadataRequest.Marshal(b, m, deterministic)
}
func (m *SetVMMetadataRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetVMMetadataRequest.Merge(m, src)
}
func (m *SetVMMetadataRequest) XXX_Size() int {
	return xxx_messageInfo_SetVMMetadataRequest.Size(m)
}
func (m *SetVMMetadataRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetVMMetadataRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetVMMetadataRequest proto.InternalMessageInfo

func (m *SetVMMetadataRequest) GetVMID() string {
	if m != nil {
		return m.VMID
	}
	return ""
}

func (m *SetVMMetadataRequest) GetMetadata() string {
	if m != nil {
		return m.Metadata
	}
	return ""
}

type UpdateVMMetadataRequest struct {
	VMID                 string   `protobuf:"bytes,1,opt,name=VMID,json=vMID,proto3" json:"VMID,omitempty"`
	Metadata             string   `protobuf:"bytes,2,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateVMMetadataRequest) Reset()         { *m = UpdateVMMetadataRequest{} }
func (m *UpdateVMMetadataRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateVMMetadataRequest) ProtoMessage()    {}
func (*UpdateVMMetadataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73317e9fb8da571, []int{6}
}
func (m *UpdateVMMetadataRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateVMMetadataRequest.Unmarshal(m, b)
}
func (m *UpdateVMMetadataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateVMMetadataRequest.Marshal(b, m, deterministic)
}
func (m *UpdateVMMetadataRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateVMMetadataRequest.Merge(m, src)
}
func (m *UpdateVMMetadataRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateVMMetadataRequest.Size(m)
}
func (m *UpdateVMMetadataRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateVMMetadataRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateVMMetadataRequest proto.InternalMessageInfo

func (m *UpdateVMMetadataRequest) GetVMID() string {
	if m != nil {
		return m.VMID
	}
	return ""
}

func (m *UpdateVMMetadataRequest) GetMetadata() string {
	if m != nil {
		return m.Metadata
	}
	return ""
}

type GetVMMetadataRequest struct {
	VMID                 string   `protobuf:"bytes,1,opt,name=VMID,json=vMID,proto3" json:"VMID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetVMMetadataRequest) Reset()         { *m = GetVMMetadataRequest{} }
func (m *GetVMMetadataRequest) String() string { return proto.CompactTextString(m) }
func (*GetVMMetadataRequest) ProtoMessage()    {}
func (*GetVMMetadataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73317e9fb8da571, []int{7}
}
func (m *GetVMMetadataRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetVMMetadataRequest.Unmarshal(m, b)
}
func (m *GetVMMetadataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetVMMetadataRequest.Marshal(b, m, deterministic)
}
func (m *GetVMMetadataRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetVMMetadataRequest.Merge(m, src)
}
func (m *GetVMMetadataRequest) XXX_Size() int {
	return xxx_messageInfo_GetVMMetadataRequest.Size(m)
}
func (m *GetVMMetadataRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetVMMetadataRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetVMMetadataRequest proto.InternalMessageInfo

func (m *GetVMMetadataRequest) GetVMID() string {
	if m != nil {
		return m.VMID
	}
	return ""
}

type GetVMMetadataResponse struct {
	Metadata             string   `protobuf:"bytes,1,opt,name=Metadata,json=metadata,proto3" json:"Metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetVMMetadataResponse) Reset()         { *m = GetVMMetadataResponse{} }
func (m *GetVMMetadataResponse) String() string { return proto.CompactTextString(m) }
func (*GetVMMetadataResponse) ProtoMessage()    {}
func (*GetVMMetadataResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73317e9fb8da571, []int{8}
}
func (m *GetVMMetadataResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetVMMetadataResponse.Unmarshal(m, b)
}
func (m *GetVMMetadataResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetVMMetadataResponse.Marshal(b, m, deterministic)
}
func (m *GetVMMetadataResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetVMMetadataResponse.Merge(m, src)
}
func (m *GetVMMetadataResponse) XXX_Size() int {
	return xxx_messageInfo_GetVMMetadataResponse.Size(m)
}
func (m *GetVMMetadataResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetVMMetadataResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetVMMetadataResponse proto.InternalMessageInfo

func (m *GetVMMetadataResponse) GetMetadata() string {
	if m != nil {
		return m.Metadata
	}
	return ""
}

type JailerConfig struct {
	NetNS                string   `protobuf:"bytes,1,opt,name=NetNS,json=netNS,proto3" json:"NetNS,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JailerConfig) Reset()         { *m = JailerConfig{} }
func (m *JailerConfig) String() string { return proto.CompactTextString(m) }
func (*JailerConfig) ProtoMessage()    {}
func (*JailerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73317e9fb8da571, []int{9}
}
func (m *JailerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JailerConfig.Unmarshal(m, b)
}
func (m *JailerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JailerConfig.Marshal(b, m, deterministic)
}
func (m *JailerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JailerConfig.Merge(m, src)
}
func (m *JailerConfig) XXX_Size() int {
	return xxx_messageInfo_JailerConfig.Size(m)
}
func (m *JailerConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_JailerConfig.DiscardUnknown(m)
}

var xxx_messageInfo_JailerConfig proto.InternalMessageInfo

func (m *JailerConfig) GetNetNS() string {
	if m != nil {
		return m.NetNS
	}
	return ""
}

func init() {
	proto.RegisterType((*CreateVMRequest)(nil), "CreateVMRequest")
	proto.RegisterType((*CreateVMResponse)(nil), "CreateVMResponse")
	proto.RegisterType((*StopVMRequest)(nil), "StopVMRequest")
	proto.RegisterType((*GetVMInfoRequest)(nil), "GetVMInfoRequest")
	proto.RegisterType((*GetVMInfoResponse)(nil), "GetVMInfoResponse")
	proto.RegisterType((*SetVMMetadataRequest)(nil), "SetVMMetadataRequest")
	proto.RegisterType((*UpdateVMMetadataRequest)(nil), "UpdateVMMetadataRequest")
	proto.RegisterType((*GetVMMetadataRequest)(nil), "GetVMMetadataRequest")
	proto.RegisterType((*GetVMMetadataResponse)(nil), "GetVMMetadataResponse")
	proto.RegisterType((*JailerConfig)(nil), "JailerConfig")
}

func init() { proto.RegisterFile("firecracker.proto", fileDescriptor_a73317e9fb8da571) }

var fileDescriptor_a73317e9fb8da571 = []byte{
	// 550 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x54, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x95, 0xdb, 0xb8, 0x4d, 0xc6, 0x6d, 0xd2, 0xac, 0x5a, 0x58, 0x55, 0xa8, 0xb2, 0x2c, 0x54,
	0x59, 0x1c, 0x2a, 0xd1, 0x5e, 0x10, 0x27, 0x4a, 0x42, 0x2b, 0xb7, 0xb8, 0x42, 0x4e, 0xe9, 0x81,
	0xdb, 0xe2, 0x8c, 0x5d, 0x63, 0x67, 0xd7, 0xec, 0x6e, 0x0a, 0x9c, 0xf9, 0x23, 0xf8, 0x41, 0xe4,
	0x4d, 0xea, 0x38, 0x11, 0x54, 0x48, 0x9c, 0x38, 0x65, 0xe7, 0xcd, 0x9b, 0x99, 0xa7, 0x99, 0x17,
	0x43, 0x3f, 0xc9, 0x24, 0xc6, 0x92, 0xc5, 0x39, 0xca, 0xa3, 0x52, 0x0a, 0x2d, 0xf6, 0x1d, 0xfd,
	0xad, 0x44, 0x35, 0x0b, 0xbc, 0xef, 0x2d, 0xe8, 0x0d, 0x24, 0x32, 0x8d, 0x37, 0x61, 0x84, 0x9f,
	0xa7, 0xa8, 0x34, 0x21, 0xd0, 0xba, 0x09, 0x83, 0x21, 0xb5, 0x5c, 0xcb, 0xef, 0x44, 0xad, 0xbb,
	0x30, 0x18, 0x92, 0x57, 0x00, 0x21, 0x8b, 0x6f, 0x33, 0x8e, 0x83, 0x24, 0xa5, 0x6b, 0xae, 0xe5,
	0x3b, 0xc7, 0xee, 0xd1, 0xd9, 0xa2, 0xf9, 0x7d, 0x56, 0xf0, 0x24, 0x4b, 0xa7, 0x92, 0xe9, 0x4c,
	0xf0, 0x08, 0x26, 0x75, 0x0d, 0xf1, 0xa1, 0x77, 0x89, 0x92, 0x63, 0x11, 0x4c, 0x58, 0x8a, 0xef,
	0x98, 0xbe, 0xa5, 0xeb, 0x66, 0x40, 0x2f, 0x5f, 0x86, 0xc9, 0x01, 0xc0, 0x8c, 0x79, 0x2a, 0x53,
	0x45, 0x5b, 0x86, 0x04, 0x79, 0x8d, 0x90, 0x13, 0xe8, 0x44, 0x42, 0xe8, 0xa1, 0xcc, 0xee, 0x90,
	0xda, 0x46, 0xca, 0x5e, 0x53, 0x4a, 0x9d, 0x8c, 0x3a, 0xf2, 0xfe, 0x49, 0x5e, 0x80, 0x63, 0x1e,
	0xa1, 0x98, 0x72, 0xad, 0xe8, 0x86, 0xbb, 0xee, 0x3b, 0xc7, 0x8f, 0x9a, 0x65, 0x8b, 0x74, 0xe4,
	0x8c, 0x17, 0x54, 0x72, 0x01, 0xfd, 0x2b, 0xd4, 0x5f, 0x84, 0xcc, 0x03, 0xae, 0x51, 0x26, 0x2c,
	0x46, 0x45, 0x37, 0x4d, 0xfd, 0x93, 0x66, 0xfd, 0x2a, 0x29, 0xea, 0xf3, 0xd5, 0x32, 0x72, 0x08,
	0xdd, 0x81, 0xe0, 0x9a, 0x65, 0x1c, 0xe5, 0xa0, 0x6a, 0x4f, 0xdb, 0xae, 0xe5, 0xdb, 0x51, 0x37,
	0x5e, 0x42, 0xc9, 0x4b, 0xa0, 0x6f, 0xbe, 0x66, 0xfa, 0x34, 0xd1, 0x28, 0x4f, 0x8b, 0xe2, 0x9a,
	0xa9, 0x5c, 0x0d, 0xb1, 0x40, 0x8d, 0x63, 0xda, 0x71, 0x2d, 0xbf, 0x1d, 0x51, 0xfc, 0x43, 0x9e,
	0x3c, 0x87, 0xad, 0x0b, 0x96, 0x15, 0x55, 0xab, 0xea, 0x16, 0x14, 0xcc, 0x86, 0xb6, 0x8f, 0x9a,
	0x60, 0xb4, 0xf5, 0xa9, 0x11, 0x79, 0x3f, 0x2c, 0xd8, 0x59, 0xb8, 0x40, 0x95, 0x82, 0x2b, 0xfc,
	0xad, 0x0d, 0x0e, 0x00, 0x46, 0x22, 0xce, 0x51, 0x9b, 0xfb, 0xad, 0xcd, 0x4e, 0xa3, 0x6a, 0x84,
	0xb8, 0xe0, 0xbc, 0x15, 0xe9, 0x59, 0x96, 0x88, 0xc6, 0x81, 0x9d, 0x62, 0x01, 0x55, 0x36, 0x08,
	0x51, 0xcb, 0x2c, 0x56, 0x35, 0x6b, 0x76, 0xe1, 0xde, 0x64, 0x19, 0xae, 0x66, 0x0d, 0x52, 0x29,
	0xa6, 0xa5, 0x21, 0xd9, 0xb3, 0x59, 0x71, 0x8d, 0x78, 0x97, 0xb0, 0x3d, 0xd2, 0xa2, 0x7c, 0xd8,
	0xb7, 0x87, 0xd0, 0xbd, 0xce, 0x26, 0x28, 0xa6, 0x7a, 0x84, 0xb1, 0xe0, 0x63, 0x65, 0x44, 0x6f,
	0x47, 0x5d, 0xbd, 0x84, 0x7a, 0x87, 0xb0, 0x73, 0x8e, 0xfa, 0x26, 0x0c, 0x78, 0x22, 0x1e, 0xe8,
	0xe7, 0xfd, 0xb4, 0xa0, 0xdf, 0x20, 0xfe, 0x27, 0xab, 0x3a, 0x83, 0xdd, 0x51, 0x25, 0x3a, 0x44,
	0xcd, 0xc6, 0x4c, 0xb3, 0x87, 0x36, 0xb6, 0x0f, 0xed, 0x7b, 0xda, 0x5c, 0x75, 0x7b, 0x32, 0x8f,
	0xbd, 0x00, 0x1e, 0xbf, 0x2f, 0xc7, 0xc6, 0x26, 0xff, 0xda, 0xea, 0x19, 0xec, 0x9e, 0xff, 0xa5,
	0x24, 0xef, 0x04, 0xf6, 0x56, 0xb8, 0xf3, 0xbd, 0x37, 0x07, 0x58, 0x2b, 0x03, 0x9e, 0x2e, 0xff,
	0x0d, 0xc8, 0x2e, 0xd8, 0x57, 0xa8, 0xaf, 0x46, 0x73, 0xa2, 0xcd, 0xab, 0xe0, 0xf5, 0xe6, 0x07,
	0xdb, 0x7c, 0x08, 0x3f, 0x6e, 0x98, 0x9f, 0x93, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x62,
	0xbf, 0xbf, 0x31, 0x05, 0x00, 0x00,
}
