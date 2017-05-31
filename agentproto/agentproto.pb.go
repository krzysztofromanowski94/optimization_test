// Code generated by protoc-gen-go. DO NOT EDIT.
// source: agentproto.proto

/*
Package agentproto is a generated protocol buffer package.

It is generated from these files:
	agentproto.proto

It has these top-level messages:
	AgentData
*/
package agentproto

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

type AgentData struct {
	X                *float64 `protobuf:"fixed64,1,req,name=x" json:"x,omitempty"`
	Y                *float64 `protobuf:"fixed64,2,req,name=y" json:"y,omitempty"`
	Fitness          *float64 `protobuf:"fixed64,3,req,name=fitness" json:"fitness,omitempty"`
	Average          *int32   `protobuf:"varint,4,req,name=average" json:"average,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *AgentData) Reset()                    { *m = AgentData{} }
func (m *AgentData) String() string            { return proto.CompactTextString(m) }
func (*AgentData) ProtoMessage()               {}
func (*AgentData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *AgentData) GetX() float64 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *AgentData) GetY() float64 {
	if m != nil && m.Y != nil {
		return *m.Y
	}
	return 0
}

func (m *AgentData) GetFitness() float64 {
	if m != nil && m.Fitness != nil {
		return *m.Fitness
	}
	return 0
}

func (m *AgentData) GetAverage() int32 {
	if m != nil && m.Average != nil {
		return *m.Average
	}
	return 0
}

func init() {
	proto.RegisterType((*AgentData)(nil), "agentproto.AgentData")
}

func init() { proto.RegisterFile("agentproto.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 96 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x48, 0x4c, 0x4f, 0xcd,
	0x2b, 0x29, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x03, 0x93, 0x42, 0x5c, 0x08, 0x11, 0x25, 0x67, 0x2e,
	0x4e, 0x47, 0x10, 0xcf, 0x25, 0xb1, 0x24, 0x51, 0x88, 0x93, 0x8b, 0xb1, 0x42, 0x82, 0x51, 0x81,
	0x49, 0x83, 0x11, 0xc4, 0xac, 0x94, 0x60, 0x02, 0x33, 0xf9, 0xb9, 0xd8, 0xd3, 0x32, 0x4b, 0xf2,
	0x52, 0x8b, 0x8b, 0x25, 0x98, 0x61, 0x02, 0x89, 0x65, 0xa9, 0x45, 0x40, 0x53, 0x24, 0x58, 0x80,
	0x02, 0xac, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x0f, 0x23, 0x59, 0x63, 0x00, 0x00, 0x00,
}