// Code generated by protoc-gen-go. DO NOT EDIT.
// source: yandex/cloud/compute/v1/zone_service.proto

package compute

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/yandex-cloud/go-genproto/yandex/cloud"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type ListZonesRequest struct {
	// The maximum number of results per page to return. If the number of available
	// results is larger than [page_size],
	// the service returns a [ListZonesResponse.next_page_token]
	// that can be used to get the next page of results in subsequent list requests.
	PageSize int64 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// Page token. To get the next page of results, set [page_token] to the
	// [ListZonesResponse.next_page_token] returned by a previous list request.
	PageToken            string   `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListZonesRequest) Reset()         { *m = ListZonesRequest{} }
func (m *ListZonesRequest) String() string { return proto.CompactTextString(m) }
func (*ListZonesRequest) ProtoMessage()    {}
func (*ListZonesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c45093f3209cbc9e, []int{0}
}

func (m *ListZonesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListZonesRequest.Unmarshal(m, b)
}
func (m *ListZonesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListZonesRequest.Marshal(b, m, deterministic)
}
func (m *ListZonesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListZonesRequest.Merge(m, src)
}
func (m *ListZonesRequest) XXX_Size() int {
	return xxx_messageInfo_ListZonesRequest.Size(m)
}
func (m *ListZonesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListZonesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListZonesRequest proto.InternalMessageInfo

func (m *ListZonesRequest) GetPageSize() int64 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ListZonesRequest) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

type ListZonesResponse struct {
	// List of availability zones.
	Zones []*Zone `protobuf:"bytes,1,rep,name=zones,proto3" json:"zones,omitempty"`
	// This token allows you to get the next page of results for list requests. If the number of results
	// is larger than [ListZonesRequest.page_size], use
	// the [ListZonesRequest.page_token] as the value
	// for the [ListZonesRequest.page_token] query parameter
	// in the next list request. Subsequent list requests will have their own
	// [ListZonesRequest.page_token] to continue paging through the results.
	NextPageToken        string   `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListZonesResponse) Reset()         { *m = ListZonesResponse{} }
func (m *ListZonesResponse) String() string { return proto.CompactTextString(m) }
func (*ListZonesResponse) ProtoMessage()    {}
func (*ListZonesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c45093f3209cbc9e, []int{1}
}

