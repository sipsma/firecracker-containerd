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
	ExitAfterAllTasksDeleted bool     `protobuf:"varint,9,opt,name=ExitAfterAllTasksDeleted,json=exitAfterAllTasksDeleted,proto3" json:"ExitAfterAllTasksDeleted,omitempty"`
	XXX_NoUnkeyedLiteral     struct{} `json:"-"`
	XXX_unrecognized         []byte   `json:"-"`
	XXX_sizecache            int32    `json:"-"`
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

type StopVMRequest struct {
	VMID                 string   `protobuf:"bytes,1,opt,name=VMID,json=vMID,proto3" json:"VMID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopVMRequest) Reset()         { *m = StopVMRequest{} }
func (m *StopVMRequest) String() string { return proto.CompactTextString(m) }
func (*StopVMRequest) ProtoMessage()    {}
func (*StopVMRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73317e9fb8da571, []int{1}
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
	return fileDescriptor_a73317e9fb8da571, []int{2}
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
	ContextID            uint32   `protobuf:"varint,2,opt,name=ContextID,json=contextID,proto3" json:"ContextID,omitempty"`
	SocketPath           string   `protobuf:"bytes,3,opt,name=SocketPath,json=socketPath,proto3" json:"SocketPath,omitempty"`
	LogFifoPath          string   `protobuf:"bytes,4,opt,name=LogFifoPath,json=logFifoPath,proto3" json:"LogFifoPath,omitempty"`
	MetricsFifoPath      string   `protobuf:"bytes,5,opt,name=MetricsFifoPath,json=metricsFifoPath,proto3" json:"MetricsFifoPath,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetVMInfoResponse) Reset()         { *m = GetVMInfoResponse{} }
func (m *GetVMInfoResponse) String() string { return proto.CompactTextString(m) }
func (*GetVMInfoResponse) ProtoMessage()    {}
func (*GetVMInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73317e9fb8da571, []int{3}
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

func (m *GetVMInfoResponse) GetContextID() uint32 {
	if m != nil {
		return m.ContextID
	}
	return 0
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
	return fileDescriptor_a73317e9fb8da571, []int{4}
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

func init() {
	proto.RegisterType((*CreateVMRequest)(nil), "CreateVMRequest")
	proto.RegisterType((*StopVMRequest)(nil), "StopVMRequest")
	proto.RegisterType((*GetVMInfoRequest)(nil), "GetVMInfoRequest")
	proto.RegisterType((*GetVMInfoResponse)(nil), "GetVMInfoResponse")
	proto.RegisterType((*SetVMMetadataRequest)(nil), "SetVMMetadataRequest")
}

func init() { proto.RegisterFile("firecracker.proto", fileDescriptor_a73317e9fb8da571) }

var fileDescriptor_a73317e9fb8da571 = []byte{
	// 450 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x93, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x86, 0x65, 0x9a, 0xb4, 0xf1, 0x58, 0x6d, 0xc8, 0x0a, 0xd0, 0xaa, 0xaa, 0x90, 0x15, 0xa4,
	0xca, 0xa7, 0x1e, 0xda, 0x0b, 0xe2, 0x44, 0x48, 0x08, 0x32, 0x60, 0x84, 0x36, 0x28, 0x07, 0x6e,
	0x8b, 0x33, 0x76, 0x2d, 0x3b, 0xbb, 0x61, 0x77, 0x52, 0xca, 0x63, 0xf1, 0x26, 0x3c, 0x12, 0xf2,
	0x86, 0x38, 0x4e, 0x05, 0x39, 0x79, 0xf7, 0x9f, 0xef, 0x1f, 0x8f, 0x35, 0xbf, 0x61, 0x90, 0x15,
	0x06, 0x53, 0x23, 0xd3, 0x12, 0xcd, 0xd5, 0xca, 0x68, 0xd2, 0xe7, 0x01, 0xfd, 0x5c, 0xa1, 0xdd,
	0x5c, 0x86, 0xbf, 0x8f, 0xa0, 0x3f, 0x36, 0x28, 0x09, 0xe7, 0x89, 0xc0, 0xef, 0x6b, 0xb4, 0xc4,
	0x18, 0x74, 0xe6, 0x49, 0x3c, 0xe1, 0x5e, 0xe8, 0x45, 0xbe, 0xe8, 0xdc, 0x25, 0xf1, 0x84, 0xbd,
	0x06, 0x48, 0x64, 0x7a, 0x5b, 0x28, 0x1c, 0x67, 0x39, 0x7f, 0x14, 0x7a, 0x51, 0x70, 0x1d, 0x5e,
	0x4d, 0x77, 0xcd, 0xb7, 0x55, 0xad, 0xb2, 0x22, 0x5f, 0x1b, 0x49, 0x85, 0x56, 0x02, 0x96, 0x8d,
	0x87, 0x45, 0xd0, 0xff, 0x80, 0x46, 0x61, 0x15, 0x2f, 0x65, 0x8e, 0x9f, 0x25, 0xdd, 0xf2, 0x23,
	0xf7, 0x82, 0x7e, 0xb9, 0x2f, 0xb3, 0xe7, 0x00, 0x1b, 0x72, 0x64, 0x72, 0xcb, 0x3b, 0x0e, 0x82,
	0xb2, 0x51, 0xd8, 0x0d, 0xf8, 0x42, 0x6b, 0x9a, 0x98, 0xe2, 0x0e, 0x79, 0xd7, 0x8d, 0xf2, 0xb4,
	0x3d, 0x4a, 0x53, 0x14, 0xbe, 0xd9, 0x1e, 0xd9, 0x4b, 0x08, 0xdc, 0x21, 0xd1, 0x6b, 0x45, 0x96,
	0x1f, 0x87, 0x47, 0x51, 0x70, 0xfd, 0xac, 0x6d, 0xdb, 0x95, 0x45, 0xb0, 0xd8, 0xa1, 0xec, 0x3d,
	0x0c, 0x3e, 0x21, 0xfd, 0xd0, 0xa6, 0x8c, 0x15, 0xa1, 0xc9, 0x64, 0x8a, 0x96, 0x9f, 0x38, 0xff,
	0x45, 0xdb, 0xff, 0x10, 0x12, 0x03, 0xf5, 0xd0, 0xc6, 0x2e, 0xe1, 0x6c, 0xac, 0x15, 0xc9, 0x42,
	0xa1, 0x19, 0xd7, 0xed, 0x79, 0x2f, 0xf4, 0xa2, 0xae, 0x38, 0x4b, 0xf7, 0x54, 0xf6, 0x0a, 0xf8,
	0xdb, 0xfb, 0x82, 0x46, 0x19, 0xa1, 0x19, 0x55, 0xd5, 0x17, 0x69, 0x4b, 0x3b, 0xc1, 0x0a, 0x09,
	0x17, 0xdc, 0x0f, 0xbd, 0xa8, 0x27, 0x38, 0xfe, 0xa7, 0x3e, 0x7c, 0x01, 0xa7, 0x33, 0xd2, 0xab,
	0x83, 0xfb, 0x1c, 0x5e, 0xc2, 0xe3, 0x77, 0x48, 0xf3, 0x24, 0x56, 0x99, 0x3e, 0xc4, 0xfd, 0xf2,
	0x60, 0xd0, 0x02, 0xed, 0x4a, 0x2b, 0x8b, 0xff, 0x4c, 0xc8, 0x05, 0xf8, 0xf5, 0xa7, 0xe1, 0x3d,
	0xc5, 0x13, 0x17, 0x90, 0x53, 0xe1, 0xa7, 0x5b, 0xa1, 0xde, 0xe9, 0x4c, 0xa7, 0x25, 0x52, 0x6b,
	0xf1, 0x60, 0x1b, 0x85, 0x85, 0x10, 0x7c, 0xd4, 0xf9, 0xb4, 0xc8, 0xb4, 0x03, 0x36, 0x4b, 0x0f,
	0xaa, 0x9d, 0x54, 0xe7, 0x27, 0x41, 0x32, 0x45, 0x6a, 0x1b, 0xaa, 0xbb, 0xc9, 0xcf, 0x72, 0x5f,
	0x1e, 0x4e, 0xe1, 0xc9, 0xac, 0x1e, 0x39, 0x41, 0x92, 0x0b, 0x49, 0xf2, 0x50, 0xae, 0xcf, 0xa1,
	0xb7, 0xc5, 0xdc, 0xd0, 0xbe, 0xe8, 0x2d, 0xff, 0xde, 0xdf, 0x9c, 0x7c, 0xed, 0xba, 0x9f, 0xe4,
	0xdb, 0xb1, 0x7b, 0xdc, 0xfc, 0x09, 0x00, 0x00, 0xff, 0xff, 0xa0, 0xec, 0x78, 0x11, 0x4d, 0x03,
	0x00, 0x00,
}
