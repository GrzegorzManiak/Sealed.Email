// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.6
// source: smtp/smtp.proto

package smtp

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

type PrivateInbox struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PepperedUser string `protobuf:"bytes,1,opt,name=pepperedUser,proto3" json:"pepperedUser,omitempty"` // Hash of the usr with a pepper (the sending domain)
	Domain       string `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	EncryptedKey string `protobuf:"bytes,3,opt,name=encryptedKey,proto3" json:"encryptedKey,omitempty"`
}

func (x *PrivateInbox) Reset() {
	*x = PrivateInbox{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smtp_smtp_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrivateInbox) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrivateInbox) ProtoMessage() {}

func (x *PrivateInbox) ProtoReflect() protoreflect.Message {
	mi := &file_smtp_smtp_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrivateInbox.ProtoReflect.Descriptor instead.
func (*PrivateInbox) Descriptor() ([]byte, []int) {
	return file_smtp_smtp_proto_rawDescGZIP(), []int{0}
}

func (x *PrivateInbox) GetPepperedUser() string {
	if x != nil {
		return x.PepperedUser
	}
	return ""
}

func (x *PrivateInbox) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *PrivateInbox) GetEncryptedKey() string {
	if x != nil {
		return x.EncryptedKey
	}
	return ""
}

type PublicInbox struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User   string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Domain string `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	Email  string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *PublicInbox) Reset() {
	*x = PublicInbox{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smtp_smtp_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublicInbox) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublicInbox) ProtoMessage() {}

func (x *PublicInbox) ProtoReflect() protoreflect.Message {
	mi := &file_smtp_smtp_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublicInbox.ProtoReflect.Descriptor instead.
func (*PublicInbox) Descriptor() ([]byte, []int) {
	return file_smtp_smtp_proto_rawDescGZIP(), []int{1}
}

func (x *PublicInbox) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *PublicInbox) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *PublicInbox) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type EncryptedEmail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EncryptedFrom     string            `protobuf:"bytes,1,opt,name=encryptedFrom,proto3" json:"encryptedFrom,omitempty"`       // Fully formatted "from" and "to" header value, e.g., "John Doe <johndoe@example.com>"
	PepperedUsername  string            `protobuf:"bytes,2,opt,name=pepperedUsername,proto3" json:"pepperedUsername,omitempty"` // Hash of the username with a pepper (the sending domain + server name)
	FromDomain        string            `protobuf:"bytes,3,opt,name=fromDomain,proto3" json:"fromDomain,omitempty"`             // The domain of the sender's email address
	To                *PrivateInbox     `protobuf:"bytes,4,opt,name=to,proto3" json:"to,omitempty"`                             // The primary recipient's private inbox
	Cc                []*PrivateInbox   `protobuf:"bytes,5,rep,name=cc,proto3" json:"cc,omitempty"`                             // Carbon copy recipients
	Bcc               []*PrivateInbox   `protobuf:"bytes,6,rep,name=bcc,proto3" json:"bcc,omitempty"`                           // Blind carbon copy recipients
	EncryptedSubject  string            `protobuf:"bytes,7,opt,name=encryptedSubject,proto3" json:"encryptedSubject,omitempty"`
	EncryptedBody     string            `protobuf:"bytes,8,opt,name=encryptedBody,proto3" json:"encryptedBody,omitempty"`
	BounceAddressHash string            `protobuf:"bytes,10,opt,name=bounceAddressHash,proto3" json:"bounceAddressHash,omitempty"`
	BoundDomain       string            `protobuf:"bytes,11,opt,name=boundDomain,proto3" json:"boundDomain,omitempty"`
	DkimSignature     string            `protobuf:"bytes,12,opt,name=dkimSignature,proto3" json:"dkimSignature,omitempty"`
	Attachments       []string          `protobuf:"bytes,13,rep,name=attachments,proto3" json:"attachments,omitempty"`
	ContentType       string            `protobuf:"bytes,14,opt,name=contentType,proto3" json:"contentType,omitempty"` // Content type, e.g., "text/plain", "text/html"
	Headers           map[string]string `protobuf:"bytes,15,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ReplyChain        string            `protobuf:"bytes,16,opt,name=replyChain,proto3" json:"replyChain,omitempty"`
	Date              int64             `protobuf:"varint,17,opt,name=date,proto3" json:"date,omitempty"`
	ServerName        string            `protobuf:"bytes,18,opt,name=serverName,proto3" json:"serverName,omitempty"` // The name of the sending server e.g., "NOISE_V1.0.0"
	Version           string            `protobuf:"bytes,19,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *EncryptedEmail) Reset() {
	*x = EncryptedEmail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smtp_smtp_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncryptedEmail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncryptedEmail) ProtoMessage() {}

