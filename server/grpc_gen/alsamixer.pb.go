// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.13.0
// source: alsamixer.proto

package grpc_gen

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Card    string  `protobuf:"bytes,1,opt,name=card,proto3" json:"card,omitempty"`
	Control string  `protobuf:"bytes,2,opt,name=control,proto3" json:"control,omitempty"`
	Volume  []int32 `protobuf:"varint,3,rep,packed,name=volume,proto3" json:"volume,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alsamixer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_alsamixer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_alsamixer_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetCard() string {
	if x != nil {
		return x.Card
	}
	return ""
}

func (x *Request) GetControl() string {
	if x != nil {
		return x.Control
	}
	return ""
}

func (x *Request) GetVolume() []int32 {
	if x != nil {
		return x.Volume
	}
	return nil
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Card     string              `protobuf:"bytes,1,opt,name=card,proto3" json:"card,omitempty"`
	Controls []*Response_Control `protobuf:"bytes,2,rep,name=controls,proto3" json:"controls,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alsamixer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_alsamixer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_alsamixer_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetCard() string {
	if x != nil {
		return x.Card
	}
	return ""
}

func (x *Response) GetControls() []*Response_Control {
	if x != nil {
		return x.Controls
	}
	return nil
}

type Response_Control struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Volume []int32 `protobuf:"varint,2,rep,packed,name=volume,proto3" json:"volume,omitempty"`
}

func (x *Response_Control) Reset() {
	*x = Response_Control{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alsamixer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response_Control) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response_Control) ProtoMessage() {}

func (x *Response_Control) ProtoReflect() protoreflect.Message {
	mi := &file_alsamixer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response_Control.ProtoReflect.Descriptor instead.
func (*Response_Control) Descriptor() ([]byte, []int) {
	return file_alsamixer_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Response_Control) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Response_Control) GetVolume() []int32 {
	if x != nil {
		return x.Volume
	}
	return nil
}

var File_alsamixer_proto protoreflect.FileDescriptor

var file_alsamixer_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x61, 0x6c, 0x73, 0x61, 0x6d, 0x69, 0x78, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x61, 0x6c, 0x73, 0x61, 0x22, 0x4f, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x61, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x63, 0x61, 0x72, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x12, 0x16, 0x0a, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x03, 0x28, 0x05,
	0x52, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x22, 0x89, 0x01, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x61, 0x72, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x61, 0x72, 0x64, 0x12, 0x32, 0x0a, 0x08, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x6c,
	0x73, 0x61, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x73, 0x1a, 0x35, 0x0a,
	0x07, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x05, 0x52, 0x06, 0x76, 0x6f,
	0x6c, 0x75, 0x6d, 0x65, 0x32, 0x3f, 0x0a, 0x09, 0x41, 0x6c, 0x73, 0x61, 0x6d, 0x69, 0x78, 0x65,
	0x72, 0x12, 0x32, 0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x12, 0x0d, 0x2e, 0x61, 0x6c, 0x73, 0x61, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x0e, 0x2e, 0x61, 0x6c, 0x73, 0x61, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x5f,
	0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_alsamixer_proto_rawDescOnce sync.Once
	file_alsamixer_proto_rawDescData = file_alsamixer_proto_rawDesc
)

func file_alsamixer_proto_rawDescGZIP() []byte {
	file_alsamixer_proto_rawDescOnce.Do(func() {
		file_alsamixer_proto_rawDescData = protoimpl.X.CompressGZIP(file_alsamixer_proto_rawDescData)
	})
	return file_alsamixer_proto_rawDescData
}

var file_alsamixer_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_alsamixer_proto_goTypes = []interface{}{
	(*Request)(nil),          // 0: alsa.Request
	(*Response)(nil),         // 1: alsa.Response
	(*Response_Control)(nil), // 2: alsa.Response.Control
}
var file_alsamixer_proto_depIdxs = []int32{
	2, // 0: alsa.Response.controls:type_name -> alsa.Response.Control
	0, // 1: alsa.Alsamixer.Communicate:input_type -> alsa.Request
	1, // 2: alsa.Alsamixer.Communicate:output_type -> alsa.Response
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_alsamixer_proto_init() }
func file_alsamixer_proto_init() {
	if File_alsamixer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_alsamixer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_alsamixer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_alsamixer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response_Control); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_alsamixer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_alsamixer_proto_goTypes,
		DependencyIndexes: file_alsamixer_proto_depIdxs,
		MessageInfos:      file_alsamixer_proto_msgTypes,
	}.Build()
	File_alsamixer_proto = out.File
	file_alsamixer_proto_rawDesc = nil
	file_alsamixer_proto_goTypes = nil
	file_alsamixer_proto_depIdxs = nil
}
