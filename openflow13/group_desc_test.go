package openflow13

import (
	"testing"

	"github.com/contiv/libOpenflow/util"
	"github.com/stretchr/testify/require"
)

func TestGroupDescStatsRoundTrip(t *testing.T) {
	bkt := NewBucket()
	desc := &GroupDescStats{
		Type:    OFPGT_ALL,
		GroupId: 0x20000001,
		Buckets: []Bucket{*bkt},
	}
	raw, err := desc.MarshalBinary()
	require.NoError(t, err)

	got := &GroupDescStats{}
	require.NoError(t, got.UnmarshalBinary(raw))
	require.Equal(t, desc.GroupId, got.GroupId)
	require.Equal(t, desc.Type, got.Type)
	require.Equal(t, uint16(len(raw)), got.Len())
}

func TestMultipartReplyGroupDesc(t *testing.T) {
	first := &GroupDescStats{Type: OFPGT_ALL, GroupId: 0x10000000}
	second := &GroupDescStats{Type: OFPGT_ALL, GroupId: 0x01000000}
	reply := &MultipartReply{
		Header: NewOfp13Header(),
		Type:   MultipartType_GroupDesc,
		Body:   []util.Message{first, second},
	}
	reply.Header.Type = Type_MultiPartReply
	raw, err := reply.MarshalBinary()
	require.NoError(t, err)

	got := &MultipartReply{}
	require.NoError(t, got.UnmarshalBinary(raw))
	require.Equal(t, uint16(MultipartType_GroupDesc), got.Type)
	require.Len(t, got.Body, 2)
	require.Equal(t, uint32(0x10000000), got.Body[0].(*GroupDescStats).GroupId)
	require.Equal(t, uint32(0x01000000), got.Body[1].(*GroupDescStats).GroupId)
}