func (x *EncryptedEmail) ProtoReflect() protoreflect.Message {
	mi := &file_smtp_smtp_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncryptedEmail.ProtoReflect.Descriptor instead.
func (*EncryptedEmail) Descriptor() ([]byte, []int) {
	return file_smtp_smtp_proto_rawDescGZIP(), []int{2}
}

func (x *EncryptedEmail) GetEncryptedFrom() string {
	if x != nil {
		return x.EncryptedFrom
	}
	return ""
}

func (x *EncryptedEmail) GetPepperedUsername() string {
	if x != nil {
		return x.PepperedUsername
	}
	return ""
}

func (x *EncryptedEmail) GetFromDomain() string {
	if x != nil {
		return x.FromDomain
	}
	return ""
}

func (x *EncryptedEmail) GetTo() *PrivateInbox {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *EncryptedEmail) GetCc() []*PrivateInbox {
	if x != nil {
		return x.Cc
	}
	return nil
}

func (x *EncryptedEmail) GetBcc() []*PrivateInbox {
	if x != nil {
		return x.Bcc
	}
	return nil
}

func (x *EncryptedEmail) GetEncryptedSubject() string {
	if x != nil {
		return x.EncryptedSubject
	}
	return ""
}

func (x *EncryptedEmail) GetEncryptedBody() string {
	if x != nil {
		return x.EncryptedBody
	}
	return ""
}

func (x *EncryptedEmail) GetBounceAddressHash() string {
	if x != nil {
		return x.BounceAddressHash
	}
	return ""
}

func (x *EncryptedEmail) GetBoundDomain() string {
	if x != nil {
		return x.BoundDomain
	}
	return ""
}

func (x *EncryptedEmail) GetDkimSignature() string {
	if x != nil {
		return x.DkimSignature
	}
	return ""
}

func (x *EncryptedEmail) GetAttachments() []string {
	if x != nil {
		return x.Attachments
	}
	return nil
}

func (x *EncryptedEmail) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

func (x *EncryptedEmail) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *EncryptedEmail) GetReplyChain() string {
	if x != nil {
		return x.ReplyChain
	}
	return ""
}

func (x *EncryptedEmail) GetDate() int64 {
	if x != nil {
		return x.Date
	}
	return 0
}

func (x *EncryptedEmail) GetServerName() string {
	if x != nil {
		return x.ServerName
	}
	return ""
}

