package encoder

import "encoding/base64"

//
// Base64 is a base64-encoder.
//
type Base64 struct{}

//
// NewBase64 returns a new Base64 encoder instance.
//
func NewBase64() Base64 {

	return Base64{}
}

//
// DecodeString decodes string message.
//
func (e Base64) DecodeString(data string) ([]byte, error) {

	return base64.StdEncoding.DecodeString(data)
}

//
// EncodeToString encodes a message to string.
//
func (e Base64) EncodeToString(data []byte) string {

	return base64.StdEncoding.EncodeToString(data)
}
