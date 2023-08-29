package controller

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/VirgilSecurity/virgil-services-core-kit/models"
	"github.com/VirgilSecurity/virgil-services-core-kit/test/helper"

	"github.com/VirgilSecurity/virgil-services-cards/src/api"
	"github.com/VirgilSecurity/virgil-services-cards/src/model"
	"github.com/VirgilSecurity/virgil-services-cards/test/mock"
)

//
// validate :: for an empty request body :: returns an error.
//
func TestDeleteValidateForAnEmptyRequest(t *testing.T) {

	validator := getDeleteCardValidatorUnderTest(validatorDeps{})

	r := api.CardBaseRequest{
		Headers: &api.Headers{
			ApplicationID: emptyScopeID,
			UserID:        emptyRequestIdentity,
		},
	}
	err := validator.Validate(mock.StartNoopSpan(), &r, new(model.CardDTO))

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRIsEmpty, err)
}

//
// validate :: for an invalid previous card ID :: returns an error.
//
func TestDeleteValidateForAnInvalidPreviousCardID(t *testing.T) {

	scopeID := strings.Repeat("a", models.IDLength)
	scr, err := getCSRMessageAsJSON(
		publicKeyEncoded,
		validIdentity,
		validCardVersion,
		validID,
		time.Now().UTC().Unix(),
	)

	assert.Nil(t, err)

	encoder, encodedCSR := presetEncoder(scr)
	encoder.On("DecodeString", publicKeyEncoded).Return(publicKeyBytes, nil)

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On(
		"DoesCardExistByPreviousIDAndScopeID",
		validID,
		scopeID,
	).Return(true)

	crypto := new(mock.Crypto)
	validator := getDeleteCardValidatorUnderTest(validatorDeps{
		encoder:        encoder,
		crypto:         crypto,
		cardRepository: cardRepositoryMock,
	})

	err = validator.Validate(mock.StartNoopSpan(), &api.CardDeleteRequest{
		Headers: &api.Headers{
			UserID:        validIdentity,
			ApplicationID: scopeID,
		},
		CSR:       encodedCSR,
		CSRStamps: []api.CSRStamp{},
	}, new(model.CardDTO))

	assert.Error(t, err)
	assert.Equal(t, api.ErrPreviousVirgilCardExistsAlready, err)
}

//
// TestDeleteValidateForEmptyPreviousCardID :: for an empty previous card ID.
//
func TestDeleteValidateForEmptyPreviousCardID(t *testing.T) {

	scr, err := getCSRMessageAsJSON(
		emptyString,
		validIdentity,
		validCardVersion,
		emptyString,
		time.Now().UTC().Unix(),
	)

	assert.Nil(t, err)

	cardID := strings.Repeat("a", models.IDLength)

	encoder, encodedCSR := presetEncoder(scr)

	crypto := new(mock.Crypto)
	crypto.On("CalculateCardID", []byte(scr)).Return(cardID)

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On("GetCardByID", cardID).Return(new(model.CardDTO), errCardGet)

	validator := getDeleteCardValidatorUnderTest(validatorDeps{
		encoder:        encoder,
		crypto:         crypto,
		cardRepository: cardRepositoryMock,
	})

	err = validator.Validate(mock.StartNoopSpan(), &api.CardDeleteRequest{
		Headers: &api.Headers{
			UserID:        validIdentity,
			ApplicationID: emptyScopeID,
		},
		CSR:       encodedCSR,
		CSRStamps: []api.CSRStamp{},
	}, new(model.CardDTO))

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRPrevCardIDMustNotBeEmpty, err)
}

//
// validate :: for a card that has duplicates :: returns an error.
//
func TestDeleteValidateForACardThatHasDuplicates(t *testing.T) {

	scr, err := getCSRMessageAsJSON(
		publicKeyEncoded,
		validIdentity,
		validCardVersion,
		emptyString,
		time.Now().UTC().Unix(),
	)

	assert.Nil(t, err)

	signatureBytes := []byte("signature")
	signatureEncoded := encodeMessage(signatureBytes)
	cardID := strings.Repeat("a", models.IDLength)

	encoder, encodedCSR := presetEncoder(scr)
	encoder.On("DecodeString", publicKeyEncoded).Return(publicKeyBytes, nil)
	encoder.On("DecodeString", signatureEncoded).Return(signatureBytes, nil)

	crypto := new(mock.Crypto)
	crypto.On("ValidateVirgilCardSignature",
		[]byte(scr), []byte{}, publicKeyBytes, signatureBytes,
	).Return(nil)
	crypto.On("CalculateCardID", []byte(scr)).Return(cardID)

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On("GetCardByID", cardID).Return(new(model.CardDTO), nil)
	cardRepositoryMock.On("DoesCardExistByPreviousIDAndScopeID",
		cardID,
		emptyString,
	).Return(true)

	validator := getDeleteCardValidatorUnderTest(validatorDeps{
		encoder:        encoder,
		crypto:         crypto,
		cardRepository: cardRepositoryMock,
	})

	err = validator.Validate(mock.StartNoopSpan(), &api.CardDeleteRequest{
		Headers: &api.Headers{
			UserID:        validIdentity,
			ApplicationID: emptyScopeID,
		},
		CSR:       encodedCSR,
		CSRStamps: []api.CSRStamp{},
	}, new(model.CardDTO))

	assert.Error(t, err)
	assert.Equal(t, api.ErrVirgilCardContentSnapshotIsNotUnique, err)
}

