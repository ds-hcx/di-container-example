package api

//
// CSR fields constants.
//
const (
	CreatedAtFieldName = "created_at"
)

//
// CSR is a Card Signing Request message.
// It gets passed an a base64-encoded JSON in the request content snapshot.
//
type CSR struct {
	PublicKey      string `json:"public_key,omitempty"`
	Identity       string `json:"identity,omitempty"`
	PreviousCardID string `json:"previous_card_id,omitempty"`
	Version        string `json:"version,omitempty"`
	CreatedAt      int64  `json:"created_at,omitempty"`
}

//
// GetPublicKey returns public key value.
//
func (m *CSR) GetPublicKey() string {

	return m.PublicKey
}

//
// GetIdentity returns identity value.
//
func (m *CSR) GetIdentity() string {

	return m.Identity
}

//
// GetPreviousCardID returns previous Virgil Card ID value.
//
func (m *CSR) GetPreviousCardID() string {

	return m.PreviousCardID
}

//
// GetVersion returns Virgil Card version.
//
func (m *CSR) GetVersion() string {

	return m.Version
}

//
// GetCreatedAt returns Virgil Card creation UTC Unix timestamp.
//
func (m *CSR) GetCreatedAt() int64 {

	return m.CreatedAt
}
