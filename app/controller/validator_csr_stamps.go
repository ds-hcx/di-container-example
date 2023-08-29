package controller

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/errors"
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"

	"github.com/VirgilSecurity/virgil-services-cards/src/api"
	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto"
	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/encoder"
	"github.com/VirgilSecurity/virgil-services-cards/src/model"
)

//
// Validation constraints.
//
const (
	CSRStampsListMinLength = 1
	CSRStampsListMaxLength = 8

	CSRStampSnapshotMaxLength = 1024
	CSRStampSignerMaxLength   = 1024
)

//
// CSRStampsValidatorProvider provides an interface to work with the CSR stamps validator.
//
type CSRStampsValidatorProvider interface {
	//
	// Validate validates the SCR stamps collection.
	//
	Validate(span tracer.Span, scrStamps []api.CSRStamp, scopeID string, csrParams *ParametersStore) (err error)
}

//
// CSRStampsValidator represents the CSR stamps validator.
//
type CSRStampsValidator struct {
	crypto  crypto.Provider
	encoder encoder.Provider
}

//
// NewCSRStampsValidator creates new instance of CSR stamps validator.
//
func NewCSRStampsValidator(crypto crypto.Provider, encoder encoder.Provider) *CSRStampsValidator {
	return &CSRStampsValidator{
		crypto:  crypto,
		encoder: encoder,
	}
}

//
// Validate validates the SCR stamps collection.
//
func (v *CSRStampsValidator) Validate(
	span tracer.Span,
	scrStamps []api.CSRStamp,
	scopeID string,
	csrParams *ParametersStore,
) (err error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentValidator,
		},
	)
	defer span.Finish()

	if CSRStampsListMinLength > len(scrStamps) {
		return tracer.SetSpanErrorAndReturn(span, api.ErrCSRStampsListIsTooSmall)
	}
	if CSRStampsListMaxLength < len(scrStamps) {
		return tracer.SetSpanErrorAndReturn(span, api.ErrCSRStampsListIsTooLarge)
	}

	var doesSelfStampExist bool
	for _, csrStamp := range scrStamps {
		if err = v.validateCSRStamp(span, csrStamp, csrParams); nil != err {
			return
		}
		if csrStamp.IsSelf() {
			if doesSelfStampExist {
				return tracer.SetSpanErrorAndReturn(span, api.ErrSelfCSRStampMustBeUnique)
			}
			doesSelfStampExist = true
		}
	}
	if !doesSelfStampExist {
		return tracer.SetSpanErrorAndReturn(span, api.ErrSelfCSRStampIsMissing)
	}

	return nil
}

//
// validateCSRStamp validates SCR stamp.
//
func (v *CSRStampsValidator) validateCSRStamp(
	span tracer.Span,
	csrStamp api.CSRStamp,
	csrParams *ParametersStore,
) (err error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentValidator,
		},
	)
	defer span.Finish()

	if err = v.validateCSRStampSnapshot(csrStamp.GetSnapshot()); nil != err {
		return tracer.SetSpanErrorAndReturn(span, err)
	}

	if err = v.validateCSRStampSignature(
		span,
		csrStamp.GetSignature(),
		csrStamp.GetSnapshot(),
		csrStamp.IsSelf() && 0 < len(csrParams.publicKey),
		csrParams,
	); nil != err {
		return err
	}

	if err := v.validateCSRSigner(csrStamp.GetSigner()); nil != err {
		return tracer.SetSpanErrorAndReturn(span, err)
	}

	return nil
}

//
// validateCSRStampSignature validates the signature correctness.
//
func (v *CSRStampsValidator) validateCSRStampSignature(
	span tracer.Span,
	signature, snapshot string,
	decodeAndValidateSignature bool,
	csrParams *ParametersStore,
) (err error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentValidator,
		},
	)
	defer span.Finish()

	if 0 == len(signature) {
		return tracer.SetSpanErrorAndReturn(span, api.ErrCSRStampSignatureIsMissing)
	}
	signatureBytes, err := v.encoder.DecodeString(signature)
	if nil != err {
		return tracer.SetSpanErrorAndReturn(span, errors.Wrap(err, api.ErrSignatureDecoding.WithMessage(
			`signature (%s) decode error`, signature,
		)))
	}

	if decodeAndValidateSignature {
		if "" == snapshot {
			csrParams.extraSnapshot = []byte{}
		} else {
			csrParams.extraSnapshot, err = v.encoder.DecodeString(snapshot)
			if err != nil {
				return tracer.SetSpanErrorAndReturn(
					span,
					api.ErrInternalError.WithMessage("impossible to decode snapshot: %+v", err),
				)
			}
		}

		if err = v.crypto.ValidateVirgilCardSignature(
			csrParams.csr,
			csrParams.extraSnapshot,
			csrParams.publicKey,
			signatureBytes,
		); nil != err {
			return tracer.SetSpanErrorAndReturn(
				span,
				api.ErrSignatureVerificationFailed.WithMessage(`self-signature is invalid: %+v`, err),
			)
		}
	}

	return nil
}

//
// validateCSRSigner validates the signature type.
//
func (v *CSRStampsValidator) validateCSRSigner(signerType string) (err error) {

	if model.VirgilSignatureType == signerType {
		return api.ErrCSRStampSignerIsIncorrect
	}
	if CSRStampSignerMaxLength < len(signerType) {
		return api.ErrCSRStampSignerIsTooLong
	}
	if 0 == len(signerType) {
		return api.ErrCSRStampSignerIsEmpty
	}

	return nil
}

//
// validateCSRStampSnapshot validates the CSR stamp snapshot parameter.
//
func (v *CSRStampsValidator) validateCSRStampSnapshot(snapshot string) (err error) {

	if "" == snapshot {
		return nil
	}

	decodedSnapshot, err := v.encoder.DecodeString(snapshot)
	if nil != err {
		return errors.Wrap(err, api.ErrCSRStampExtraSnapshotDecoding.WithMessage(
			`decode snapshot (%s) error`, snapshot,
		))
	}

	if CSRStampSnapshotMaxLength < len(decodedSnapshot) {
		return api.ErrExtraContentSnapshotIsTooLong
	}

	return nil
}
