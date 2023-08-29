package controller

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"

	"github.com/VirgilSecurity/virgil-services-cards/src/api"
	"github.com/VirgilSecurity/virgil-services-cards/src/dao"
	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto"
	"github.com/VirgilSecurity/virgil-services-cards/src/model"
)

//
// ParametersStore is an object that contains validated and decoded parameters.
// It is used to prevent multiple decoding on the same structure.
//
type ParametersStore struct {
	publicKey     []byte
	csr           []byte
	extraSnapshot []byte
}

//
// BaseCardValidator performs base card request validations.
//
type BaseCardValidator struct {
	crypto             crypto.Provider
	csrValidator       CSRValidatorProvider
	cardRepository     dao.CardRepositoryProvider
	csrStampsValidator CSRStampsValidatorProvider
}

//
// Validate performs a validation of the Virgil Card base request object.
//
func (v *BaseCardValidator) Validate(
	span tracer.Span,
	cardBaseRequest *api.CardBaseRequest,
	virgilCard *model.CardDTO,
	verifySignatures bool,
) error {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentValidator,
		},
	)
	defer span.Finish()

	var (
		csr = new(api.CSR)
	)

	decodedParams, err := v.csrValidator.Validate(span, cardBaseRequest.CSR, csr, cardBaseRequest.UserID)
	if nil != err {
		return err
	}

	if verifySignatures {
		if err := v.csrStampsValidator.Validate(
			span,
			cardBaseRequest.CSRStamps,
			cardBaseRequest.ApplicationID,
			decodedParams,
		); nil != err {
			return err
		}
	}

	if err := v.validatePreviousCardID(
		span,
		csr.GetPreviousCardID(),
		cardBaseRequest.ApplicationID,
		csr.GetIdentity(),
	); nil != err {
		return err
	}

	if err := v.fillVirgilCardModel(
		span,
		virgilCard,
		csr,
		cardBaseRequest,
		cardBaseRequest.ApplicationID,
		decodedParams,
	); nil != err {
		return err
	}

	if err := v.validateVirgilCardDuplicates(span, virgilCard.GetID()); nil != err {
		return err
	}

	return nil
}

//
// fillVirgilCardModel fills Card object with already validated data from the request.
//
func (v *BaseCardValidator) fillVirgilCardModel(
	span tracer.Span,
	card *model.CardDTO,
	csr *api.CSR,
	request *api.CardBaseRequest,
	scopeID string,
	decodedParams *ParametersStore,
) (err error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentValidator,
		},
	)
	defer span.Finish()

	card.ID = v.crypto.CalculateCardID(decodedParams.csr)
	card.ContentSnapshot = request.GetCSR()
	card.PublicKey = decodedParams.publicKey
	card.Identity = csr.GetIdentity()
	card.Version = csr.GetVersion()
	card.CreatedAt = csr.GetCreatedAt()
	card.PreviousCardID = csr.GetPreviousCardID()
	card.ApplicationID = scopeID

	// Convert the CSR stamps.
	for _, s := range request.CSRStamps {
		card.AppendSignature(&model.CardSignatureDTO{
			Signer:    s.GetSigner(),
			Snapshot:  s.GetSnapshot(),
			Signature: s.GetSignature(),
		})
	}

	return
}

//
// validatePreviousCardID performs the validation of the previous card ID parameter.
//
func (v *BaseCardValidator) validatePreviousCardID(span tracer.Span, previousCardID, scopeID, identity string) error {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentValidator,
		},
	)
	defer span.Finish()

	if "" == previousCardID {
		return nil
	}

	// search by previous card ID
	isSuperseeded, err := v.cardRepository.DoesCardExistByPreviousIDAndScopeID(span, previousCardID, scopeID)
	if err != nil {
		return api.ErrInternalError.WithMessage("error checking if card already exists: %+v", err)
	}

	if isSuperseeded {
		return tracer.SetSpanErrorAndReturn(span, api.ErrPreviousVirgilCardExistsAlready)
	}

	card, err := v.cardRepository.GetCardByID(span, previousCardID)
	if nil != err {
		return api.ErrPreviousVirgilCardDoesNotExist.WithMessage("error: %+v", err)
	}
	if scopeID != card.GetApplicationID() {
		return tracer.SetSpanErrorAndReturn(span, api.ErrPreviousVirgilCardIsRegisteredForAnotherScope)
	}
	if identity != card.Identity {
		return tracer.SetSpanErrorAndReturn(span, api.ErrPreviousVirgilCardIdentityIsIncorrect)
	}

	return nil
}

//
// validateVirgilCardDuplicates performs the validation that identical Virgil Cards do not exist already.
//
func (v *BaseCardValidator) validateVirgilCardDuplicates(span tracer.Span, virgilCardID string) error {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentValidator,
		},
	)
	defer span.Finish()

	_, err := v.cardRepository.GetCardByID(span, virgilCardID)
	if nil == err {
		return tracer.SetSpanErrorAndReturn(span, api.ErrVirgilCardContentSnapshotIsNotUnique)
	}

	return nil
}
