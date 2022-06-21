// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.1
// source: package.proto

package versatilis

import (
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

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StringMessage string `protobuf:"bytes,10,opt,name=stringMessage,proto3" json:"stringMessage,omitempty"`
	BytesMessage  []byte `protobuf:"bytes,20,opt,name=bytesMessage,proto3" json:"bytesMessage,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_package_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_package_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_package_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetStringMessage() string {
	if x != nil {
		return x.StringMessage
	}
	return ""
}

func (x *Message) GetBytesMessage() []byte {
	if x != nil {
		return x.BytesMessage
	}
	return nil
}

type Package struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version            string `protobuf:"bytes,10,opt,name=version,proto3" json:"version,omitempty"`
	NoiseHandshakeInfo []byte `protobuf:"bytes,20,opt,name=noiseHandshakeInfo,proto3" json:"noiseHandshakeInfo,omitempty"`
	NoiseCiphertext    []byte `protobuf:"bytes,23,opt,name=noiseCiphertext,proto3" json:"noiseCiphertext,omitempty"`
	NoiseAuthTag       []byte `protobuf:"bytes,25,opt,name=noiseAuthTag,proto3" json:"noiseAuthTag,omitempty"`
}

func (x *Package) Reset() {
	*x = Package{}
	if protoimpl.UnsafeEnabled {
		mi := &file_package_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Package) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Package) ProtoMessage() {}

func (x *Package) ProtoReflect() protoreflect.Message {
	mi := &file_package_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Package.ProtoReflect.Descriptor instead.
func (*Package) Descriptor() ([]byte, []int) {
	return file_package_proto_rawDescGZIP(), []int{1}
}

func (x *Package) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Package) GetNoiseHandshakeInfo() []byte {
	if x != nil {
		return x.NoiseHandshakeInfo
	}
	return nil
}

func (x *Package) GetNoiseCiphertext() []byte {
	if x != nil {
		return x.NoiseCiphertext
	}
	return nil
}

func (x *Package) GetNoiseAuthTag() []byte {
	if x != nil {
		return x.NoiseAuthTag
	}
	return nil
}

var File_package_proto protoreflect.FileDescriptor

var file_package_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0a, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6c, 0x69, 0x73, 0x22, 0x53, 0x0a, 0x07, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x22, 0x0a, 0x0c,
	0x62, 0x79, 0x74, 0x65, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x14, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x0c, 0x62, 0x79, 0x74, 0x65, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0xa1, 0x01, 0x0a, 0x07, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2e, 0x0a, 0x12, 0x6e, 0x6f, 0x69, 0x73, 0x65, 0x48,
	0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x14, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x12, 0x6e, 0x6f, 0x69, 0x73, 0x65, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61,
	0x6b, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x28, 0x0a, 0x0f, 0x6e, 0x6f, 0x69, 0x73, 0x65, 0x43,
	0x69, 0x70, 0x68, 0x65, 0x72, 0x74, 0x65, 0x78, 0x74, 0x18, 0x17, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0f, 0x6e, 0x6f, 0x69, 0x73, 0x65, 0x43, 0x69, 0x70, 0x68, 0x65, 0x72, 0x74, 0x65, 0x78, 0x74,
	0x12, 0x22, 0x0a, 0x0c, 0x6e, 0x6f, 0x69, 0x73, 0x65, 0x41, 0x75, 0x74, 0x68, 0x54, 0x61, 0x67,
	0x18, 0x19, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x6e, 0x6f, 0x69, 0x73, 0x65, 0x41, 0x75, 0x74,
	0x68, 0x54, 0x61, 0x67, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2e, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x61,
	0x74, 0x69, 0x6c, 0x69, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_package_proto_rawDescOnce sync.Once
	file_package_proto_rawDescData = file_package_proto_rawDesc
)

func file_package_proto_rawDescGZIP() []byte {
	file_package_proto_rawDescOnce.Do(func() {
		file_package_proto_rawDescData = protoimpl.X.CompressGZIP(file_package_proto_rawDescData)
	})
	return file_package_proto_rawDescData
}

var file_package_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_package_proto_goTypes = []interface{}{
	(*Message)(nil), // 0: versatilis.Message
	(*Package)(nil), // 1: versatilis.Package
}
var file_package_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_package_proto_init() }
func file_package_proto_init() {
	if File_package_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_package_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_package_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Package); i {
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
			RawDescriptor: file_package_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_package_proto_goTypes,
		DependencyIndexes: file_package_proto_depIdxs,
		MessageInfos:      file_package_proto_msgTypes,
	}.Build()
	File_package_proto = out.File
	file_package_proto_rawDesc = nil
	file_package_proto_goTypes = nil
	file_package_proto_depIdxs = nil
}
