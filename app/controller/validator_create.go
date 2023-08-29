package controller

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"

	"github.com/VirgilSecurity/virgil-services-cards/src/api"
	"github.com/VirgilSecurity/virgil-services-cards/src/dao"
	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto"
	"github.com/VirgilSecurity/virgil-services-cards/src/model"
)

//
// CreateCardValidatorProvider is the create card validation interface.
//
type CreateCardValidatorProvider interface {
	//
	// validate validates Virgil Card create request.
	//
	Validate(
		span tracer.Span,
		request *api.CardCreateRequest,
		virgilCard *model.CardDTO,
	) (err error)
}

//
// CreateCardValidator performs card create request validations.
//
type CreateCardValidator struct {
	BaseCardValidator
}

//
// NewCreateCardValidator creates new instance of validator for create card request.
//
func NewCreateCardValidator(
	cardRepository dao.CardRepositoryProvider,
	crypto crypto.Provider,
	csrValidator CSRValidatorProvider,
	csrStampsValidator CSRStampsValidatorProvider,
) *CreateCardValidator {

	return &CreateCardValidator{
		BaseCardValidator: BaseCardValidator{
			crypto:             crypto,
			csrValidator:       csrValidator,
			cardRepository:     cardRepository,
			csrStampsValidator: csrStampsValidator,
		},
	}
}

//
// Validate performs a validation of the search Virgil Card request object.
//
func (v *CreateCardValidator) Validate(
	span tracer.Span,
	request *api.CardCreateRequest,
	virgilCard *model.CardDTO,
) error {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentValidator,
		},
	)
	defer span.Finish()

	if err := v.BaseCardValidator.Validate(span, request, virgilCard, true); nil != err {
		return err
	}

	if err := v.validateCSRPublicKeyLength(virgilCard.PublicKey); nil != err {
		return err
	}

	return nil
}

//
// validateCSRPublicKeyLength validates public key value length.
//
func (v *CreateCardValidator) validateCSRPublicKeyLength(key []byte) error {

	if PublicKeyMinLength > len(key) {
		return api.ErrCSRPublicKeyIsTooShort
	}

	return nil
}