func (x *EncryptedEmail) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type PublicEmail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromDomain    string            `protobuf:"bytes,1,opt,name=fromDomain,proto3" json:"fromDomain,omitempty"`
	FromUser      string            `protobuf:"bytes,2,opt,name=fromUser,proto3" json:"fromUser,omitempty"`
	FromName      string            `protobuf:"bytes,3,opt,name=fromName,proto3" json:"fromName,omitempty"`
	To            *PublicInbox      `protobuf:"bytes,4,opt,name=to,proto3" json:"to,omitempty"`
	Cc            []*PublicInbox    `protobuf:"bytes,5,rep,name=cc,proto3" json:"cc,omitempty"`
	Bcc           []*PublicInbox    `protobuf:"bytes,6,rep,name=bcc,proto3" json:"bcc,omitempty"`
	Subject       string            `protobuf:"bytes,7,opt,name=subject,proto3" json:"subject,omitempty"`
	Body          string            `protobuf:"bytes,8,opt,name=body,proto3" json:"body,omitempty"`
	ContentType   string            `protobuf:"bytes,9,opt,name=contentType,proto3" json:"contentType,omitempty"` // Content type, e.g., "text/plain", "text/html"
	Headers       map[string]string `protobuf:"bytes,10,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Attachments   []string          `protobuf:"bytes,11,rep,name=attachments,proto3" json:"attachments,omitempty"`
	ReplyTo       string            `protobuf:"bytes,12,opt,name=replyTo,proto3" json:"replyTo,omitempty"`
	DkimSignature string            `protobuf:"bytes,13,opt,name=dkimSignature,proto3" json:"dkimSignature,omitempty"`
	Date          int64             `protobuf:"varint,14,opt,name=date,proto3" json:"date,omitempty"`
	Version       string            `protobuf:"bytes,15,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *PublicEmail) Reset() {
	*x = PublicEmail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smtp_smtp_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublicEmail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublicEmail) ProtoMessage() {}

func (x *PublicEmail) ProtoReflect() protoreflect.Message {
	mi := &file_smtp_smtp_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublicEmail.ProtoReflect.Descriptor instead.
func (*PublicEmail) Descriptor() ([]byte, []int) {
	return file_smtp_smtp_proto_rawDescGZIP(), []int{3}
}

func (x *PublicEmail) GetFromDomain() string {
	if x != nil {
		return x.FromDomain
	}
	return ""
}

func (x *PublicEmail) GetFromUser() string {
	if x != nil {
		return x.FromUser
	}
	return ""
}

func (x *PublicEmail) GetFromName() string {
	if x != nil {
		return x.FromName
	}
	return ""
}

func (x *PublicEmail) GetTo() *PublicInbox {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *PublicEmail) GetCc() []*PublicInbox {
	if x != nil {
		return x.Cc
	}
	return nil
}

func (x *PublicEmail) GetBcc() []*PublicInbox {
	if x != nil {
		return x.Bcc
	}
	return nil
}

func (x *PublicEmail) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *PublicEmail) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *PublicEmail) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

func (x *PublicEmail) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *PublicEmail) GetAttachments() []string {
	if x != nil {
		return x.Attachments
	}
	return nil
}

func (x *PublicEmail) GetReplyTo() string {
	if x != nil {
		return x.ReplyTo
	}
	return ""
}

func (x *PublicEmail) GetDkimSignature() string {
	if x != nil {
		return x.DkimSignature
	}
	return ""
}

func (x *PublicEmail) GetDate() int64 {
	if x != nil {
		return x.Date
	}
	return 0
}

func (x *PublicEmail) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type SendEmailResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error   string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *SendEmailResponse) Reset() {
	*x = SendEmailResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smtp_smtp_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendEmailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendEmailResponse) ProtoMessage() {}

