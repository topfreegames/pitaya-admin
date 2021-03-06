// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protoD.proto

package protos

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

type LastMessage struct {
	Last                 string    `protobuf:"bytes,1,opt,name=last,proto3" json:"last,omitempty"`
	Final                *MessageE `protobuf:"bytes,2,opt,name=final,proto3" json:"final,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *LastMessage) Reset()         { *m = LastMessage{} }
func (m *LastMessage) String() string { return proto.CompactTextString(m) }
func (*LastMessage) ProtoMessage()    {}
func (*LastMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_protoD_cba76a45f743b8ef, []int{0}
}
func (m *LastMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LastMessage.Unmarshal(m, b)
}
func (m *LastMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LastMessage.Marshal(b, m, deterministic)
}
func (dst *LastMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LastMessage.Merge(dst, src)
}
func (m *LastMessage) XXX_Size() int {
	return xxx_messageInfo_LastMessage.Size(m)
}
func (m *LastMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_LastMessage.DiscardUnknown(m)
}

var xxx_messageInfo_LastMessage proto.InternalMessageInfo

func (m *LastMessage) GetLast() string {
	if m != nil {
		return m.Last
	}
	return ""
}

func (m *LastMessage) GetFinal() *MessageE {
	if m != nil {
		return m.Final
	}
	return nil
}

func init() {
	proto.RegisterType((*LastMessage)(nil), "protos.LastMessage")
}

func init() { proto.RegisterFile("protoD.proto", fileDescriptor_protoD_cba76a45f743b8ef) }

var fileDescriptor_protoD_cba76a45f743b8ef = []byte{
	// 108 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x28, 0xca, 0x2f,
	0xc9, 0x77, 0xd1, 0x03, 0x53, 0x42, 0x6c, 0x60, 0xaa, 0x58, 0x0a, 0x22, 0xea, 0x0a, 0x11, 0x55,
	0xf2, 0xe4, 0xe2, 0xf6, 0x49, 0x2c, 0x2e, 0xf1, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x15, 0x12,
	0xe2, 0x62, 0xc9, 0x49, 0x2c, 0x2e, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x85,
	0xd4, 0xb8, 0x58, 0xd3, 0x32, 0xf3, 0x12, 0x73, 0x24, 0x98, 0x14, 0x18, 0x35, 0xb8, 0x8d, 0x04,
	0x20, 0x3a, 0x8b, 0xf5, 0xa0, 0x7a, 0x5c, 0x83, 0x20, 0xd2, 0x49, 0x10, 0x0b, 0x8c, 0x01, 0x01,
	0x00, 0x00, 0xff, 0xff, 0xd4, 0xf3, 0x63, 0x89, 0x77, 0x00, 0x00, 0x00,
}
