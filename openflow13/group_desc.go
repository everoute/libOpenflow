package openflow13

import (
	"encoding/binary"
	"fmt"

	"github.com/contiv/libOpenflow/util"
)

// emptyMultipartBody is the empty body of a group description request.
type emptyMultipartBody struct{}

func (e *emptyMultipartBody) Len() uint16 { return 0 }

func (e *emptyMultipartBody) MarshalBinary() ([]byte, error) { return nil, nil }

func (e *emptyMultipartBody) UnmarshalBinary([]byte) error { return nil }

// NewGroupDescRequest builds an OF1.3 group description multipart request.
// The request body is empty.
func NewGroupDescRequest() *MultipartRequest {
	req := new(MultipartRequest)
	req.Header = NewOfp13Header()
	req.Header.Type = Type_MultiPartRequest
	req.Type = MultipartType_GroupDesc
	req.Body = &emptyMultipartBody{}
	return req
}

// GroupDescStats is an ofp_group_desc entry.
type GroupDescStats struct {
	Length  uint16
	Type    uint8
	pad     uint8
	GroupId uint32
	Buckets []Bucket
}

func (g *GroupDescStats) Len() uint16 {
	if g.Length != 0 {
		return g.Length
	}
	n := uint16(8)
	for i := range g.Buckets {
		n += g.Buckets[i].Len()
	}
	return n
}

func (g *GroupDescStats) MarshalBinary() ([]byte, error) {
	var buckets []byte
	for i := range g.Buckets {
		raw, err := g.Buckets[i].MarshalBinary()
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, raw...)
	}
	g.Length = uint16(8 + len(buckets))
	data := make([]byte, g.Length)
	binary.BigEndian.PutUint16(data[0:], g.Length)
	data[2] = g.Type
	binary.BigEndian.PutUint32(data[4:], g.GroupId)
	copy(data[8:], buckets)
	return data, nil
}

func (g *GroupDescStats) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return fmt.Errorf("group description is too short: %d", len(data))
	}
	g.Length = binary.BigEndian.Uint16(data[0:])
	g.Type = data[2]
	g.pad = data[3]
	g.GroupId = binary.BigEndian.Uint32(data[4:])
	if g.Length < 8 || int(g.Length) > len(data) {
		return fmt.Errorf("invalid group description length %d", g.Length)
	}

	n := 8
	g.Buckets = nil
	for n < int(g.Length) {
		bkt := new(Bucket)
		if err := bkt.UnmarshalBinary(data[n:]); err != nil {
			return err
		}
		if bkt.Length == 0 || n+int(bkt.Length) > int(g.Length) {
			return fmt.Errorf("invalid bucket length %d", bkt.Length)
		}
		g.Buckets = append(g.Buckets, *bkt)
		n += int(bkt.Length)
	}
	return nil
}

var _ util.Message = &emptyMultipartBody{}
