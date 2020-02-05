// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common.proto

package common_go_proto

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

// Kind represents the kinds of notes supported.
type NoteKind int32

const (
	// Unknown.
	NoteKind_NOTE_KIND_UNSPECIFIED NoteKind = 0
	// The note and occurrence represent a package vulnerability.
	NoteKind_VULNERABILITY NoteKind = 1
	// The note and occurrence assert build provenance.
	NoteKind_BUILD NoteKind = 2
	// This represents an image basis relationship.
	NoteKind_IMAGE NoteKind = 3
	// This represents a package installed via a package manager.
	NoteKind_PACKAGE NoteKind = 4
	// The note and occurrence track deployment events.
	NoteKind_DEPLOYMENT NoteKind = 5
	// The note and occurrence track the initial discovery status of a resource.
	NoteKind_DISCOVERY NoteKind = 6
	// This represents a logical "role" that can attest to artifacts.
	NoteKind_ATTESTATION NoteKind = 7
)

var NoteKind_name = map[int32]string{
	0: "NOTE_KIND_UNSPECIFIED",
	1: "VULNERABILITY",
	2: "BUILD",
	3: "IMAGE",
	4: "PACKAGE",
	5: "DEPLOYMENT",
	6: "DISCOVERY",
	7: "ATTESTATION",
}

var NoteKind_value = map[string]int32{
	"NOTE_KIND_UNSPECIFIED": 0,
	"VULNERABILITY":         1,
	"BUILD":                 2,
	"IMAGE":                 3,
	"PACKAGE":               4,
	"DEPLOYMENT":            5,
	"DISCOVERY":             6,
	"ATTESTATION":           7,
}

func (x NoteKind) String() string {
	return proto.EnumName(NoteKind_name, int32(x))
}

func (NoteKind) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0}
}

