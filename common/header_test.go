package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHelloElemVersionBitmap(t *testing.T) {
	tests := []struct {
		name     string
		bitmap   *HelloElemVersionBitmap
		expected []uint32
	}{
		{
			name:     "default advertises OpenFlow 1.0 and 1.3",
			bitmap:   NewHelloElemVersionBitmap(),
			expected: []uint32{1<<1 | 1<<4},
		},
		{
			name:     "single version",
			bitmap:   NewHelloElemVersionBitmapWithVersions(6),
			expected: []uint32{1 << 6},
		},
		{
			name:     "multiple versions",
			bitmap:   NewHelloElemVersionBitmapWithVersions(1, 4, 6),
			expected: []uint32{1<<1 | 1<<4 | 1<<6},
		},
		{
			name:     "duplicate versions",
			bitmap:   NewHelloElemVersionBitmapWithVersions(6, 6),
			expected: []uint32{1 << 6},
		},
		{
			name:     "version 0 is ignored",
			bitmap:   NewHelloElemVersionBitmapWithVersions(0, 6),
			expected: []uint32{1 << 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.bitmap.Bitmaps)
			require.Equal(t, uint16(4+4*len(tt.expected)), tt.bitmap.Len())
		})
	}
}

func TestHelloElemVersionBitmapMarshalRoundTrip(t *testing.T) {
	bitmap := NewHelloElemVersionBitmapWithVersions(4, 6)
	data, err := bitmap.MarshalBinary()
	require.NoError(t, err)
	require.Len(t, data, int(bitmap.Len()))

	var decoded HelloElemVersionBitmap
	require.NoError(t, decoded.UnmarshalBinary(data))
	require.Equal(t, uint16(HelloElemType_VersionBitmap), decoded.Type)
	require.Equal(t, bitmap.Bitmaps, decoded.Bitmaps)
}

func TestNewHelloWithVersions(t *testing.T) {
	hello, err := NewHelloWithVersions(1, 4)
	require.NoError(t, err)
	require.Equal(t, uint8(4), hello.Header.Version)
	require.Len(t, hello.Elements, 1)

	bitmap, ok := hello.Elements[0].(*HelloElemVersionBitmap)
	require.True(t, ok)
	require.Equal(t, []uint32{1<<1 | 1<<4}, bitmap.Bitmaps)
}

func TestNewHelloWithVersionsRequiresVersion(t *testing.T) {
	_, err := NewHelloWithVersions()
	require.Error(t, err)

	_, err = NewHelloWithVersions(0)
	require.Error(t, err)
}

func TestNewHelloKeepsDefaultBitmap(t *testing.T) {
	hello, err := NewHello(4)
	require.NoError(t, err)
	require.Equal(t, uint8(4), hello.Header.Version)

	bitmap, ok := hello.Elements[0].(*HelloElemVersionBitmap)
	require.True(t, ok)
	require.Equal(t, []uint32{1<<1 | 1<<4}, bitmap.Bitmaps)
}
