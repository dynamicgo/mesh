// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mesh.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type LockupRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LockupRequest) Reset()         { *m = LockupRequest{} }
func (m *LockupRequest) String() string { return proto.CompactTextString(m) }
func (*LockupRequest) ProtoMessage()    {}
func (*LockupRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mesh_db396b7f2924b68e, []int{0}
}
func (m *LockupRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LockupRequest.Unmarshal(m, b)
}
func (m *LockupRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LockupRequest.Marshal(b, m, deterministic)
}
func (dst *LockupRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LockupRequest.Merge(dst, src)
}
func (m *LockupRequest) XXX_Size() int {
	return xxx_messageInfo_LockupRequest.Size(m)
}
func (m *LockupRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LockupRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LockupRequest proto.InternalMessageInfo

func (m *LockupRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type LookupResponse struct {
	Addrs                []string `protobuf:"bytes,1,rep,name=addrs" json:"addrs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LookupResponse) Reset()         { *m = LookupResponse{} }
func (m *LookupResponse) String() string { return proto.CompactTextString(m) }
func (*LookupResponse) ProtoMessage()    {}
func (*LookupResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_mesh_db396b7f2924b68e, []int{1}
}
func (m *LookupResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LookupResponse.Unmarshal(m, b)
}
func (m *LookupResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LookupResponse.Marshal(b, m, deterministic)
}
func (dst *LookupResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LookupResponse.Merge(dst, src)
}
func (m *LookupResponse) XXX_Size() int {
	return xxx_messageInfo_LookupResponse.Size(m)
}
func (m *LookupResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LookupResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LookupResponse proto.InternalMessageInfo

func (m *LookupResponse) GetAddrs() []string {
	if m != nil {
		return m.Addrs
	}
	return nil
}

type RegisterRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Addrs                []string `protobuf:"bytes,2,rep,name=addrs" json:"addrs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterRequest) Reset()         { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterRequest) ProtoMessage()    {}
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mesh_db396b7f2924b68e, []int{2}
}
func (m *RegisterRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRequest.Unmarshal(m, b)
}
func (m *RegisterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRequest.Marshal(b, m, deterministic)
}
func (dst *RegisterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRequest.Merge(dst, src)
}
func (m *RegisterRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterRequest.Size(m)
}
func (m *RegisterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRequest proto.InternalMessageInfo

func (m *RegisterRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RegisterRequest) GetAddrs() []string {
	if m != nil {
		return m.Addrs
	}
	return nil
}

type RegisterResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterResponse) Reset()         { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterResponse) ProtoMessage()    {}
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_mesh_db396b7f2924b68e, []int{3}
}
func (m *RegisterResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterResponse.Unmarshal(m, b)
}
func (m *RegisterResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterResponse.Marshal(b, m, deterministic)
}
func (dst *RegisterResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterResponse.Merge(dst, src)
}
func (m *RegisterResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterResponse.Size(m)
}
func (m *RegisterResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*LockupRequest)(nil), "proto.LockupRequest")
	proto.RegisterType((*LookupResponse)(nil), "proto.LookupResponse")
	proto.RegisterType((*RegisterRequest)(nil), "proto.RegisterRequest")
	proto.RegisterType((*RegisterResponse)(nil), "proto.RegisterResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ServiceHub service

type ServiceHubClient interface {
	Lookup(ctx context.Context, in *LockupRequest, opts ...grpc.CallOption) (*LookupResponse, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
}

type serviceHubClient struct {
	cc *grpc.ClientConn
}

func NewServiceHubClient(cc *grpc.ClientConn) ServiceHubClient {
	return &serviceHubClient{cc}
}

func (c *serviceHubClient) Lookup(ctx context.Context, in *LockupRequest, opts ...grpc.CallOption) (*LookupResponse, error) {
	out := new(LookupResponse)
	err := grpc.Invoke(ctx, "/proto.ServiceHub/Lookup", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceHubClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := grpc.Invoke(ctx, "/proto.ServiceHub/Register", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ServiceHub service

type ServiceHubServer interface {
	Lookup(context.Context, *LockupRequest) (*LookupResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
}

func RegisterServiceHubServer(s *grpc.Server, srv ServiceHubServer) {
	s.RegisterService(&_ServiceHub_serviceDesc, srv)
}

func _ServiceHub_Lookup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LockupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceHubServer).Lookup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ServiceHub/Lookup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceHubServer).Lookup(ctx, req.(*LockupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceHub_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceHubServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ServiceHub/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceHubServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ServiceHub_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ServiceHub",
	HandlerType: (*ServiceHubServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Lookup",
			Handler:    _ServiceHub_Lookup_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _ServiceHub_Register_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mesh.proto",
}

func init() { proto.RegisterFile("mesh.proto", fileDescriptor_mesh_db396b7f2924b68e) }

var fileDescriptor_mesh_db396b7f2924b68e = []byte{
	// 193 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xca, 0x4d, 0x2d, 0xce,
	0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0xca, 0x5c, 0xbc, 0x3e, 0xf9,
	0xc9, 0xd9, 0xa5, 0x05, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x42, 0x5c, 0x2c, 0x79,
	0x89, 0xb9, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x92, 0x1a, 0x17, 0x9f,
	0x4f, 0x7e, 0x3e, 0x58, 0x51, 0x71, 0x41, 0x7e, 0x5e, 0x71, 0xaa, 0x90, 0x08, 0x17, 0x6b, 0x62,
	0x4a, 0x4a, 0x51, 0xb1, 0x04, 0xa3, 0x02, 0xb3, 0x06, 0x67, 0x10, 0x84, 0xa3, 0x64, 0xcd, 0xc5,
	0x1f, 0x94, 0x9a, 0x9e, 0x59, 0x5c, 0x92, 0x5a, 0x84, 0xc7, 0x38, 0x84, 0x66, 0x26, 0x64, 0xcd,
	0x42, 0x5c, 0x02, 0x08, 0xcd, 0x10, 0x6b, 0x8c, 0x5a, 0x18, 0xb9, 0xb8, 0x82, 0x53, 0x8b, 0xca,
	0x32, 0x93, 0x53, 0x3d, 0x4a, 0x93, 0x84, 0xcc, 0xb9, 0xd8, 0x20, 0xee, 0x10, 0x12, 0x81, 0xf8,
	0x42, 0x0f, 0xc5, 0xed, 0x52, 0xa2, 0x70, 0x51, 0x64, 0xc7, 0x2a, 0x31, 0x08, 0xd9, 0x72, 0x71,
	0xc0, 0xcc, 0x16, 0x12, 0x83, 0x2a, 0x42, 0x73, 0xa9, 0x94, 0x38, 0x86, 0x38, 0x4c, 0x7b, 0x12,
	0x1b, 0x58, 0xc6, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x18, 0x4d, 0xac, 0x02, 0x40, 0x01, 0x00,
	0x00,
}
