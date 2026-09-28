package openflow15

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/contiv/libOpenflow/common"
)

func TestNewHelloElementBitmap(t *testing.T) {
	bitmap := NewHelloElemVersionBitmap()
	require.Equal(t, []uint32{1 << VERSION}, bitmap.Bitmaps)
}

func TestNewHello(t *testing.T) {
	hello, err := NewHello()
	require.NoError(t, err)
	require.Equal(t, uint8(VERSION), hello.Header.Version)
	require.Len(t, hello.Elements, 1)

	bitmap, ok := hello.Elements[0].(*common.HelloElemVersionBitmap)
	require.True(t, ok)
	require.Equal(t, []uint32{1 << VERSION}, bitmap.Bitmaps)

	data, err := hello.MarshalBinary()
	require.NoError(t, err)

	var decoded common.Hello
	require.NoError(t, decoded.UnmarshalBinary(data))
	require.Equal(t, uint8(VERSION), decoded.Header.Version)
	require.Len(t, decoded.Elements, 1)

	decodedBitmap, ok := decoded.Elements[0].(*common.HelloElemVersionBitmap)
	require.True(t, ok)
	require.Equal(t, []uint32{1 << VERSION}, decodedBitmap.Bitmaps)
}
