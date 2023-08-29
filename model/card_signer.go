package model

import (
	"encoding/base64"

	"github.com/VirgilSecurity/virgil-services-core-kit/errors"
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"

	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto"
)

//
// CardSigner is an interface.
//
type CardSigner interface {
	//
	// SignCardByCardsService signs the Virgil Card with VirgilCards service.
	//
	SignCardByCardsService(tracer.Span, *CardDTO) error
}

//
// DefaultSigner is a default Virgil Cards signer interface.
//
type DefaultSigner struct {
	crypto           crypto.Provider
	cards5PrivateKey crypto.PrivateKey
	cards5SignerID   string
}

//
// NewSigner returns a new instance of card signer.
//
func NewSigner(
	crypto crypto.Provider,
	cards5SignerID string,
	cards5PrivateKey crypto.PrivateKey,
) DefaultSigner {

	return DefaultSigner{
		crypto:           crypto,
		cards5SignerID:   cards5SignerID,
		cards5PrivateKey: cards5PrivateKey,
	}
}

//
// SignCardByCardsService signs the Virgil Card with VirgilCards service.
//
func (cs DefaultSigner) SignCardByCardsService(span tracer.Span, c *CardDTO) (err error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentModel,
		},
	)
	defer span.Finish()

	cardContentSnapshotBytes, err := base64.StdEncoding.DecodeString(c.GetContentSnapshot())
	if nil != err {
		return tracer.SetSpanErrorAndReturn(span, errors.Wrap(err, errors.New(
			"virgil card content snapshot (%s) decode error", c.GetContentSnapshot()),
		))
	}

	signature, err := cs.crypto.SignVirgilCard(cardContentSnapshotBytes, []byte{}, cs.cards5PrivateKey)
	if nil != err {
		return tracer.SetSpanErrorAndReturn(span, errors.Wrap(err, errors.New(
			"virgil card content snapshot (%s) sign error", c.GetContentSnapshot()),
		))
	}

	c.AppendSignature(&CardSignatureDTO{
		Signer:    VirgilSignatureType,
		Signature: base64.StdEncoding.EncodeToString(signature),
	})

	return nil
}
