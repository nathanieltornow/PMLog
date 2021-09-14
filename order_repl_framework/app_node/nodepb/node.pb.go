// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.3
// source: nodepb/node.proto

package nodepb

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

type Prep struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LocalToken uint64 `protobuf:"varint,1,opt,name=localToken,proto3" json:"localToken,omitempty"`
	Color      uint32 `protobuf:"varint,2,opt,name=color,proto3" json:"color,omitempty"`
	Content    string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *Prep) Reset() {
	*x = Prep{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nodepb_node_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Prep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Prep) ProtoMessage() {}

func (x *Prep) ProtoReflect() protoreflect.Message {
	mi := &file_nodepb_node_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Prep.ProtoReflect.Descriptor instead.
func (*Prep) Descriptor() ([]byte, []int) {
	return file_nodepb_node_proto_rawDescGZIP(), []int{0}
}

func (x *Prep) GetLocalToken() uint64 {
	if x != nil {
		return x.LocalToken
	}
	return 0
}

func (x *Prep) GetColor() uint32 {
	if x != nil {
		return x.Color
	}
	return 0
}

func (x *Prep) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type Ack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LocalToken uint64 `protobuf:"varint,1,opt,name=localToken,proto3" json:"localToken,omitempty"`
}

func (x *Ack) Reset() {
	*x = Ack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nodepb_node_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ack) ProtoMessage() {}

func (x *Ack) ProtoReflect() protoreflect.Message {
	mi := &file_nodepb_node_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ack.ProtoReflect.Descriptor instead.
func (*Ack) Descriptor() ([]byte, []int) {
	return file_nodepb_node_proto_rawDescGZIP(), []int{1}
}

func (x *Ack) GetLocalToken() uint64 {
	if x != nil {
		return x.LocalToken
	}
	return 0
}

type AckReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID uint32 `protobuf:"varint,1,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
}

func (x *AckReq) Reset() {
	*x = AckReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nodepb_node_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AckReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckReq) ProtoMessage() {}

func (x *AckReq) ProtoReflect() protoreflect.Message {
	mi := &file_nodepb_node_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckReq.ProtoReflect.Descriptor instead.
func (*AckReq) Descriptor() ([]byte, []int) {
	return file_nodepb_node_proto_rawDescGZIP(), []int{2}
}

func (x *AckReq) GetNodeID() uint32 {
	if x != nil {
		return x.NodeID
	}
	return 0
}

type Com struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LocalToken uint64 `protobuf:"varint,1,opt,name=localToken,proto3" json:"localToken,omitempty"`
}

func (x *Com) Reset() {
	*x = Com{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nodepb_node_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Com) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Com) ProtoMessage() {}

func (x *Com) ProtoReflect() protoreflect.Message {
	mi := &file_nodepb_node_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Com.ProtoReflect.Descriptor instead.
func (*Com) Descriptor() ([]byte, []int) {
	return file_nodepb_node_proto_rawDescGZIP(), []int{3}
}

func (x *Com) GetLocalToken() uint64 {
	if x != nil {
		return x.LocalToken
	}
	return 0
}

type RegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IP string `protobuf:"bytes,1,opt,name=IP,proto3" json:"IP,omitempty"`
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nodepb_node_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nodepb_node_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_nodepb_node_proto_rawDescGZIP(), []int{4}
}

func (x *RegisterRequest) GetIP() string {
	if x != nil {
		return x.IP
	}
	return ""
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nodepb_node_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_nodepb_node_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_nodepb_node_proto_rawDescGZIP(), []int{5}
}

var File_nodepb_node_proto protoreflect.FileDescriptor

var file_nodepb_node_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6e, 0x6f, 0x64, 0x65, 0x70, 0x62, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x70, 0x70, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0x56, 0x0a,
	0x04, 0x50, 0x72, 0x65, 0x70, 0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x6c, 0x6f, 0x63, 0x61, 0x6c,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x25, 0x0a, 0x03, 0x41, 0x63, 0x6b, 0x12, 0x1e, 0x0a, 0x0a,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0a, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x20, 0x0a, 0x06,
	0x41, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x22, 0x25,
	0x0a, 0x03, 0x43, 0x6f, 0x6d, 0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x6c, 0x6f, 0x63, 0x61, 0x6c,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x21, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x50, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x50, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x32, 0xc6, 0x01, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x36, 0x0a, 0x08, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x61, 0x70, 0x70, 0x5f, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x70, 0x70, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x2c, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12, 0x0e, 0x2e,
	0x61, 0x70, 0x70, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x1a, 0x0f, 0x2e,
	0x61, 0x70, 0x70, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x28, 0x01,
	0x12, 0x2c, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x41, 0x63, 0x6b, 0x73, 0x12, 0x10, 0x2e, 0x61, 0x70,
	0x70, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x41, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x1a, 0x0d, 0x2e,
	0x61, 0x70, 0x70, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x41, 0x63, 0x6b, 0x30, 0x01, 0x12, 0x2a,
	0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x0d, 0x2e, 0x61, 0x70, 0x70, 0x5f, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x1a, 0x0f, 0x2e, 0x61, 0x70, 0x70, 0x5f, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x28, 0x01, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x74, 0x68, 0x61, 0x6e, 0x69,
	0x65, 0x6c, 0x74, 0x6f, 0x72, 0x6e, 0x6f, 0x77, 0x2f, 0x61, 0x70, 0x70, 0x5f, 0x6e, 0x6f, 0x64,
	0x65, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_nodepb_node_proto_rawDescOnce sync.Once
	file_nodepb_node_proto_rawDescData = file_nodepb_node_proto_rawDesc
)

func file_nodepb_node_proto_rawDescGZIP() []byte {
	file_nodepb_node_proto_rawDescOnce.Do(func() {
		file_nodepb_node_proto_rawDescData = protoimpl.X.CompressGZIP(file_nodepb_node_proto_rawDescData)
	})
	return file_nodepb_node_proto_rawDescData
}

var file_nodepb_node_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_nodepb_node_proto_goTypes = []interface{}{
	(*Prep)(nil),            // 0: app_node.Prep
	(*Ack)(nil),             // 1: app_node.Ack
	(*AckReq)(nil),          // 2: app_node.AckReq
	(*Com)(nil),             // 3: app_node.Com
	(*RegisterRequest)(nil), // 4: app_node.RegisterRequest
	(*Empty)(nil),           // 5: app_node.Empty
}
var file_nodepb_node_proto_depIdxs = []int32{
	4, // 0: app_node.Node.Register:input_type -> app_node.RegisterRequest
	0, // 1: app_node.Node.Prepare:input_type -> app_node.Prep
	2, // 2: app_node.Node.GetAcks:input_type -> app_node.AckReq
	3, // 3: app_node.Node.Commit:input_type -> app_node.Com
	5, // 4: app_node.Node.Register:output_type -> app_node.Empty
	5, // 5: app_node.Node.Prepare:output_type -> app_node.Empty
	1, // 6: app_node.Node.GetAcks:output_type -> app_node.Ack
	5, // 7: app_node.Node.Commit:output_type -> app_node.Empty
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_nodepb_node_proto_init() }
func file_nodepb_node_proto_init() {
	if File_nodepb_node_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_nodepb_node_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Prep); i {
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
		file_nodepb_node_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ack); i {
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
		file_nodepb_node_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AckReq); i {
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
		file_nodepb_node_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Com); i {
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
		file_nodepb_node_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterRequest); i {
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
		file_nodepb_node_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_nodepb_node_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_nodepb_node_proto_goTypes,
		DependencyIndexes: file_nodepb_node_proto_depIdxs,
		MessageInfos:      file_nodepb_node_proto_msgTypes,
	}.Build()
	File_nodepb_node_proto = out.File
	file_nodepb_node_proto_rawDesc = nil
	file_nodepb_node_proto_goTypes = nil
	file_nodepb_node_proto_depIdxs = nil
}
