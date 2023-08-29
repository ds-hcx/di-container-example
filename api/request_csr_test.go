package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//
// CSR getters :: on a configured object :: return preset values.
//
func TestCSRGetters(t *testing.T) {

	publicKey := "public key"
	previousCardID := "previous Card ID"
	identity := "identity"
	version := "version"
	createdAt := time.Now().UTC().Unix()
	csr := CSR{
		PublicKey:      publicKey,
		PreviousCardID: previousCardID,
		Identity:       identity,
		Version:        version,
		CreatedAt:      createdAt,
	}

	assert.Equal(t, publicKey, csr.GetPublicKey())
	assert.Equal(t, previousCardID, csr.GetPreviousCardID())
	assert.Equal(t, identity, csr.GetIdentity())
	assert.Equal(t, version, csr.GetVersion())
	assert.Equal(t, createdAt, csr.GetCreatedAt())
}
