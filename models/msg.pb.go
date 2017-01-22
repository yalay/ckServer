// Code generated by protoc-gen-go.
// source: msg.proto
// DO NOT EDIT!

/*
Package models is a generated protocol buffer package.

It is generated from these files:
	msg.proto

It has these top-level messages:
	Msg
*/
package models

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

type Msg struct {
	ArticleId int32 `protobuf:"varint,1,opt,name=articleId" json:"articleId,omitempty"`
	PkgIndex  int32 `protobuf:"varint,2,opt,name=pkgIndex" json:"pkgIndex,omitempty"`
}

func (m *Msg) Reset()                    { *m = Msg{} }
func (m *Msg) String() string            { return proto.CompactTextString(m) }
func (*Msg) ProtoMessage()               {}
func (*Msg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Msg) GetArticleId() int32 {
	if m != nil {
		return m.ArticleId
	}
	return 0
}

func (m *Msg) GetPkgIndex() int32 {
	if m != nil {
		return m.PkgIndex
	}
	return 0
}

func init() {
	proto.RegisterType((*Msg)(nil), "models.Msg")
}

func init() { proto.RegisterFile("msg.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 95 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcc, 0x2d, 0x4e, 0xd7,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcb, 0xcd, 0x4f, 0x49, 0xcd, 0x29, 0x56, 0xb2, 0xe7,
	0x62, 0xf6, 0x2d, 0x4e, 0x17, 0x92, 0xe1, 0xe2, 0x4c, 0x2c, 0x2a, 0xc9, 0x4c, 0xce, 0x49, 0xf5,
	0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0x60, 0x0d, 0x42, 0x08, 0x08, 0x49, 0x71, 0x71, 0x14, 0x64,
	0xa7, 0x7b, 0xe6, 0xa5, 0xa4, 0x56, 0x48, 0x30, 0x81, 0x25, 0xe1, 0xfc, 0x24, 0x36, 0xb0, 0x79,
	0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8c, 0xa8, 0x53, 0xbc, 0x5c, 0x00, 0x00, 0x00,
}