//
// validatePreviousCardID :: for an empty ID :: passes.
//
func TestDeleteValidatePreviousCardIDForAnEmptyID(t *testing.T) {

	emptyID := ""
	validator := getDeleteCardValidatorUnderTest(validatorDeps{})

	err := validator.validatePreviousCardID(mock.StartNoopSpan(), emptyID, emptyString, emptyString)

	assert.Empty(t, err)
}

//
// validatePreviousCardID :: for a previous card ID that has already been previous :: returns an error.
//
func TestDeleteValidatePreviousCardIDForAnExistingPreviousCard(t *testing.T) {

	scopeID := validID
	alreadyClaimedPreviousCardID := validID

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On("DoesCardExistByPreviousIDAndScopeID",
		mock.Anything,
		mock.Anything,
	).Return(true)

	validator := getDeleteCardValidatorUnderTest(validatorDeps{cardRepository: cardRepositoryMock})

	err := validator.validatePreviousCardID(mock.StartNoopSpan(), alreadyClaimedPreviousCardID, scopeID, "")

	assert.Error(t, err)
	assert.Equal(t, api.ErrPreviousVirgilCardExistsAlready, err)
}

//
// validatePreviousCardID :: for an empty ID :: returns an error.
//
func TestDeleteValidatePreviousCardIDForANotExistingCardID(t *testing.T) {

	nonExistingCardID := "Some non-existing card ID"
	scopeID := validID

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On("GetCardByID", nonExistingCardID).Return(new(model.CardDTO), errCardGet)
	cardRepositoryMock.On(
		"DoesCardExistByPreviousIDAndScopeID",
		mock.Anything,
		mock.Anything,
	).Return(false)

	validator := getDeleteCardValidatorUnderTest(validatorDeps{cardRepository: cardRepositoryMock})

	err := validator.validatePreviousCardID(mock.StartNoopSpan(), nonExistingCardID, scopeID, "")

	assert.Error(t, err)
	assert.Equal(t, api.ErrPreviousVirgilCardDoesNotExist, helper.ExtractHTTPError(err))
}

//
// TestDeleteValidatePreviousCardIDForAPreviousCardIDRegisteredForAnotherScope
// validatePreviousCardID :: for a previous card ID registered for another scope :: returns an error.
//
func TestDeleteValidatePreviousCardIDForAPreviousCardIDRegisteredForAnotherScope(t *testing.T) {

	previousCardID := "Some non-existing card ID"
	incorrectScopeID := "Some incorrect application ID"
	correctScopeID := "Some correct application ID"

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On(
		"DoesCardExistByPreviousIDAndScopeID",
		mock.Anything,
		mock.Anything,
	).Return(false)
	cardRepositoryMock.On(
		"GetCardByID",
		previousCardID,
	).Return(&model.CardDTO{
		ApplicationID: incorrectScopeID,
	}, nil)

	validator := getDeleteCardValidatorUnderTest(validatorDeps{cardRepository: cardRepositoryMock})

	err := validator.validatePreviousCardID(mock.StartNoopSpan(), previousCardID, correctScopeID, "")

	assert.Error(t, err)
	assert.Equal(t, err, api.ErrPreviousVirgilCardIsRegisteredForAnotherScope)
}

//
// TestDeleteValidatePreviousCardIDForAPreviousCardIdentityNotEqualsCurrentOne
// validatePreviousCardID :: for a previous card ID which identity differs from current one :: returns an error.
//
func TestDeleteValidatePreviousCardIDForAPreviousCardIdentityNotEqualsCurrentOne(t *testing.T) {

	previousCardID := "Some non-existing card ID"
	scopeID := "Some scope ID"
	cardRepositoryMock := mock.CardRepository{}
	cardRepositoryMock.On(
		"DoesCardExistByPreviousIDAndScopeID",
		mock.Anything,
		mock.Anything,
	).Return(false)
	cardRepositoryMock.On("GetCardByID", previousCardID).Return(&model.CardDTO{
		Identity:      "another identity",
		ApplicationID: scopeID,
	}, nil)

	validator := getDeleteCardValidatorUnderTest(validatorDeps{cardRepository: &cardRepositoryMock})

	err := validator.validatePreviousCardID(mock.StartNoopSpan(), previousCardID, scopeID, "identity")

	assert.Error(t, err)
	assert.Equal(t, api.ErrPreviousVirgilCardIdentityIsIncorrect, err)
}

