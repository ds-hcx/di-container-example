package controller

import (
	"encoding/json"

	"github.com/VirgilSecurity/virgil-services-core-kit/errors"
	"github.com/VirgilSecurity/virgil-services-core-kit/models"
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"

	"github.com/VirgilSecurity/virgil-services-cards/src/api"
	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/encoder"
	"github.com/VirgilSecurity/virgil-services-cards/src/model"
)

//
// Validation constraints.
//
const (
	PublicKeyMinLength = 16
	PublicKeyMaxLength = 4096

	IdentityMinLength = 1
	IdentityMaxLength = 1024
)

//
// CSRValidatorProvider provides an interface to work with the CSR validator.
//
type CSRValidatorProvider interface {
	//
	// Validate validates Virgil Card CSR.
	//
	Validate(span tracer.Span, csrData string, csr *api.CSR, identity string) (params *ParametersStore, err error)
}

//
// CSRValidator represents the CSR validator.
//
type CSRValidator struct {
	encoder encoder.Provider
}

//
// NewCSRValidator creates new instance of CSR validator.
//
func NewCSRValidator(encoder encoder.Provider) *CSRValidator {
	return &CSRValidator{
		encoder: encoder,
	}
}

//
// Validate extends Virgil Card CSR to ParametersStore and validates its fields.
//
func (v *CSRValidator) Validate(
	span tracer.Span,
	csrData string,
	csr *api.CSR,
	identity string,
) (params *ParametersStore, err error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentValidator,
		},
	)
	defer span.Finish()

	if "" == csrData {
		return nil, tracer.SetSpanErrorAndReturn(span, api.ErrCSRIsEmpty)
	}

	params = new(ParametersStore)
	if err = v.parseCSR(csrData, params, csr); nil != err {
		return nil, tracer.SetSpanErrorAndReturn(span, err)
	}

	params.publicKey, err = v.decodeCSRPublicKeyAndCheckMaxLength(csr.GetPublicKey())
	if nil != err {
		return nil, tracer.SetSpanErrorAndReturn(span, err)
	}
	if err := v.validateCSRIdentity(csr.GetIdentity(), identity); nil != err {
		return nil, tracer.SetSpanErrorAndReturn(span, err)
	}
	if err := v.validateCSRPreviousCardID(csr.GetPreviousCardID()); nil != err {
		return nil, tracer.SetSpanErrorAndReturn(span, err)
	}
	if err := v.validateCSRCreatedAt(csr.GetCreatedAt()); nil != err {
		return nil, tracer.SetSpanErrorAndReturn(span, err)
	}
	if err := v.validateCSRVersion(csr.GetVersion()); nil != err {
		return nil, tracer.SetSpanErrorAndReturn(span, err)
	}

	return params, nil
}

//
// parseCSR parses CSR message string.
//
func (v *CSRValidator) parseCSR(csrMsg string, params *ParametersStore, csr *api.CSR) (err error) {

	if params.csr, err = v.encoder.DecodeString(csrMsg); nil != err {
		return errors.Wrap(err, api.ErrContentSnapshotIsNotABase64EncodedString.WithMessage(
			"content_snapshot (%s) decode error", csrMsg,
		))
	}

	if err = json.Unmarshal(params.csr, csr); nil != err {
		// Corner case to catch up int64 timestamp parsing error.
		if e, ok := err.(*json.UnmarshalTypeError); ok && e.Field == api.CreatedAtFieldName {
			return errors.Wrap(err, api.ErrCSRCreationTimeIsIncorrect.WithMessage(
				"decoded content_snapshot (%s) unmarshal error", csrMsg,
			))
		}

		return errors.Wrap(err, api.ErrContentSnapshotIsNotAJSONMessage.WithMessage(
			"unmarshal decoded content_snapshot (%s)", csrMsg,
		))
	}

	return nil
}

//
// decodeCSRPublicKeyAndCheckMaxLength validates public key value.
//
func (v *CSRValidator) decodeCSRPublicKeyAndCheckMaxLength(key string) (decodedKey []byte, err error) {

	if "" == key {
		return []byte{}, nil
	}

	decodedKey, err = v.encoder.DecodeString(key)
	if nil != err {
		return []byte{}, api.ErrCSRPublicKeyDecoding.WithMessage(
			"public key value (%s) decode error:%+v", key, err,
		)
	}

	if PublicKeyMaxLength < len(decodedKey) {
		return []byte{}, api.ErrCSRPublicKeyIsTooLong
	}

	return decodedKey, nil
}

//
// validateCSRIdentity validates identity value.
//
func (v *CSRValidator) validateCSRIdentity(identity, requestIdentity string) (err error) {

	if "" == identity {
		return api.ErrCSRIdentityIsEmpty
	}
	if IdentityMinLength > len(identity) || IdentityMaxLength < len(identity) {
		return api.ErrCSRIdentityIsIncorrect
	}
	if identity != requestIdentity {
		return api.ErrCSRIdentityDoesNotMatchRequestIdentity
	}

	return nil
}

//
// validateCSRPreviousCardID validates previous card id value.
//
func (v *CSRValidator) validateCSRPreviousCardID(id string) (err error) {

	if "" == id {
		return
	}
	if models.IDLength != len(id) {
		return api.ErrCSRPreviousCardIDIsIncorrect
	}

	return nil
}

//
// validateCSRCreatedAt validates creation timestamp value.
//
func (v *CSRValidator) validateCSRCreatedAt(createdAt int64) (err error) {

	if 0 >= createdAt {
		return api.ErrCSRCreationTimeIsIncorrect
	}

	return nil
}

//
// validateCSRVersion validates version value.
//
func (v *CSRValidator) validateCSRVersion(ver string) (err error) {

	if model.CardVersion5 != ver {
		return api.ErrCSRVersionIsIncorrect
	}

	return nil
}
