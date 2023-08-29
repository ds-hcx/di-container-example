package encoder

import (
	"encoding/hex"
)

//
// Hex is a hexadecimal encoder implementation.
//
type Hex struct{}

//
// NewHex returns a hex encoder object.
//
func NewHex() *Hex {

	return &Hex{}
}

//
// EncodeToString encodes a message to string.
//
func (h *Hex) EncodeToString(data []byte) string {

	return hex.EncodeToString(data)
}

//
// DecodeString decodes string message.
//
func (h *Hex) DecodeString(data string) ([]byte, error) {

	return hex.DecodeString(data)
}