func (m *ListZonesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListZonesResponse.Unmarshal(m, b)
}
func (m *ListZonesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListZonesResponse.Marshal(b, m, deterministic)
}
func (m *ListZonesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListZonesResponse.Merge(m, src)
}
func (m *ListZonesResponse) XXX_Size() int {
	return xxx_messageInfo_ListZonesResponse.Size(m)
}
func (m *ListZonesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListZonesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListZonesResponse proto.InternalMessageInfo

func (m *ListZonesResponse) GetZones() []*Zone {
	if m != nil {
		return m.Zones
	}
	return nil
}

func (m *ListZonesResponse) GetNextPageToken() string {
	if m != nil {
		return m.NextPageToken
	}
	return ""
}

type GetZoneRequest struct {
	// ID of the availability zone to return information about.
	ZoneId               string   `protobuf:"bytes,1,opt,name=zone_id,json=zoneId,proto3" json:"zone_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetZoneRequest) Reset()         { *m = GetZoneRequest{} }
func (m *GetZoneRequest) String() string { return proto.CompactTextString(m) }
func (*GetZoneRequest) ProtoMessage()    {}
func (*GetZoneRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c45093f3209cbc9e, []int{2}
}

func (m *GetZoneRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetZoneRequest.Unmarshal(m, b)
}
func (m *GetZoneRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetZoneRequest.Marshal(b, m, deterministic)
}
func (m *GetZoneRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetZoneRequest.Merge(m, src)
}
func (m *GetZoneRequest) XXX_Size() int {
	return xxx_messageInfo_GetZoneRequest.Size(m)
}
func (m *GetZoneRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetZoneRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetZoneRequest proto.InternalMessageInfo

func (m *GetZoneRequest) GetZoneId() string {
	if m != nil {
		return m.ZoneId
	}
	return ""
}

func init() {
	proto.RegisterType((*ListZonesRequest)(nil), "yandex.cloud.compute.v1.ListZonesRequest")
	proto.RegisterType((*ListZonesResponse)(nil), "yandex.cloud.compute.v1.ListZonesResponse")
	proto.RegisterType((*GetZoneRequest)(nil), "yandex.cloud.compute.v1.GetZoneRequest")
}

func init() {
	proto.RegisterFile("yandex/cloud/compute/v1/zone_service.proto", fileDescriptor_c45093f3209cbc9e)
}

var fileDescriptor_c45093f3209cbc9e = []byte{
	// 428 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xbf, 0x8f, 0xd3, 0x30,
	0x14, 0xc7, 0x95, 0xeb, 0x5d, 0x21, 0x3e, 0x7e, 0x9d, 0x19, 0x28, 0x39, 0x2a, 0x55, 0x41, 0x70,
	0xe1, 0xa4, 0x8b, 0x93, 0x3b, 0x21, 0x06, 0xda, 0x25, 0x0c, 0x15, 0x12, 0x03, 0x4a, 0x99, 0xba,
	0x54, 0x69, 0xf3, 0x14, 0x2c, 0x8a, 0x1d, 0x6a, 0x27, 0x2a, 0x45, 0x2c, 0x8c, 0x5d, 0xf9, 0xa3,
	0xda, 0x9d, 0x7f, 0x81, 0x81, 0xbf, 0x01, 0x16, 0x64, 0x3b, 0xa0, 0xb6, 0x28, 0x15, 0x9b, 0xe5,
	0xf7, 0x79, 0xef, 0xfb, 0x7d, 0x3f, 0xd0, 0xf9, 0xc7, 0x84, 0xa5, 0x30, 0x27, 0x93, 0x29, 0x2f,
	0x52, 0x32, 0xe1, 0xef, 0xf3, 0x42, 0x02, 0x29, 0x43, 0xb2, 0xe0, 0x0c, 0x46, 0x02, 0x66, 0x25,
	0x9d, 0x80, 0x9f, 0xcf, 0xb8, 0xe4, 0xf8, 0x9e, 0x61, 0x7d, 0xcd, 0xfa, 0x15, 0xeb, 0x97, 0xa1,
	0xf3, 0x20, 0xe3, 0x3c, 0x9b, 0x02, 0x49, 0x72, 0x4a, 0x12, 0xc6, 0xb8, 0x4c, 0x24, 0xe5, 0x4c,
	0x98, 0x34, 0xc7, 0xdd, 0x27, 0x51, 0x31, 0xed, 0x2d, 0xa6, 0x4c, 0xa6, 0x34, 0xd5, 0x35, 0x4c,
	0xd8, 0x05, 0x74, 0xe7, 0x15, 0x15, 0x72, 0xc8, 0x19, 0x88, 0x18, 0x3e, 0x14, 0x20, 0x24, 0x3e,
	0x43, 0x76, 0x9e, 0x64, 0x30, 0x12, 0x74, 0x01, 0x2d, 0xab, 0x63, 0x79, 0x8d, 0x08, 0xfd, 0x5c,
	0x85, 0xcd, 0x6e, 0x2f, 0x0c, 0x82, 0x20, 0xbe, 0xae, 0x82, 0x03, 0xba, 0x00, 0xec, 0x21, 0xa4,
	0x41, 0xc9, 0xdf, 0x01, 0x6b, 0x1d, 0x74, 0x2c, 0xcf, 0x8e, 0xec, 0xe5, 0x3a, 0x3c, 0xd2, 0x64,
	0xac, 0xab, 0xbc, 0x51, 0x31, 0x37, 0x47, 0x27, 0x1b, 0x32, 0x22, 0xe7, 0x4c, 0x00, 0xbe, 0x42,
	0x47, 0xca, 0xa8, 0x68, 0x59, 0x9d, 0x86, 0x77, 0x7c, 0xd9, 0xf6, 0x6b, 0xa6, 0xe0, 0xab, 0xb4,
	0xd8, 0xb0, 0xf8, 0x31, 0xba, 0xcd, 0x60, 0x2e, 0x47, 0xbb, 0xc2, 0xf1, 0x4d, 0xf5, 0xfd, 0xfa,
	0xaf, 0xe2, 0x33, 0x74, 0xab, 0x0f, 0x5a, 0xf0, 0x4f, 0x5b, 0x8f, 0xd0, 0x35, 0x3d, 0x7a, 0x9a,
	0xea, 0xa6, 0xec, 0xe8, 0xc6, 0x8f, 0x55, 0x68, 0x2d, 0xd7, 0xe1, 0x61, 0xb7, 0xf7, 0x34, 0x88,
	0x9b, 0x2a, 0xf8, 0x32, 0xbd, 0xfc, 0x65, 0xa1, 0x63, 0x95, 0x36, 0x30, 0x1b, 0xc2, 0x33, 0xd4,
	0xe8, 0x83, 0xc4, 0x67, 0xb5, 0xee, 0xb6, 0x65, 0x9c, 0xfd, 0x6d, 0xb8, 0x0f, 0xbf, 0x7c, 0xfb,
	0xfe, 0xf5, 0xa0, 0x8d, 0x4f, 0x77, 0xf7, 0x25, 0xc8, 0xa7, 0xca, 0xde, 0x67, 0x3c, 0x47, 0x87,
	0x6a, 0x5c, 0xf8, 0x49, 0x6d, 0xad, 0xdd, 0xa5, 0x39, 0xe7, 0xff, 0x83, 0x9a, 0xc1, 0xbb, 0xf7,
	0xb5, 0x87, 0xbb, 0xf8, 0xe4, 0x1f, 0x0f, 0xd1, 0x18, 0x9d, 0x6e, 0xd5, 0x49, 0x72, 0xba, 0x51,
	0x6b, 0xf8, 0x22, 0xa3, 0xf2, 0x6d, 0x31, 0x56, 0x5f, 0xc4, 0x70, 0x17, 0xe6, 0xb0, 0x32, 0x7e,
	0x91, 0x01, 0xd3, 0x37, 0x45, 0x6a, 0xae, 0xf2, 0x79, 0xf5, 0x1c, 0x37, 0x35, 0x76, 0xf5, 0x3b,
	0x00, 0x00, 0xff, 0xff, 0x41, 0xfa, 0xcd, 0x99, 0x22, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ZoneServiceClient is the client API for ZoneService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ZoneServiceClient interface {
	// Returns the information about the specified availability zone.
	//
	// To get the list of availability zones, make a [List] request.
	Get(ctx context.Context, in *GetZoneRequest, opts ...grpc.CallOption) (*Zone, error)
	// Retrieves the list of availability zones.
	List(ctx context.Context, in *ListZonesRequest, opts ...grpc.CallOption) (*ListZonesResponse, error)
}

type zoneServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewZoneServiceClient(cc grpc.ClientConnInterface) ZoneServiceClient {
	return &zoneServiceClient{cc}
}

func (c *zoneServiceClient) Get(ctx context.Context, in *GetZoneRequest, opts ...grpc.CallOption) (*Zone, error) {
	out := new(Zone)
	err := c.cc.Invoke(ctx, "/yandex.cloud.compute.v1.ZoneService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zoneServiceClient) List(ctx context.Context, in *ListZonesRequest, opts ...grpc.CallOption) (*ListZonesResponse, error) {
	out := new(ListZonesResponse)
	err := c.cc.Invoke(ctx, "/yandex.cloud.compute.v1.ZoneService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ZoneServiceServer is the server API for ZoneService service.
type ZoneServiceServer interface {
	// Returns the information about the specified availability zone.
	//
	// To get the list of availability zones, make a [List] request.
	Get(context.Context, *GetZoneRequest) (*Zone, error)
	// Retrieves the list of availability zones.
	List(context.Context, *ListZonesRequest) (*ListZonesResponse, error)
}

// UnimplementedZoneServiceServer can be embedded to have forward compatible implementations.
type UnimplementedZoneServiceServer struct {
}

func (*UnimplementedZoneServiceServer) Get(ctx context.Context, req *GetZoneRequest) (*Zone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedZoneServiceServer) List(ctx context.Context, req *ListZonesRequest) (*ListZonesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

func RegisterZoneServiceServer(s *grpc.Server, srv ZoneServiceServer) {
	s.RegisterService(&_ZoneService_serviceDesc, srv)
}

func _ZoneService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetZoneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZoneServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/yandex.cloud.compute.v1.ZoneService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZoneServiceServer).Get(ctx, req.(*GetZoneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ZoneService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListZonesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZoneServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/yandex.cloud.compute.v1.ZoneService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZoneServiceServer).List(ctx, req.(*ListZonesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ZoneService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "yandex.cloud.compute.v1.ZoneService",
	HandlerType: (*ZoneServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _ZoneService_Get_Handler,
		},
		{
			MethodName: "List",
			Handler:    _ZoneService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "yandex/cloud/compute/v1/zone_service.proto",
}