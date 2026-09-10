package openflow13

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIpTtlFieldRoundTrip(t *testing.T) {
	field := NewIpTtlField(1)
	data, err := field.MarshalBinary()
	require.NoError(t, err)
	assert.Equal(t, 5, len(data))

	decoded := new(MatchField)
	require.NoError(t, decoded.UnmarshalBinary(data))
	assert.Equal(t, uint16(OXM_CLASS_NXM_1), decoded.Class)
	assert.Equal(t, uint8(NXM_NX_IP_TTL), decoded.Field)
	assert.Equal(t, uint16(1), decoded.Value.Len())
	ttl, ok := decoded.Value.(*Uint8Message)
	assert.True(t, ok)
	assert.Equal(t, uint8(1), ttl.Data)
}

func TestIpEcnFieldRoundTrip(t *testing.T) {
	field := NewIpEcnField(3)
	data, err := field.MarshalBinary()
	require.NoError(t, err)
	assert.Equal(t, 5, len(data))

	decoded := new(MatchField)
	require.NoError(t, decoded.UnmarshalBinary(data))
	assert.Equal(t, uint8(NXM_NX_IP_ECN), decoded.Field)
	ecn, ok := decoded.Value.(*Uint8Message)
	assert.True(t, ok)
	assert.Equal(t, uint8(3), ecn.Data)
	assert.Equal(t, uint16(1), decoded.Value.Len())
}

func TestMplsTtlFieldRoundTrip(t *testing.T) {
	field := NewMplsTtlField(64)
	data, err := field.MarshalBinary()
	require.NoError(t, err)
	assert.Equal(t, 5, len(data))

	decoded := new(MatchField)
	require.NoError(t, decoded.UnmarshalBinary(data))
	assert.Equal(t, uint8(NXM_NX_MPLS_TTL), decoded.Field)
	ttl, ok := decoded.Value.(*Uint8Message)
	assert.True(t, ok)
	assert.Equal(t, uint8(64), ttl.Data)
	assert.Equal(t, uint16(1), decoded.Value.Len())
}

func TestMatchKeepsAlignmentAfterIpTtl(t *testing.T) {
	match := NewMatch()
	match.AddField(*NewEthTypeField(0x0800))
	match.AddField(*NewIpTtlField(1))
	match.AddField(*NewEthTypeField(0x0800))

	data, err := match.MarshalBinary()
	require.NoError(t, err)

	decoded := NewMatch()
	require.NoError(t, decoded.UnmarshalBinary(data))
	require.Len(t, decoded.Fields, 3)
	assert.Equal(t, uint8(NXM_NX_IP_TTL), decoded.Fields[1].Field)
	assert.Equal(t, uint16(1), decoded.Fields[1].Value.Len())
	assert.Equal(t, uint8(OXM_FIELD_ETH_TYPE), decoded.Fields[2].Field)
}

func TestTunnelIpv4SrcStillFourBytes(t *testing.T) {
	field := NewTunnelIpv4SrcField(net.ParseIP("10.0.0.1"), nil)
	data, err := field.MarshalBinary()
	require.NoError(t, err)
	assert.Equal(t, 8, len(data))

	decoded := new(MatchField)
	require.NoError(t, decoded.UnmarshalBinary(data))
	assert.Equal(t, uint8(NXM_NX_TUN_IPV4_SRC), decoded.Field)
	src, ok := decoded.Value.(*TunnelIpv4SrcField)
	require.True(t, ok)
	assert.Equal(t, uint16(4), src.Len())
	assert.Equal(t, "10.0.0.1", src.TunnelIpv4Src.String())
}
