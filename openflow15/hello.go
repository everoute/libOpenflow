package openflow15

import "github.com/contiv/libOpenflow/common"

// NewHelloElemVersionBitmap returns a Hello version bitmap element that
// advertises OpenFlow 1.5 only.
func NewHelloElemVersionBitmap() *common.HelloElemVersionBitmap {
	return common.NewHelloElemVersionBitmapWithVersions(VERSION)
}

// NewHello returns an OpenFlow 1.5 Hello message whose version bitmap
// advertises OpenFlow 1.5 only.
func NewHello() (*common.Hello, error) {
	return common.NewHelloWithVersions(VERSION)
}
