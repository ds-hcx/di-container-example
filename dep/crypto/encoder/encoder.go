package encoder

//
// Provider is a request encoder.
//
type Provider interface {
	//
	// DecodeString decodes string message.
	//
	DecodeString(string) ([]byte, error)

	//
	// EncodeToString encodes a message to string.
	//
	EncodeToString([]byte) string
}
