// Code generated by protoc-gen-go. DO NOT EDIT.
// source: transport.proto

/*
Package grpc is a generated protocol buffer package.

It is generated from these files:
	transport.proto

It has these top-level messages:
	Request
	Response
*/
package grpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Request struct {
	ServiceMethod    *string `protobuf:"bytes,1,opt,name=service_method,json=serviceMethod" json:"service_method,omitempty"`
	Seq              *uint64 `protobuf:"fixed64,2,opt,name=seq" json:"seq,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetServiceMethod() string {
	if m != nil && m.ServiceMethod != nil {
		return *m.ServiceMethod
	}
	return ""
}

func (m *Request) GetSeq() uint64 {
	if m != nil && m.Seq != nil {
		return *m.Seq
	}
	return 0
}

type Response struct {
	ServiceMethod    *string `protobuf:"bytes,1,opt,name=service_method,json=serviceMethod" json:"service_method,omitempty"`
	Seq              *uint64 `protobuf:"fixed64,2,opt,name=seq" json:"seq,omitempty"`
	Error            *string `protobuf:"bytes,3,opt,name=error" json:"error,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Response) GetServiceMethod() string {
	if m != nil && m.ServiceMethod != nil {
		return *m.ServiceMethod
	}
	return ""
}

func (m *Response) GetSeq() uint64 {
	if m != nil && m.Seq != nil {
		return *m.Seq
	}
	return 0
}

func (m *Response) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "grpc.Request")
	proto.RegisterType((*Response)(nil), "grpc.Response")
}

func init() { proto.RegisterFile("transport.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 134 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0x29, 0x4a, 0xcc,
	0x2b, 0x2e, 0xc8, 0x2f, 0x2a, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x49, 0x2f, 0x2a,
	0x48, 0x56, 0x72, 0xe2, 0x62, 0x0f, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x52, 0xe5, 0xe2,
	0x2b, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0x8d, 0xcf, 0x4d, 0x2d, 0xc9, 0xc8, 0x4f, 0x91, 0x60,
	0x54, 0x60, 0xd4, 0xe0, 0x0c, 0xe2, 0x85, 0x8a, 0xfa, 0x82, 0x05, 0x85, 0x04, 0xb8, 0x98, 0x8b,
	0x53, 0x0b, 0x25, 0x98, 0x14, 0x18, 0x35, 0xd8, 0x82, 0x40, 0x4c, 0xa5, 0x48, 0x2e, 0x8e, 0xa0,
	0xd4, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0xb2, 0x0d, 0x11, 0x12, 0xe1, 0x62, 0x4d, 0x2d, 0x2a,
	0xca, 0x2f, 0x92, 0x60, 0x06, 0xab, 0x87, 0x70, 0x00, 0x01, 0x00, 0x00, 0xff, 0xff, 0x77, 0x12,
	0xe1, 0x1e, 0xb6, 0x00, 0x00, 0x00,
}