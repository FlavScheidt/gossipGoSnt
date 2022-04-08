// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gossip_message.proto

package proto

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

type Gossip struct {
	Message              []byte   `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
	Validator_Key        []byte   `protobuf:"bytes,2,opt,name=Validator_Key,json=ValidatorKey,proto3" json:"Validator_Key,omitempty"`
	Hash                 string   `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Gossip) Reset()         { *m = Gossip{} }
func (m *Gossip) String() string { return proto.CompactTextString(m) }
func (*Gossip) ProtoMessage()    {}
func (*Gossip) Descriptor() ([]byte, []int) {
	return fileDescriptor_288e7560062f96b3, []int{0}
}

func (m *Gossip) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Gossip.Unmarshal(m, b)
}
func (m *Gossip) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Gossip.Marshal(b, m, deterministic)
}
func (m *Gossip) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Gossip.Merge(m, src)
}
func (m *Gossip) XXX_Size() int {
	return xxx_messageInfo_Gossip.Size(m)
}
func (m *Gossip) XXX_DiscardUnknown() {
	xxx_messageInfo_Gossip.DiscardUnknown(m)
}

var xxx_messageInfo_Gossip proto.InternalMessageInfo

func (m *Gossip) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *Gossip) GetValidator_Key() []byte {
	if m != nil {
		return m.Validator_Key
	}
	return nil
}

func (m *Gossip) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

type Control struct {
	Stream               bool     `protobuf:"varint,1,opt,name=stream,proto3" json:"stream,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Control) Reset()         { *m = Control{} }
func (m *Control) String() string { return proto.CompactTextString(m) }
func (*Control) ProtoMessage()    {}
func (*Control) Descriptor() ([]byte, []int) {
	return fileDescriptor_288e7560062f96b3, []int{1}
}

func (m *Control) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Control.Unmarshal(m, b)
}
func (m *Control) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Control.Marshal(b, m, deterministic)
}
func (m *Control) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Control.Merge(m, src)
}
func (m *Control) XXX_Size() int {
	return xxx_messageInfo_Control.Size(m)
}
func (m *Control) XXX_DiscardUnknown() {
	xxx_messageInfo_Control.DiscardUnknown(m)
}

var xxx_messageInfo_Control proto.InternalMessageInfo

func (m *Control) GetStream() bool {
	if m != nil {
		return m.Stream
	}
	return false
}

func init() {
	proto.RegisterType((*Gossip)(nil), "gossipmessage.Gossip")
	proto.RegisterType((*Control)(nil), "gossipmessage.Control")
}

func init() { proto.RegisterFile("gossip_message.proto", fileDescriptor_288e7560062f96b3) }

var fileDescriptor_288e7560062f96b3 = []byte{
	// 239 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x49, 0xcf, 0x2f, 0x2e,
	0xce, 0x2c, 0x88, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0xe2, 0x85, 0x88, 0x42, 0x05, 0x95, 0xa2, 0xb9, 0xd8, 0xdc, 0xc1, 0x02, 0x42, 0x12, 0x5c,
	0xec, 0xbe, 0x10, 0x41, 0x09, 0x46, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x18, 0x57, 0x48, 0x99, 0x8b,
	0x37, 0x2c, 0x31, 0x27, 0x33, 0x25, 0xb1, 0x24, 0xbf, 0x28, 0xde, 0x3b, 0xb5, 0x52, 0x82, 0x09,
	0x2c, 0xcf, 0x03, 0x17, 0xf4, 0x4e, 0xad, 0x14, 0x12, 0xe2, 0x62, 0xc9, 0x48, 0x2c, 0xce, 0x90,
	0x60, 0x56, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x95, 0x14, 0xb9, 0xd8, 0x9d, 0xf3, 0xf3, 0x4a,
	0x8a, 0xf2, 0x73, 0x84, 0xc4, 0xb8, 0xd8, 0x8a, 0x4b, 0x8a, 0x52, 0x13, 0x73, 0xc1, 0x86, 0x73,
	0x04, 0x41, 0x79, 0x46, 0x5d, 0x8c, 0x5c, 0xbc, 0x10, 0x07, 0xc0, 0x6c, 0xb3, 0xe6, 0xe2, 0x28,
	0xc9, 0xf7, 0xc9, 0x4c, 0x0a, 0x30, 0x0a, 0x10, 0x12, 0xd5, 0x43, 0x71, 0xad, 0x1e, 0x44, 0xa5,
	0x94, 0x18, 0x9a, 0x30, 0xd4, 0x12, 0x25, 0x06, 0x21, 0x1b, 0x2e, 0xce, 0x92, 0xfc, 0xa0, 0xcc,
	0x82, 0x82, 0x9c, 0xd4, 0x14, 0x92, 0x75, 0x3b, 0x69, 0x45, 0x69, 0xa4, 0x67, 0x96, 0x64, 0x94,
	0x26, 0xe9, 0x25, 0xe7, 0xe7, 0xea, 0xbb, 0xe5, 0x24, 0x96, 0x05, 0x27, 0x67, 0xa4, 0x66, 0xa6,
	0x94, 0xe8, 0x43, 0x74, 0xb8, 0xe7, 0x07, 0xe7, 0x95, 0xe8, 0x83, 0xc3, 0x31, 0x89, 0x0d, 0x4c,
	0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x2e, 0x00, 0xe3, 0x74, 0x66, 0x01, 0x00, 0x00,
}