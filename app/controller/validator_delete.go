package controller

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"

	"github.com/VirgilSecurity/virgil-services-cards/src/api"
	"github.com/VirgilSecurity/virgil-services-cards/src/dao"
	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto"
	"github.com/VirgilSecurity/virgil-services-cards/src/model"
)

//
// DeleteCardValidatorProvider is the delete card validation interface.
//
type DeleteCardValidatorProvider interface {
	//
	// validate validates Virgil Card delete request.
	//
	Validate(
		span tracer.Span,
		request *api.CardDeleteRequest,
		virgilCard *model.CardDTO,
	) (err error)
}

//
// DeleteCardValidator performs card create request validations.
//
type DeleteCardValidator struct {
	BaseCardValidator
}

//
// NewDeleteCardValidator creates new instance of validator for create card request.
//
func NewDeleteCardValidator(
	cardRepository dao.CardRepositoryProvider,
	crypto crypto.Provider,
	csrValidator CSRValidatorProvider,
	csrStampsValidator CSRStampsValidatorProvider,
) *DeleteCardValidator {

	return &DeleteCardValidator{
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
func (v *DeleteCardValidator) Validate(
	span tracer.Span,
	request *api.CardDeleteRequest,
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

	if err := v.BaseCardValidator.Validate(span, request, virgilCard, false); nil != err {
		return err
	}

	if err := v.validateCSRPublicKeyLength(virgilCard.PublicKey); nil != err {
		return err
	}

	if err := v.validateCSRCardIDToDelete(virgilCard.PreviousCardID); nil != err {
		return err
	}

	return nil
}

//
// validateCSRPublicKeyLength validates public key value length.
//
func (v *DeleteCardValidator) validateCSRPublicKeyLength(key []byte) error {

	if len(key) > 0 {
		return api.ErrCSRPublicKeyMustBeEmpty
	}
	return nil
}

//
// validateCSRCardIDToDelete validates card ID to delete.
//
func (v *DeleteCardValidator) validateCSRCardIDToDelete(previousCardID string) error {

	if 0 >= len(previousCardID) {
		return api.ErrCSRPrevCardIDMustNotBeEmpty
	}
	return nil
}
