package api

import "github.com/VirgilSecurity/virgil-services-cards/src/model"

//
// CSRStamp is an object that proves CSR validity.
// CSRStamps are expected to come from the Virgil Card issuer (type=self), 3rd-party Application services (type=app)
// and from the Virgil Cards Service (type=virgil).
//
type CSRStamp struct {
	Signer    string `json:"signer"`
	Signature string `json:"signature"`
	Snapshot  string `json:"snapshot,omitempty"`
}

//
// GetSigner returns signer value.
//
func (s *CSRStamp) GetSigner() string {

	return s.Signer
}

//
// GetSignature returns signature value.
//
func (s *CSRStamp) GetSignature() string {

	return s.Signature
}

//
// GetSnapshot returns a snapshot value.
//
func (s *CSRStamp) GetSnapshot() string {

	return s.Snapshot
}

//
// IsSelf returns true if CSRStamp is a self-signed one.
//
func (s *CSRStamp) IsSelf() bool {

	return model.SelfSignatureType == s.Signer
}

//
// IsApplication returns true if CSRStamp is an application-signed one.
//
func (s *CSRStamp) IsApplication() bool {

	return model.ApplicationSignatureType == s.Signer
}

//
// IsVirgil returns true if CSRStamp a Virgil-signed one.
//
func (s *CSRStamp) IsVirgil() bool {

	return model.VirgilSignatureType == s.Signer
}