// Metadata for any related URL information.
type RelatedUrl struct {
	// Specific URL associated with the resource.
	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	// Label to describe usage of the URL.
	Label                string   `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RelatedUrl) Reset()         { *m = RelatedUrl{} }
func (m *RelatedUrl) String() string { return proto.CompactTextString(m) }
func (*RelatedUrl) ProtoMessage()    {}
func (*RelatedUrl) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0}
}

func (m *RelatedUrl) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RelatedUrl.Unmarshal(m, b)
}
func (m *RelatedUrl) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RelatedUrl.Marshal(b, m, deterministic)
}
func (m *RelatedUrl) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelatedUrl.Merge(m, src)
}
func (m *RelatedUrl) XXX_Size() int {
	return xxx_messageInfo_RelatedUrl.Size(m)
}
func (m *RelatedUrl) XXX_DiscardUnknown() {
	xxx_messageInfo_RelatedUrl.DiscardUnknown(m)
}

var xxx_messageInfo_RelatedUrl proto.InternalMessageInfo

func (m *RelatedUrl) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *RelatedUrl) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

// Verifiers (e.g. Kritis implementations) MUST verify signatures
// with respect to the trust anchors defined in policy (e.g. a Kritis policy).
// Typically this means that the verifier has been configured with a map from
// `public_key_id` to public key material (and any required parameters, e.g.
// signing algorithm).
//
// In particular, verification implementations MUST NOT treat the signature
// `public_key_id` as anything more than a key lookup hint. The `public_key_id`
// DOES NOT validate or authenticate a public key; it only provides a mechanism
// for quickly selecting a public key ALREADY CONFIGURED on the verifier through
// a trusted channel. Verification implementations MUST reject signatures in any
// of the following circumstances:
//   * The `public_key_id` is not recognized by the verifier.
//   * The public key that `public_key_id` refers to does not verify the
//     signature with respect to the payload.
//
// The `signature` contents SHOULD NOT be "attached" (where the payload is
// included with the serialized `signature` bytes). Verifiers MUST ignore any
// "attached" payload and only verify signatures with respect to explicitly
// provided payload (e.g. a `payload` field on the proto message that holds
// this Signature, or the canonical serialization of the proto message that
// holds this signature).
type Signature struct {
	// The content of the signature, an opaque bytestring.
	// The payload that this signature verifies MUST be unambiguously provided
	// with the Signature during verification. A wrapper message might provide
	// the payload explicitly. Alternatively, a message might have a canonical
	// serialization that can always be unambiguously computed to derive the
	// payload.
	Signature []byte `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	// The identifier for the public key that verifies this signature.
	//   * The `public_key_id` is required.
	//   * The `public_key_id` MUST be an RFC3986 conformant URI.
	//   * When possible, the `public_key_id` SHOULD be an immutable reference,
	//     such as a cryptographic digest.
	//
	// Examples of valid `public_key_id`s:
	//
	// OpenPGP V4 public key fingerprint:
	//   * "openpgp4fpr:74FAF3B861BDA0870C7B6DEF607E48D2A663AEEA"
	// See https://www.iana.org/assignments/uri-schemes/prov/openpgp4fpr for more
	// details on this scheme.
	//
	// RFC6920 digest-named SubjectPublicKeyInfo (digest of the DER
	// serialization):
	//   * "ni:///sha-256;cD9o9Cq6LG3jD0iKXqEi_vdjJGecm_iXkbqVoScViaU"
	//   * "nih:///sha-256;703f68f42aba2c6de30f488a5ea122fef76324679c9bf89791ba95a1271589a5"
	PublicKeyId          string   `protobuf:"bytes,2,opt,name=public_key_id,json=publicKeyId,proto3" json:"public_key_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Signature) Reset()         { *m = Signature{} }
func (m *Signature) String() string { return proto.CompactTextString(m) }
func (*Signature) ProtoMessage()    {}
func (*Signature) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{1}
}

func (m *Signature) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Signature.Unmarshal(m, b)
}
func (m *Signature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Signature.Marshal(b, m, deterministic)
}
func (m *Signature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Signature.Merge(m, src)
}
func (m *Signature) XXX_Size() int {
	return xxx_messageInfo_Signature.Size(m)
}
func (m *Signature) XXX_DiscardUnknown() {
	xxx_messageInfo_Signature.DiscardUnknown(m)
}

var xxx_messageInfo_Signature proto.InternalMessageInfo

func (m *Signature) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *Signature) GetPublicKeyId() string {
	if m != nil {
		return m.PublicKeyId
	}
	return ""
}

func init() {
	proto.RegisterEnum("grafeas.v1beta1.NoteKind", NoteKind_name, NoteKind_value)
	proto.RegisterType((*RelatedUrl)(nil), "grafeas.v1beta1.RelatedUrl")
	proto.RegisterType((*Signature)(nil), "grafeas.v1beta1.Signature")
}

func init() { proto.RegisterFile("common.proto", fileDescriptor_555bd8c177793206) }

var fileDescriptor_555bd8c177793206 = []byte{
	// 328 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xc1, 0x4e, 0xc2, 0x40,
	0x10, 0x86, 0x2d, 0x08, 0xd8, 0x01, 0x64, 0xdd, 0x68, 0x82, 0x89, 0x07, 0xc3, 0xc9, 0x78, 0x28,
	0x21, 0x7a, 0xf0, 0x5a, 0xe8, 0x4a, 0x36, 0x2d, 0xdb, 0xa6, 0xdd, 0x92, 0xe0, 0xa5, 0x69, 0xe9,
	0x5a, 0x1b, 0x0b, 0x4b, 0x4a, 0x6b, 0xc2, 0x33, 0xf8, 0x16, 0x3e, 0xa9, 0xa1, 0xa0, 0x26, 0x9e,
	0xf6, 0xff, 0xbe, 0x9d, 0x99, 0xc3, 0x0f, 0x9d, 0xa5, 0x5c, 0xad, 0xe4, 0x5a, 0xdb, 0xe4, 0xb2,
	0x90, 0xb8, 0x97, 0xe4, 0xe1, 0xab, 0x08, 0xb7, 0xda, 0xc7, 0x28, 0x12, 0x45, 0x38, 0x1a, 0x3c,
	0x02, 0xb8, 0x22, 0x0b, 0x0b, 0x11, 0xfb, 0x79, 0x86, 0x11, 0xd4, 0xcb, 0x3c, 0xeb, 0x2b, 0xb7,
	0xca, 0x9d, 0xea, 0xee, 0x23, 0xbe, 0x84, 0x46, 0x16, 0x46, 0x22, 0xeb, 0xd7, 0x2a, 0x77, 0x80,
	0xc1, 0x0c, 0x54, 0x2f, 0x4d, 0xd6, 0x61, 0x51, 0xe6, 0x02, 0xdf, 0x80, 0xba, 0xfd, 0x81, 0x6a,
	0xb5, 0xe3, 0xfe, 0x09, 0x3c, 0x80, 0xee, 0xa6, 0x8c, 0xb2, 0x74, 0x19, 0xbc, 0x8b, 0x5d, 0x90,
	0xc6, 0xc7, 0x43, 0xed, 0x83, 0x34, 0xc5, 0x8e, 0xc6, 0xf7, 0x9f, 0x0a, 0x9c, 0x31, 0x59, 0x08,
	0x33, 0x5d, 0xc7, 0xf8, 0x1a, 0xae, 0x98, 0xcd, 0x49, 0x60, 0x52, 0x66, 0x04, 0x3e, 0xf3, 0x1c,
	0x32, 0xa1, 0xcf, 0x94, 0x18, 0xe8, 0x04, 0x5f, 0x40, 0x77, 0xee, 0x5b, 0x8c, 0xb8, 0xfa, 0x98,
	0x5a, 0x94, 0x2f, 0x90, 0x82, 0x55, 0x68, 0x8c, 0x7d, 0x6a, 0x19, 0xa8, 0xb6, 0x8f, 0x74, 0xa6,
	0x4f, 0x09, 0xaa, 0xe3, 0x36, 0xb4, 0x1c, 0x7d, 0x62, 0xee, 0xe1, 0x14, 0x9f, 0x03, 0x18, 0xc4,
	0xb1, 0xec, 0xc5, 0x8c, 0x30, 0x8e, 0x1a, 0xb8, 0x0b, 0xaa, 0x41, 0xbd, 0x89, 0x3d, 0x27, 0xee,
	0x02, 0x35, 0x71, 0x0f, 0xda, 0x3a, 0xe7, 0xc4, 0xe3, 0x3a, 0xa7, 0x36, 0x43, 0xad, 0xf1, 0x1c,
	0x70, 0x2a, 0xb5, 0x7f, 0x45, 0x39, 0xca, 0xcb, 0x53, 0x92, 0x16, 0x6f, 0x65, 0xa4, 0x2d, 0xe5,
	0x6a, 0x78, 0xfc, 0xfd, 0x7d, 0xab, 0x76, 0x87, 0xc7, 0xd9, 0xe1, 0xa1, 0xf2, 0x20, 0x91, 0x41,
	0xe5, 0xbf, 0x6a, 0xf5, 0xa9, 0xab, 0x47, 0xcd, 0x0a, 0x1e, 0xbe, 0x03, 0x00, 0x00, 0xff, 0xff,
	0xd9, 0xa1, 0x36, 0xff, 0x92, 0x01, 0x00, 0x00,
}
