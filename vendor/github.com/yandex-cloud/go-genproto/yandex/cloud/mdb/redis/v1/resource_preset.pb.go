// Code generated by protoc-gen-go. DO NOT EDIT.
// source: yandex/cloud/mdb/redis/v1/resource_preset.proto

package redis

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// A resource preset that describes hardware configuration for a host.
type ResourcePreset struct {
	// ID of the resource preset.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// IDs of availability zones where the resource preset is available.
	ZoneIds []string `protobuf:"bytes,2,rep,name=zone_ids,json=zoneIds,proto3" json:"zone_ids,omitempty"`
	// RAM volume for a Redis host created with the preset, in bytes.
	Memory               int64    `protobuf:"varint,3,opt,name=memory,proto3" json:"memory,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResourcePreset) Reset()         { *m = ResourcePreset{} }
func (m *ResourcePreset) String() string { return proto.CompactTextString(m) }
func (*ResourcePreset) ProtoMessage()    {}
func (*ResourcePreset) Descriptor() ([]byte, []int) {
	return fileDescriptor_2997d340901a0fe3, []int{0}
}

func (m *ResourcePreset) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResourcePreset.Unmarshal(m, b)
}
func (m *ResourcePreset) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResourcePreset.Marshal(b, m, deterministic)
}
func (m *ResourcePreset) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResourcePreset.Merge(m, src)
}
func (m *ResourcePreset) XXX_Size() int {
	return xxx_messageInfo_ResourcePreset.Size(m)
}
func (m *ResourcePreset) XXX_DiscardUnknown() {
	xxx_messageInfo_ResourcePreset.DiscardUnknown(m)
}

var xxx_messageInfo_ResourcePreset proto.InternalMessageInfo

func (m *ResourcePreset) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ResourcePreset) GetZoneIds() []string {
	if m != nil {
		return m.ZoneIds
	}
	return nil
}

func (m *ResourcePreset) GetMemory() int64 {
	if m != nil {
		return m.Memory
	}
	return 0
}

func init() {
	proto.RegisterType((*ResourcePreset)(nil), "yandex.cloud.mdb.redis.v1.ResourcePreset")
}

func init() {
	proto.RegisterFile("yandex/cloud/mdb/redis/v1/resource_preset.proto", fileDescriptor_2997d340901a0fe3)
}

var fileDescriptor_2997d340901a0fe3 = []byte{
	// 203 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x8f, 0x31, 0x4f, 0x80, 0x30,
	0x10, 0x85, 0x03, 0x24, 0x28, 0x1d, 0x18, 0x3a, 0x18, 0x18, 0x4c, 0x88, 0x13, 0x0b, 0x6d, 0x88,
	0xa3, 0x9b, 0x4e, 0x6e, 0xa6, 0x6e, 0x2e, 0x84, 0x72, 0x17, 0x6c, 0x62, 0x39, 0xd2, 0x02, 0x11,
	0x7f, 0xbd, 0xb1, 0xb0, 0x30, 0xb8, 0xdd, 0x4b, 0xee, 0xcb, 0x7b, 0x1f, 0x93, 0x7b, 0x3f, 0x01,
	0x7e, 0xcb, 0xe1, 0x8b, 0x56, 0x90, 0x16, 0xb4, 0x74, 0x08, 0xc6, 0xcb, 0xad, 0x95, 0x0e, 0x3d,
	0xad, 0x6e, 0xc0, 0x6e, 0x76, 0xe8, 0x71, 0x11, 0xb3, 0xa3, 0x85, 0x78, 0x79, 0x00, 0x22, 0x00,
	0xc2, 0x82, 0x16, 0x01, 0x10, 0x5b, 0xfb, 0xf0, 0xce, 0x72, 0x75, 0x32, 0x6f, 0x01, 0xe1, 0x39,
	0x8b, 0x0d, 0x14, 0x51, 0x15, 0xd5, 0x99, 0x8a, 0x0d, 0xf0, 0x92, 0xdd, 0xfe, 0xd0, 0x84, 0x9d,
	0x01, 0x5f, 0xc4, 0x55, 0x52, 0x67, 0xea, 0xe6, 0x2f, 0xbf, 0x82, 0xe7, 0x77, 0x2c, 0xb5, 0x68,
	0xc9, 0xed, 0x45, 0x52, 0x45, 0x75, 0xa2, 0xce, 0xf4, 0x0c, 0xec, 0xfe, 0xd2, 0xd8, 0xcf, 0xe6,
	0xd2, 0xfa, 0xf1, 0x32, 0x9a, 0xe5, 0x73, 0xd5, 0x62, 0x20, 0x7b, 0xca, 0x34, 0x87, 0xcc, 0x48,
	0xcd, 0x88, 0x53, 0x58, 0xfd, 0xbf, 0xe5, 0x53, 0x38, 0x74, 0x1a, 0xde, 0x1e, 0x7f, 0x03, 0x00,
	0x00, 0xff, 0xff, 0xf2, 0xc6, 0xfa, 0x9c, 0x0f, 0x01, 0x00, 0x00,
}