func (x *SendEmailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_smtp_smtp_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendEmailResponse.ProtoReflect.Descriptor instead.
func (*SendEmailResponse) Descriptor() ([]byte, []int) {
	return file_smtp_smtp_proto_rawDescGZIP(), []int{4}
}

func (x *SendEmailResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *SendEmailResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_smtp_smtp_proto protoreflect.FileDescriptor

var file_smtp_smtp_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x73, 0x6d, 0x74, 0x70, 0x2f, 0x73, 0x6d, 0x74, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x73, 0x6d, 0x74, 0x70, 0x22, 0x6e, 0x0a, 0x0c, 0x50, 0x72, 0x69, 0x76, 0x61,
	0x74, 0x65, 0x49, 0x6e, 0x62, 0x6f, 0x78, 0x12, 0x22, 0x0a, 0x0c, 0x70, 0x65, 0x70, 0x70, 0x65,
	0x72, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70,
	0x65, 0x70, 0x70, 0x65, 0x72, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64,
	0x4b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x6e, 0x63, 0x72, 0x79,
	0x70, 0x74, 0x65, 0x64, 0x4b, 0x65, 0x79, 0x22, 0x4f, 0x0a, 0x0b, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x49, 0x6e, 0x62, 0x6f, 0x78, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0xe3, 0x05, 0x0a, 0x0e, 0x45, 0x6e, 0x63,
	0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x24, 0x0a, 0x0d, 0x65,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x46, 0x72, 0x6f,
	0x6d, 0x12, 0x2a, 0x0a, 0x10, 0x70, 0x65, 0x70, 0x70, 0x65, 0x72, 0x65, 0x64, 0x55, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x70, 0x65, 0x70,
	0x70, 0x65, 0x72, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a,
	0x0a, 0x66, 0x72, 0x6f, 0x6d, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x66, 0x72, 0x6f, 0x6d, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x22, 0x0a,
	0x02, 0x74, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x6d, 0x74, 0x70,
	0x2e, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x62, 0x6f, 0x78, 0x52, 0x02, 0x74,
	0x6f, 0x12, 0x22, 0x0a, 0x02, 0x63, 0x63, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x73, 0x6d, 0x74, 0x70, 0x2e, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x62, 0x6f,
	0x78, 0x52, 0x02, 0x63, 0x63, 0x12, 0x24, 0x0a, 0x03, 0x62, 0x63, 0x63, 0x18, 0x06, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x6d, 0x74, 0x70, 0x2e, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74,
	0x65, 0x49, 0x6e, 0x62, 0x6f, 0x78, 0x52, 0x03, 0x62, 0x63, 0x63, 0x12, 0x2a, 0x0a, 0x10, 0x65,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x53, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64,
	0x53, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x65, 0x6e, 0x63, 0x72, 0x79,
	0x70, 0x74, 0x65, 0x64, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x2c, 0x0a,
	0x11, 0x62, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x48, 0x61,
	0x73, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x62, 0x6f, 0x75, 0x6e, 0x63, 0x65,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x48, 0x61, 0x73, 0x68, 0x12, 0x20, 0x0a, 0x0b, 0x62,
	0x6f, 0x75, 0x6e, 0x64, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x24, 0x0a,
	0x0d, 0x64, 0x6b, 0x69, 0x6d, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x64, 0x6b, 0x69, 0x6d, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x3b, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x73, 0x18, 0x0f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x73, 0x6d, 0x74, 0x70, 0x2e,
	0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x43, 0x68, 0x61,
	0x69, 0x6e, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x43,
	0x68, 0x61, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x11, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x13, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x1a, 0x3a, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xa6,
	0x04, 0x0a, 0x0b, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1e,
	0x0a, 0x0a, 0x66, 0x72, 0x6f, 0x6d, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x66, 0x72, 0x6f, 0x6d, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x1a,
	0x0a, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x55, 0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x72,
	0x6f, 0x6d, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x72,
	0x6f, 0x6d, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73, 0x6d, 0x74, 0x70, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x49, 0x6e, 0x62, 0x6f, 0x78, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x21, 0x0a, 0x02, 0x63, 0x63, 0x18,
	0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73, 0x6d, 0x74, 0x70, 0x2e, 0x50, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x49, 0x6e, 0x62, 0x6f, 0x78, 0x52, 0x02, 0x63, 0x63, 0x12, 0x23, 0x0a, 0x03,
	0x62, 0x63, 0x63, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73, 0x6d, 0x74, 0x70,
	0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x49, 0x6e, 0x62, 0x6f, 0x78, 0x52, 0x03, 0x62, 0x63,
	0x63, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12,
	0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x38, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x0a, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x73, 0x6d, 0x74, 0x70, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x45, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x61,
	0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0b, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x72, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x12, 0x24, 0x0a, 0x0d, 0x64, 0x6b, 0x69, 0x6d, 0x53,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x64, 0x6b, 0x69, 0x6d, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0f, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x3a, 0x0a, 0x0c, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x43, 0x0a, 0x11, 0x53, 0x65, 0x6e, 0x64, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x52, 0x0a, 0x0b,
	0x53, 0x6d, 0x74, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x12, 0x53,
	0x65, 0x6e, 0x64, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x45, 0x6d, 0x61, 0x69,
	0x6c, 0x12, 0x14, 0x2e, 0x73, 0x6d, 0x74, 0x70, 0x2e, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74,
	0x65, 0x64, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x1a, 0x17, 0x2e, 0x73, 0x6d, 0x74, 0x70, 0x2e, 0x53,
	0x65, 0x6e, 0x64, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x47,
	0x72, 0x7a, 0x65, 0x67, 0x6f, 0x72, 0x7a, 0x4d, 0x61, 0x6e, 0x69, 0x61, 0x6b, 0x2f, 0x4e, 0x6f,
	0x69, 0x73, 0x65, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x6d, 0x74, 0x70, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_smtp_smtp_proto_rawDescOnce sync.Once
	file_smtp_smtp_proto_rawDescData = file_smtp_smtp_proto_rawDesc
)

func file_smtp_smtp_proto_rawDescGZIP() []byte {
	file_smtp_smtp_proto_rawDescOnce.Do(func() {
		file_smtp_smtp_proto_rawDescData = protoimpl.X.CompressGZIP(file_smtp_smtp_proto_rawDescData)
	})
	return file_smtp_smtp_proto_rawDescData
}

var file_smtp_smtp_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_smtp_smtp_proto_goTypes = []interface{}{
	(*PrivateInbox)(nil),      // 0: smtp.PrivateInbox
	(*PublicInbox)(nil),       // 1: smtp.PublicInbox
	(*EncryptedEmail)(nil),    // 2: smtp.EncryptedEmail
	(*PublicEmail)(nil),       // 3: smtp.PublicEmail
	(*SendEmailResponse)(nil), // 4: smtp.SendEmailResponse
	nil,                       // 5: smtp.EncryptedEmail.HeadersEntry
	nil,                       // 6: smtp.PublicEmail.HeadersEntry
}
var file_smtp_smtp_proto_depIdxs = []int32{
	0, // 0: smtp.EncryptedEmail.to:type_name -> smtp.PrivateInbox
	0, // 1: smtp.EncryptedEmail.cc:type_name -> smtp.PrivateInbox
	0, // 2: smtp.EncryptedEmail.bcc:type_name -> smtp.PrivateInbox
	5, // 3: smtp.EncryptedEmail.headers:type_name -> smtp.EncryptedEmail.HeadersEntry
	1, // 4: smtp.PublicEmail.to:type_name -> smtp.PublicInbox
	1, // 5: smtp.PublicEmail.cc:type_name -> smtp.PublicInbox
	1, // 6: smtp.PublicEmail.bcc:type_name -> smtp.PublicInbox
	6, // 7: smtp.PublicEmail.headers:type_name -> smtp.PublicEmail.HeadersEntry
	2, // 8: smtp.SmtpService.SendEncryptedEmail:input_type -> smtp.EncryptedEmail
	4, // 9: smtp.SmtpService.SendEncryptedEmail:output_type -> smtp.SendEmailResponse
	9, // [9:10] is the sub-list for method output_type
	8, // [8:9] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_smtp_smtp_proto_init() }
func file_smtp_smtp_proto_init() {
	if File_smtp_smtp_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_smtp_smtp_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrivateInbox); i {
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
		file_smtp_smtp_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublicInbox); i {
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
		file_smtp_smtp_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EncryptedEmail); i {
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
		file_smtp_smtp_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublicEmail); i {
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
		file_smtp_smtp_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendEmailResponse); i {
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
			RawDescriptor: file_smtp_smtp_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_smtp_smtp_proto_goTypes,
		DependencyIndexes: file_smtp_smtp_proto_depIdxs,
		MessageInfos:      file_smtp_smtp_proto_msgTypes,
	}.Build()
	File_smtp_smtp_proto = out.File
	file_smtp_smtp_proto_rawDesc = nil
	file_smtp_smtp_proto_goTypes = nil
	file_smtp_smtp_proto_depIdxs = nil
}