//
// TestDeleteValidatePreviousCardIDWithAValidPreviousCardIDPasses: validatePreviousCardID :: for a valid previous card ID :: passes.
//
func TestDeleteValidatePreviousCardIDWithAValidPreviousCardIDPasses(t *testing.T) {

	previousCardID := "Some non-existing card ID"
	correctScopeID := "Some correct application ID"
	cardRepositoryMock := mock.CardRepository{}
	cardRepositoryMock.On(
		"DoesCardExistByPreviousIDAndScopeID",
		mock.Anything,
		mock.Anything,
	).Return(false)
	cardRepositoryMock.On("GetCardByID", previousCardID).Return(&model.CardDTO{
		ID:            previousCardID,
		Identity:      validIdentity,
		ApplicationID: correctScopeID,
	}, nil)

	validator := getDeleteCardValidatorUnderTest(validatorDeps{cardRepository: &cardRepositoryMock})

	err := validator.validatePreviousCardID(mock.StartNoopSpan(), previousCardID, correctScopeID, validIdentity)

	assert.Empty(t, err)
}

//
// TestDeleteValidateVirgilCardDuplicatesForAnExistingCardID :: for an existing card ID :: returns an error.
//
func TestDeleteValidateVirgilCardDuplicatesForAnExistingCardID(t *testing.T) {

	existingCardID := "Existing card ID"

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On("GetCardByID", existingCardID).Return(&model.CardDTO{
		ID: existingCardID,
	}, nil)

	validator := getDeleteCardValidatorUnderTest(validatorDeps{cardRepository: cardRepositoryMock})

	err := validator.validateVirgilCardDuplicates(mock.StartNoopSpan(), existingCardID)

	assert.Error(t, err)
	assert.Equal(t, err, api.ErrVirgilCardContentSnapshotIsNotUnique)
}

//
// TestDeleteValidateVirgilCardDuplicatesForANotExistingCardID :: for a not existing card ID :: returns an error.
//
func TestDeleteValidateVirgilCardDuplicatesForANotExistingCardID(t *testing.T) {

	notExistingCardID := "Not existing card ID"

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On(
		"GetCardByID",
		notExistingCardID,
		mock.Anything,
	).Return(new(model.CardDTO), errCardGet)

	validator := getDeleteCardValidatorUnderTest(validatorDeps{cardRepository: cardRepositoryMock})

	err := validator.validateVirgilCardDuplicates(mock.StartNoopSpan(), notExistingCardID)

	assert.Empty(t, err)
}

//
// TestDeleteValidateWithNotEmptyPublicKey :: for an empty public key :: returns an error.
//
func TestDeleteValidateWithNotEmptyPublicKey(t *testing.T) {

	scr, err := getCSRMessageAsJSON(
		encodedPublicKey,
		validIdentity,
		validCardVersion,
		validID,
		time.Now().UTC().Unix(),
	)

	assert.Nil(t, err)

	signatureBytes := []byte("signature")
	signatureEncoded := encodeMessage(signatureBytes)
	cardID := strings.Repeat("a", models.IDLength)

	encoder, encodedCSR := presetEncoder(scr)
	encoder.On("DecodeString", encodedPublicKey).Return([]byte(validPublicKey), nil)
	encoder.On("DecodeString", signatureEncoded).Return(signatureBytes, nil)

	crypto := new(mock.Crypto)
	crypto.On("ValidateVirgilCardSignature",
		[]byte(scr), []byte{}, []byte(validPublicKey), signatureBytes,
	).Return(nil)
	crypto.On("CalculateCardID", []byte(scr)).Return(cardID)

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On("GetCardByID", cardID).Return(new(model.CardDTO), errCardGet)
	cardRepositoryMock.On("GetCardByID", validID).Return(&model.CardDTO{Identity: validIdentity}, nil)
	cardRepositoryMock.On("DoesCardExistByPreviousIDAndScopeID",
		validID,
		emptyScopeID,
	).Return(false)

	validator := getDeleteCardValidatorUnderTest(validatorDeps{
		encoder:        encoder,
		crypto:         crypto,
		cardRepository: cardRepositoryMock,
	})

	err = validator.Validate(mock.StartNoopSpan(), &api.CardDeleteRequest{
		Headers: &api.Headers{
			ApplicationID: emptyScopeID,
			UserID:        validIdentity,
		},
		CSR:       encodedCSR,
		CSRStamps: []api.CSRStamp{},
	}, new(model.CardDTO))

	assert.Equal(t, err, api.ErrCSRPublicKeyMustBeEmpty)
}

//
// getDeleteCardValidatorUnderTest returns validator instance.
//
func getDeleteCardValidatorUnderTest(deps validatorDeps) *DeleteCardValidator {

	if nil == deps.cardRepository {
		deps.cardRepository = &mock.CardRepository{}
	}

	if nil == deps.encoder {
		deps.encoder = &mock.Base64Encoder{}
	}

	if nil == deps.crypto {
		deps.crypto = &mock.Crypto{}
	}

	return NewDeleteCardValidator(deps.cardRepository, deps.crypto, &CSRValidator{
		encoder: deps.encoder,
	}, &CSRStampsValidator{
		crypto:  deps.crypto,
		encoder: deps.encoder,
	})
}
