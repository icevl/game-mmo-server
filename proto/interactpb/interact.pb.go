// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: proto/interactpb/interact.proto

package interactpb

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

type Interact struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UUID string `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"` //"hand", "contact"
}

func (x *Interact) Reset() {
	*x = Interact{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_interactpb_interact_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Interact) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Interact) ProtoMessage() {}

func (x *Interact) ProtoReflect() protoreflect.Message {
	mi := &file_proto_interactpb_interact_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Interact.ProtoReflect.Descriptor instead.
func (*Interact) Descriptor() ([]byte, []int) {
	return file_proto_interactpb_interact_proto_rawDescGZIP(), []int{0}
}

func (x *Interact) GetUUID() string {
	if x != nil {
		return x.UUID
	}
	return ""
}

func (x *Interact) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type InteractWith struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UUID       string `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	TargetUuid string `protobuf:"bytes,2,opt,name=target_uuid,json=targetUuid,proto3" json:"target_uuid,omitempty"`
	Type       string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"` //"hand", "contact"
}

func (x *InteractWith) Reset() {
	*x = InteractWith{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_interactpb_interact_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InteractWith) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InteractWith) ProtoMessage() {}

func (x *InteractWith) ProtoReflect() protoreflect.Message {
	mi := &file_proto_interactpb_interact_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InteractWith.ProtoReflect.Descriptor instead.
func (*InteractWith) Descriptor() ([]byte, []int) {
	return file_proto_interactpb_interact_proto_rawDescGZIP(), []int{1}
}

func (x *InteractWith) GetUUID() string {
	if x != nil {
		return x.UUID
	}
	return ""
}

func (x *InteractWith) GetTargetUuid() string {
	if x != nil {
		return x.TargetUuid
	}
	return ""
}

func (x *InteractWith) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type InteractQueue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UUID string `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
}

func (x *InteractQueue) Reset() {
	*x = InteractQueue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_interactpb_interact_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InteractQueue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InteractQueue) ProtoMessage() {}

func (x *InteractQueue) ProtoReflect() protoreflect.Message {
	mi := &file_proto_interactpb_interact_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InteractQueue.ProtoReflect.Descriptor instead.
func (*InteractQueue) Descriptor() ([]byte, []int) {
	return file_proto_interactpb_interact_proto_rawDescGZIP(), []int{2}
}

func (x *InteractQueue) GetUUID() string {
	if x != nil {
		return x.UUID
	}
	return ""
}

var File_proto_interactpb_interact_proto protoreflect.FileDescriptor

var file_proto_interactpb_interact_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74,
	0x70, 0x62, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22, 0x32, 0x0a, 0x08, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x55, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x55, 0x55, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22,
	0x57, 0x0a, 0x0c, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x57, 0x69, 0x74, 0x68, 0x12,
	0x12, 0x0a, 0x04, 0x55, 0x55, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x55,
	0x55, 0x49, 0x44, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x75,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x55, 0x75, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x23, 0x0a, 0x0d, 0x49, 0x6e, 0x74, 0x65,
	0x72, 0x61, 0x63, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x55, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x55, 0x55, 0x49, 0x44, 0x42, 0x19, 0x5a,
	0x17, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_interactpb_interact_proto_rawDescOnce sync.Once
	file_proto_interactpb_interact_proto_rawDescData = file_proto_interactpb_interact_proto_rawDesc
)

func file_proto_interactpb_interact_proto_rawDescGZIP() []byte {
	file_proto_interactpb_interact_proto_rawDescOnce.Do(func() {
		file_proto_interactpb_interact_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_interactpb_interact_proto_rawDescData)
	})
	return file_proto_interactpb_interact_proto_rawDescData
}

var file_proto_interactpb_interact_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_interactpb_interact_proto_goTypes = []interface{}{
	(*Interact)(nil),      // 0: messages.Interact
	(*InteractWith)(nil),  // 1: messages.InteractWith
	(*InteractQueue)(nil), // 2: messages.InteractQueue
}
var file_proto_interactpb_interact_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_interactpb_interact_proto_init() }
func file_proto_interactpb_interact_proto_init() {
	if File_proto_interactpb_interact_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_interactpb_interact_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Interact); i {
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
		file_proto_interactpb_interact_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InteractWith); i {
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
		file_proto_interactpb_interact_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InteractQueue); i {
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
			RawDescriptor: file_proto_interactpb_interact_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_interactpb_interact_proto_goTypes,
		DependencyIndexes: file_proto_interactpb_interact_proto_depIdxs,
		MessageInfos:      file_proto_interactpb_interact_proto_msgTypes,
	}.Build()
	File_proto_interactpb_interact_proto = out.File
	file_proto_interactpb_interact_proto_rawDesc = nil
	file_proto_interactpb_interact_proto_goTypes = nil
	file_proto_interactpb_interact_proto_depIdxs = nil
}
