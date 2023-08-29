package model

import (
	"encoding/base64"
	"time"
)

//
// Card versions
//
const (
	CardVersion5 = "5.0"
)

//
// CardDTO represents the Virgil Card object persisted in the database.
//
type CardDTO struct {
	ID              string              `json:"-"`
	ContentSnapshot string              `json:"content_snapshot"`
	Identity        string              `json:"-"`
	PreviousCardID  string              `json:"-"`
	ApplicationID   string              `json:"-"`
	Version         string              `json:"-"`
	ChainID         string              `json:"-"`
	CreatedAt       int64               `json:"-"`
	Signatures      []*CardSignatureDTO `json:"signatures"`
	PublicKey       []byte              `json:"-"`
	IsSuperseeded   bool                `json:"-"`
}

//
// NewCardDTO returns a new Virgil CArd DTO instance.
//
func NewCardDTO() *CardDTO {

	return &CardDTO{}
}

//
// GetID returns an ID value.
//
func (c *CardDTO) GetID() string {

	return c.ID
}

//
// GetContentSnapshot returns content snapshot value.
//
func (c *CardDTO) GetContentSnapshot() string {

	return c.ContentSnapshot
}

//
// GetVersion returns card version value.
//
func (c *CardDTO) GetVersion() string {

	return c.Version
}

//
// GetApplicationID returns an application ID.
//
func (c *CardDTO) GetApplicationID() string {

	return c.ApplicationID
}

//
// GetIdentity returns an Identity value.
//
func (c *CardDTO) GetIdentity() string {

	return c.Identity
}

//
// GetEncodedPublicKey returns an encoded public key value.
//
func (c *CardDTO) GetEncodedPublicKey() string {

	return base64.StdEncoding.EncodeToString(c.PublicKey)
}

//
// GetPreviousCardID returns previous card id.
//
func (c *CardDTO) GetPreviousCardID() string {

	return c.PreviousCardID
}

//
// GetPublicKey returns a public key value.
//
func (c *CardDTO) GetPublicKey() []byte {

	return c.PublicKey
}

//
// GetSignatures returns signatures collection.
//
func (c *CardDTO) GetSignatures() []*CardSignatureDTO {

	return c.Signatures
}

//
// AppendSignature appends new signature value.
//
func (c *CardDTO) AppendSignature(signature *CardSignatureDTO) {

	c.Signatures = append(c.Signatures, signature)
}

//
// GetCreatedAt returns creation time.
//
func (c *CardDTO) GetCreatedAt() time.Time {
	return time.Unix(c.CreatedAt, 0)
}

//
// SetCreatedAt sets creation time.
//
func (c *CardDTO) SetCreatedAt(created time.Time) {
	c.CreatedAt = created.Unix()
}

//
// GetChainID returns card's chain id.
//
func (c *CardDTO) GetChainID() string {

	return c.ChainID
}

//
// SetChainID sets card's chain id.
//
func (c *CardDTO) SetChainID(chainID string) {
	c.ChainID = chainID
}

//
// DoesScopeMatch returns true if Virgil Card application ID matches the authorization scope application IDs.
//
func (c *CardDTO) DoesScopeMatch(scopeID string) bool {

	return scopeID == c.GetApplicationID()
}
