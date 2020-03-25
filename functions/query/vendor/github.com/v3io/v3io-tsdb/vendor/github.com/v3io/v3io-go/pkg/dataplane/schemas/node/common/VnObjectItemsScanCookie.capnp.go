// Code generated by capnpc-go. DO NOT EDIT.

package node_common_capnp

import (
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
)

type VnObjectItemsScanCookie struct{ capnp.Struct }

// VnObjectItemsScanCookie_TypeID is the unique identifier for the type VnObjectItemsScanCookie.
const VnObjectItemsScanCookie_TypeID = 0x8402081dabc79360

func NewVnObjectItemsScanCookie(s *capnp.Segment) (VnObjectItemsScanCookie, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 0})
	return VnObjectItemsScanCookie{st}, err
}

func NewRootVnObjectItemsScanCookie(s *capnp.Segment) (VnObjectItemsScanCookie, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 0})
	return VnObjectItemsScanCookie{st}, err
}

func ReadRootVnObjectItemsScanCookie(msg *capnp.Message) (VnObjectItemsScanCookie, error) {
	root, err := msg.RootPtr()
	return VnObjectItemsScanCookie{root.Struct()}, err
}

func (s VnObjectItemsScanCookie) String() string {
	str, _ := text.Marshal(0x8402081dabc79360, s.Struct)
	return str
}

func (s VnObjectItemsScanCookie) SliceId() uint16 {
	return s.Struct.Uint16(0)
}

func (s VnObjectItemsScanCookie) SetSliceId(v uint16) {
	s.Struct.SetUint16(0, v)
}

func (s VnObjectItemsScanCookie) InodeNumber() uint32 {
	return s.Struct.Uint32(4)
}

func (s VnObjectItemsScanCookie) SetInodeNumber(v uint32) {
	s.Struct.SetUint32(4, v)
}

func (s VnObjectItemsScanCookie) ClientSliceListPos() uint64 {
	return s.Struct.Uint64(8)
}

func (s VnObjectItemsScanCookie) SetClientSliceListPos(v uint64) {
	s.Struct.SetUint64(8, v)
}

func (s VnObjectItemsScanCookie) ClientSliceListEndPos() uint64 {
	return s.Struct.Uint64(16)
}

func (s VnObjectItemsScanCookie) SetClientSliceListEndPos(v uint64) {
	s.Struct.SetUint64(16, v)
}

// VnObjectItemsScanCookie_List is a list of VnObjectItemsScanCookie.
type VnObjectItemsScanCookie_List struct{ capnp.List }

// NewVnObjectItemsScanCookie creates a new list of VnObjectItemsScanCookie.
func NewVnObjectItemsScanCookie_List(s *capnp.Segment, sz int32) (VnObjectItemsScanCookie_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 0}, sz)
	return VnObjectItemsScanCookie_List{l}, err
}

func (s VnObjectItemsScanCookie_List) At(i int) VnObjectItemsScanCookie {
	return VnObjectItemsScanCookie{s.List.Struct(i)}
}

func (s VnObjectItemsScanCookie_List) Set(i int, v VnObjectItemsScanCookie) error {
	return s.List.SetStruct(i, v.Struct)
}

func (s VnObjectItemsScanCookie_List) String() string {
	str, _ := text.MarshalList(0x8402081dabc79360, s.List)
	return str
}

// VnObjectItemsScanCookie_Promise is a wrapper for a VnObjectItemsScanCookie promised by a client call.
type VnObjectItemsScanCookie_Promise struct{ *capnp.Pipeline }

func (p VnObjectItemsScanCookie_Promise) Struct() (VnObjectItemsScanCookie, error) {
	s, err := p.Pipeline.Struct()
	return VnObjectItemsScanCookie{s}, err
}

const schema_b56ec2d13b48b7cb = "x\xdal\xce\xb1J+A\x1c\x85\xf1s\xfe\x93\xdcI" +
	"\xba\xccMZ\xd1\xda\xc2 v\xda\x88\"\x18\x90\x98!" +
	"`ea2;\xe0hvfq\xd7B\x08\xd8X\xd8" +
	"\xdb\xf9\x04\x82\x9d \xf66)\xac|\x0a\x1fc%\x95" +
	"M\xda\xef\xd7|\x9d\x97}\xd9n.\x04\xb0\x1b\xcd\x7f" +
	"\xf5\xc5\xd3\xe2u\xad%\x0f\xb0\x9bT\xf5\xd7\xc7\xf1\xde" +
	"\xf7g|GC\x03;\x03\xfe\x97n\xa0\x06\xba\x9e?" +
	"\x18\xd61e\xbe\xefR\xae\xf3\x14\xfbg\xf1tz\xe5" +
	"]5\xa8|^\x8e\xdd$\x1e\xa6t\x1d\xfc\x96\x9b\x14" +
	"\xb1\xd8]\xad\xf4#\xd2vT\x03h\x100\x93\x03\xc0" +
	"\x9e+\xdaK\xa1!{\\F?\x05l\xa6h\x0b\xa1" +
	"\x11\xf6(\x80\xc9\x9f\x01[(\xda\xb9\xd0(\xe9Q\x01" +
	"\xe6\xee\x0d\xb0sE\xfb(\xbc/g\xc1\xf9AF\x0d" +
	"\xa1\x06\xeb\xb0\xbc\x1d\xde\xe6\xd0S\x7f\xc3\x16\x84-\xb0" +
	"v\xb3\xe0c5\x9e18\x7f\x12\xcaj\xa4R\xc96" +
	"\x84\xed\x15x\xb4\x1e\xb3\xd1\x9f\xff\x06\x00\x00\xff\xff\xbc" +
	"\xd4O\xa2"

func init() {
	schemas.Register(schema_b56ec2d13b48b7cb,
		0x8402081dabc79360)
}
