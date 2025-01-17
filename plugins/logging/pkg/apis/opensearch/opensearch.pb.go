// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	ragù          v0.2.3
// source: plugins/logging/pkg/apis/opensearch/opensearch.proto

package opensearch

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

type ClusterReference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthorizedClusterID string `protobuf:"bytes,1,opt,name=AuthorizedClusterID,proto3" json:"AuthorizedClusterID,omitempty"`
}

func (x *ClusterReference) Reset() {
	*x = ClusterReference{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugins_logging_pkg_apis_opensearch_opensearch_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterReference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterReference) ProtoMessage() {}

func (x *ClusterReference) ProtoReflect() protoreflect.Message {
	mi := &file_plugins_logging_pkg_apis_opensearch_opensearch_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterReference.ProtoReflect.Descriptor instead.
func (*ClusterReference) Descriptor() ([]byte, []int) {
	return file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDescGZIP(), []int{0}
}

func (x *ClusterReference) GetAuthorizedClusterID() string {
	if x != nil {
		return x.AuthorizedClusterID
	}
	return ""
}

type OpensearchDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username    string `protobuf:"bytes,1,opt,name=Username,proto3" json:"Username,omitempty"`
	Password    string `protobuf:"bytes,2,opt,name=Password,proto3" json:"Password,omitempty"`
	ExternalURL string `protobuf:"bytes,3,opt,name=ExternalURL,proto3" json:"ExternalURL,omitempty"`
}

func (x *OpensearchDetails) Reset() {
	*x = OpensearchDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugins_logging_pkg_apis_opensearch_opensearch_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpensearchDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpensearchDetails) ProtoMessage() {}

func (x *OpensearchDetails) ProtoReflect() protoreflect.Message {
	mi := &file_plugins_logging_pkg_apis_opensearch_opensearch_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpensearchDetails.ProtoReflect.Descriptor instead.
func (*OpensearchDetails) Descriptor() ([]byte, []int) {
	return file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDescGZIP(), []int{1}
}

func (x *OpensearchDetails) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *OpensearchDetails) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *OpensearchDetails) GetExternalURL() string {
	if x != nil {
		return x.ExternalURL
	}
	return ""
}

var File_plugins_logging_pkg_apis_opensearch_opensearch_proto protoreflect.FileDescriptor

var file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDesc = []byte{
	0x0a, 0x34, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x2f, 0x6c, 0x6f, 0x67, 0x67, 0x69, 0x6e,
	0x67, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x73,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6f, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x22, 0x33, 0x0a, 0x10, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x66,
	0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x1d, 0x0a, 0x13, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x69, 0x7a, 0x65, 0x64, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x00, 0x3a, 0x00, 0x22, 0x54, 0x0a, 0x11, 0x4f, 0x70, 0x65, 0x6e, 0x73,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x12, 0x0a, 0x08,
	0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x00,
	0x12, 0x12, 0x0a, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x00, 0x12, 0x15, 0x0a, 0x0b, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x55, 0x52, 0x4c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x00, 0x3a, 0x00, 0x32, 0x5f, 0x0a,
	0x0a, 0x4f, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x4f, 0x0a, 0x0a, 0x47,
	0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x1c, 0x2e, 0x6f, 0x70, 0x65, 0x6e,
	0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x1a, 0x1d, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x73, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x2e, 0x4f, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x44,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x00, 0x28, 0x00, 0x30, 0x00, 0x1a, 0x00, 0x42, 0x3d,
	0x5a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x6e,
	0x63, 0x68, 0x65, 0x72, 0x2f, 0x6f, 0x70, 0x6e, 0x69, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x73, 0x2f, 0x6c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70,
	0x69, 0x73, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDescOnce sync.Once
	file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDescData = file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDesc
)

func file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDescGZIP() []byte {
	file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDescOnce.Do(func() {
		file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDescData = protoimpl.X.CompressGZIP(file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDescData)
	})
	return file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDescData
}

var file_plugins_logging_pkg_apis_opensearch_opensearch_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_plugins_logging_pkg_apis_opensearch_opensearch_proto_goTypes = []interface{}{
	(*ClusterReference)(nil),  // 0: opensearch.ClusterReference
	(*OpensearchDetails)(nil), // 1: opensearch.OpensearchDetails
}
var file_plugins_logging_pkg_apis_opensearch_opensearch_proto_depIdxs = []int32{
	0, // 0: opensearch.Opensearch.GetDetails:input_type -> opensearch.ClusterReference
	1, // 1: opensearch.Opensearch.GetDetails:output_type -> opensearch.OpensearchDetails
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_plugins_logging_pkg_apis_opensearch_opensearch_proto_init() }
func file_plugins_logging_pkg_apis_opensearch_opensearch_proto_init() {
	if File_plugins_logging_pkg_apis_opensearch_opensearch_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_plugins_logging_pkg_apis_opensearch_opensearch_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterReference); i {
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
		file_plugins_logging_pkg_apis_opensearch_opensearch_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OpensearchDetails); i {
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
			RawDescriptor: file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_plugins_logging_pkg_apis_opensearch_opensearch_proto_goTypes,
		DependencyIndexes: file_plugins_logging_pkg_apis_opensearch_opensearch_proto_depIdxs,
		MessageInfos:      file_plugins_logging_pkg_apis_opensearch_opensearch_proto_msgTypes,
	}.Build()
	File_plugins_logging_pkg_apis_opensearch_opensearch_proto = out.File
	file_plugins_logging_pkg_apis_opensearch_opensearch_proto_rawDesc = nil
	file_plugins_logging_pkg_apis_opensearch_opensearch_proto_goTypes = nil
	file_plugins_logging_pkg_apis_opensearch_opensearch_proto_depIdxs = nil
}
