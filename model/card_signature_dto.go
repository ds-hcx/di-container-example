package model

//
// Signature types.
//
const (
	SelfSignatureType        = "self"
	ApplicationSignatureType = "app"
	VirgilSignatureType      = "virgil"
)

//
// CardSignatureDTO represents the signature entry for the Virgil Card.
//
type CardSignatureDTO struct {
	Signer    string `json:"signer"`
	Snapshot  string `json:"snapshot,omitempty"`
	Signature string `json:"signature"`
}

//
// GetSigner returns signer value.
//
func (cs *CardSignatureDTO) GetSigner() string {

	return cs.Signer
}

//
// GetExtraContent returns extra content value.
//
func (cs *CardSignatureDTO) GetExtraContent() string {

	return cs.Snapshot
}

//
// GetSignature returns signature value.
//
func (cs *CardSignatureDTO) GetSignature() string {

	return cs.Signature
}

//
// IsApp returns true if signature is application one.
//
func (cs *CardSignatureDTO) IsApp() bool {

	return ApplicationSignatureType == cs.Signer
}

//
// IsVirgil returns true if signature is Virgil one.
//
func (cs *CardSignatureDTO) IsVirgil() bool {

	return VirgilSignatureType == cs.Signer
}
