// Code generated by protoc-gen-go. DO NOT EDIT.
// source: yandex/cloud/compute/v1/snapshot.proto

package compute

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Snapshot_Status int32

const (
	Snapshot_STATUS_UNSPECIFIED Snapshot_Status = 0
	// Snapshot is being created.
	Snapshot_CREATING Snapshot_Status = 1
	// Snapshot is ready to use.
	Snapshot_READY Snapshot_Status = 2
	// Snapshot encountered a problem and cannot operate.
	Snapshot_ERROR Snapshot_Status = 3
	// Snapshot is being deleted.
	Snapshot_DELETING Snapshot_Status = 4
)

var Snapshot_Status_name = map[int32]string{
	0: "STATUS_UNSPECIFIED",
	1: "CREATING",
	2: "READY",
	3: "ERROR",
	4: "DELETING",
}

var Snapshot_Status_value = map[string]int32{
	"STATUS_UNSPECIFIED": 0,
	"CREATING":           1,
	"READY":              2,
	"ERROR":              3,
	"DELETING":           4,
}

func (x Snapshot_Status) String() string {
	return proto.EnumName(Snapshot_Status_name, int32(x))
}

func (Snapshot_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cac027a1005d8550, []int{0, 0}
}

// A Snapshot resource. For more information, see [Snapshots](/docs/compute/concepts/snapshot).
type Snapshot struct {
	// ID of the snapshot.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// ID of the folder that the snapshot belongs to.
	FolderId  string               `protobuf:"bytes,2,opt,name=folder_id,json=folderId,proto3" json:"folder_id,omitempty"`
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// Name of the snapshot. 1-63 characters long.
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the snapshot. 0-256 characters long.
	Description string `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	// Resource labels as `key:value` pairs. Maximum of 64 per resource.
	Labels map[string]string `protobuf:"bytes,6,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Size of the snapshot, specified in bytes.
	StorageSize int64 `protobuf:"varint,7,opt,name=storage_size,json=storageSize,proto3" json:"storage_size,omitempty"`
	// Size of the disk when the snapshot was created, specified in bytes.
	DiskSize int64 `protobuf:"varint,8,opt,name=disk_size,json=diskSize,proto3" json:"disk_size,omitempty"`
	// License IDs that indicate which licenses are attached to this resource.
	// License IDs are used to calculate additional charges for the use of the virtual machine.
	//
	// The correct license ID is generated by Yandex.Cloud. IDs are inherited by new resources created from this resource.
	//
	// If you know the license IDs, specify them when you create the image.
	// For example, if you create a disk image using a third-party utility and load it into Yandex Object Storage, the license IDs will be lost.
	// You can specify them in the [yandex.cloud.compute.v1.ImageService.Create] request.
	ProductIds []string `protobuf:"bytes,9,rep,name=product_ids,json=productIds,proto3" json:"product_ids,omitempty"`
	// Current status of the snapshot.
	Status Snapshot_Status `protobuf:"varint,10,opt,name=status,proto3,enum=yandex.cloud.compute.v1.Snapshot_Status" json:"status,omitempty"`
	// ID of the source disk used to create this snapshot.
	SourceDiskId         string   `protobuf:"bytes,11,opt,name=source_disk_id,json=sourceDiskId,proto3" json:"source_disk_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Snapshot) Reset()         { *m = Snapshot{} }
func (m *Snapshot) String() string { return proto.CompactTextString(m) }
func (*Snapshot) ProtoMessage()    {}
func (*Snapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_cac027a1005d8550, []int{0}
}

func (m *Snapshot) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Snapshot.Unmarshal(m, b)
}
func (m *Snapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Snapshot.Marshal(b, m, deterministic)
}
func (m *Snapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Snapshot.Merge(m, src)
}
func (m *Snapshot) XXX_Size() int {
	return xxx_messageInfo_Snapshot.Size(m)
}
func (m *Snapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_Snapshot.DiscardUnknown(m)
}

var xxx_messageInfo_Snapshot proto.InternalMessageInfo

func (m *Snapshot) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Snapshot) GetFolderId() string {
	if m != nil {
		return m.FolderId
	}
	return ""
}

func (m *Snapshot) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *Snapshot) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Snapshot) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Snapshot) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *Snapshot) GetStorageSize() int64 {
	if m != nil {
		return m.StorageSize
	}
	return 0
}

func (m *Snapshot) GetDiskSize() int64 {
	if m != nil {
		return m.DiskSize
	}
	return 0
}

func (m *Snapshot) GetProductIds() []string {
	if m != nil {
		return m.ProductIds
	}
	return nil
}

func (m *Snapshot) GetStatus() Snapshot_Status {
	if m != nil {
		return m.Status
	}
	return Snapshot_STATUS_UNSPECIFIED
}

func (m *Snapshot) GetSourceDiskId() string {
	if m != nil {
		return m.SourceDiskId
	}
	return ""
}

func init() {
	proto.RegisterEnum("yandex.cloud.compute.v1.Snapshot_Status", Snapshot_Status_name, Snapshot_Status_value)
	proto.RegisterType((*Snapshot)(nil), "yandex.cloud.compute.v1.Snapshot")
	proto.RegisterMapType((map[string]string)(nil), "yandex.cloud.compute.v1.Snapshot.LabelsEntry")
}

func init() {
	proto.RegisterFile("yandex/cloud/compute/v1/snapshot.proto", fileDescriptor_cac027a1005d8550)
}

var fileDescriptor_cac027a1005d8550 = []byte{
	// 492 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0x41, 0x8f, 0x9b, 0x3e,
	0x10, 0xc5, 0xff, 0x84, 0x24, 0xff, 0x30, 0x44, 0x51, 0x64, 0x55, 0x2d, 0xca, 0x1e, 0x96, 0xae,
	0xaa, 0x8a, 0x4b, 0x40, 0x9b, 0x5e, 0xba, 0xed, 0xa5, 0x69, 0x42, 0x2b, 0xa4, 0xd5, 0xb6, 0x32,
	0xd9, 0x43, 0x7b, 0x41, 0x0e, 0xf6, 0xb2, 0x56, 0x08, 0x46, 0xd8, 0x44, 0xcd, 0x7e, 0xc1, 0x7e,
	0xad, 0x0a, 0x43, 0xa4, 0xf4, 0xb0, 0xea, 0x6d, 0x78, 0xf3, 0x1b, 0x3f, 0x3f, 0x3c, 0xf0, 0xf6,
	0x48, 0x0a, 0xca, 0x7e, 0x05, 0x69, 0x2e, 0x6a, 0x1a, 0xa4, 0x62, 0x5f, 0xd6, 0x8a, 0x05, 0x87,
	0xeb, 0x40, 0x16, 0xa4, 0x94, 0x8f, 0x42, 0xf9, 0x65, 0x25, 0x94, 0x40, 0xaf, 0x5a, 0xce, 0xd7,
	0x9c, 0xdf, 0x71, 0xfe, 0xe1, 0x7a, 0x76, 0x99, 0x09, 0x91, 0xe5, 0x2c, 0xd0, 0xd8, 0xb6, 0x7e,
	0x08, 0x14, 0xdf, 0x33, 0xa9, 0xc8, 0xbe, 0x6c, 0x27, 0xaf, 0x7e, 0xf7, 0x61, 0x14, 0x77, 0x87,
	0xa1, 0x09, 0xf4, 0x38, 0x75, 0x0c, 0xd7, 0xf0, 0x2c, 0xdc, 0xe3, 0x14, 0x5d, 0x80, 0xf5, 0x20,
	0x72, 0xca, 0xaa, 0x84, 0x53, 0xa7, 0xa7, 0xe5, 0x51, 0x2b, 0x44, 0x14, 0xdd, 0x00, 0xa4, 0x15,
	0x23, 0x8a, 0xd1, 0x84, 0x28, 0xc7, 0x74, 0x0d, 0xcf, 0x5e, 0xcc, 0xfc, 0xd6, 0xcf, 0x3f, 0xf9,
	0xf9, 0x9b, 0x93, 0x1f, 0xb6, 0x3a, 0x7a, 0xa9, 0x10, 0x82, 0x7e, 0x41, 0xf6, 0xcc, 0xe9, 0xeb,
	0x23, 0x75, 0x8d, 0x5c, 0xb0, 0x29, 0x93, 0x69, 0xc5, 0x4b, 0xc5, 0x45, 0xe1, 0x0c, 0x74, 0xeb,
	0x5c, 0x42, 0x21, 0x0c, 0x73, 0xb2, 0x65, 0xb9, 0x74, 0x86, 0xae, 0xe9, 0xd9, 0x8b, 0xb9, 0xff,
	0x4c, 0x6a, 0xff, 0x14, 0xc8, 0xbf, 0xd5, 0x7c, 0x58, 0xa8, 0xea, 0x88, 0xbb, 0x61, 0xf4, 0x1a,
	0xc6, 0x52, 0x89, 0x8a, 0x64, 0x2c, 0x91, 0xfc, 0x89, 0x39, 0xff, 0xbb, 0x86, 0x67, 0x62, 0xbb,
	0xd3, 0x62, 0xfe, 0xc4, 0x9a, 0xdc, 0x94, 0xcb, 0x5d, 0xdb, 0x1f, 0xe9, 0xfe, 0xa8, 0x11, 0x74,
	0xf3, 0x12, 0xec, 0xb2, 0x12, 0xb4, 0x4e, 0x55, 0xc2, 0xa9, 0x74, 0x2c, 0xd7, 0xf4, 0x2c, 0x0c,
	0x9d, 0x14, 0x51, 0x89, 0x3e, 0xc1, 0x50, 0x2a, 0xa2, 0x6a, 0xe9, 0x80, 0x6b, 0x78, 0x93, 0x85,
	0xf7, 0xef, 0x7b, 0xc6, 0x9a, 0xc7, 0xdd, 0x1c, 0x7a, 0x03, 0x13, 0x29, 0xea, 0x2a, 0x65, 0x89,
	0xbe, 0x06, 0xa7, 0x8e, 0xad, 0x7f, 0xc7, 0xb8, 0x55, 0xd7, 0x5c, 0xee, 0x22, 0x3a, 0xbb, 0x01,
	0xfb, 0x2c, 0x1f, 0x9a, 0x82, 0xb9, 0x63, 0xc7, 0xee, 0xf5, 0x9a, 0x12, 0xbd, 0x80, 0xc1, 0x81,
	0xe4, 0x35, 0xeb, 0x9e, 0xae, 0xfd, 0xf8, 0xd0, 0x7b, 0x6f, 0x5c, 0x61, 0x18, 0xb6, 0x96, 0xe8,
	0x25, 0xa0, 0x78, 0xb3, 0xdc, 0xdc, 0xc7, 0xc9, 0xfd, 0x5d, 0xfc, 0x3d, 0x5c, 0x45, 0x5f, 0xa2,
	0x70, 0x3d, 0xfd, 0x0f, 0x8d, 0x61, 0xb4, 0xc2, 0xe1, 0x72, 0x13, 0xdd, 0x7d, 0x9d, 0x1a, 0xc8,
	0x82, 0x01, 0x0e, 0x97, 0xeb, 0x1f, 0xd3, 0x5e, 0x53, 0x86, 0x18, 0x7f, 0xc3, 0x53, 0xb3, 0x61,
	0xd6, 0xe1, 0x6d, 0xa8, 0x99, 0xfe, 0xe7, 0x2d, 0x5c, 0xfc, 0x95, 0x93, 0x94, 0xfc, 0x2c, 0xeb,
	0xcf, 0x55, 0xc6, 0xd5, 0x63, 0xbd, 0x6d, 0xa4, 0xa0, 0xe5, 0xe6, 0xed, 0x56, 0x67, 0x62, 0x9e,
	0xb1, 0x42, 0x2f, 0x4c, 0xf0, 0xcc, 0xba, 0x7f, 0xec, 0xca, 0xed, 0x50, 0x63, 0xef, 0xfe, 0x04,
	0x00, 0x00, 0xff, 0xff, 0xf5, 0x2f, 0x5e, 0x1f, 0x18, 0x03, 0x00, 0x00,
}