package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/VirgilSecurity/virgil-services-cards/src/model"
)

//
// CSR stamp type checkers :: for application signature :: behaves properly.
//
func TestCSRStampTypeCheckersForApplicationContentStamp(t *testing.T) {

	contentStamp := CSRStamp{
		Signer: model.ApplicationSignatureType,
	}

	assert.False(t, contentStamp.IsSelf())
	assert.True(t, contentStamp.IsApplication())
	assert.False(t, contentStamp.IsVirgil())
}

//
// CSR stamp type checkers :: for self signature :: behaves properly.
//
func TestCSRStampTypeCheckersForSelfContentStamp(t *testing.T) {

	contentStamp := CSRStamp{
		Signer: model.SelfSignatureType,
	}

	assert.True(t, contentStamp.IsSelf())
	assert.False(t, contentStamp.IsApplication())
	assert.False(t, contentStamp.IsVirgil())
}

//
// CSR stamp type checkers :: for Virgil signature :: behaves properly.
//
func TestCSRStampTypeCheckersForVirgilContentStamp(t *testing.T) {

	contentStamp := CSRStamp{
		Signer: model.VirgilSignatureType,
	}

	assert.False(t, contentStamp.IsSelf())
	assert.False(t, contentStamp.IsApplication())
	assert.True(t, contentStamp.IsVirgil())
}

//
// Getters :: for a configured object :: returns set values.
//
func TestCSRStampTypeCheckersForVirgilSignature(t *testing.T) {

	signer := "Signer type"
	signerSnapshot := "Signer snapshot"
	signerSignature := "Signer signature"
	stamp := CSRStamp{
		Signer:    signer,
		Snapshot:  signerSnapshot,
		Signature: signerSignature,
	}

	assert.Equal(t, signer, stamp.GetSigner())
	assert.Equal(t, signerSnapshot, stamp.GetSnapshot())
	assert.Equal(t, signerSignature, stamp.GetSignature())
}
