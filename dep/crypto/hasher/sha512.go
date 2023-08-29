package hasher

import (
	"crypto/sha512"
)

//
// Provider is a hashing algorithm interface.
//
type Provider interface {
	//
	// Hash calculates hash for a data.
	//
	Hash([]byte) []byte
}

//
// SHA512 hasher object.
//
type SHA512 struct{}

//
// NewSHA512 returns an instance of SHA512 hasher.
//
func NewSHA512() *SHA512 {

	return &SHA512{}
}

//
// Hash calculates hash for a data.
//
func (h *SHA512) Hash(data []byte) []byte {

	hash := sha512.Sum512(data)

	return hash[:]
}
