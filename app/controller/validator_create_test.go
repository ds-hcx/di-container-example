package controller

import (
	"errors"
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
// Testing constants.
//
const (
	encodedString  = "b3JpZ2luYWwgc3RyaW5n"
	originalString = "original string"
	tooShortPK     = "PKLessThan16b--"
)

//
// Testing variables.
//
var (
	errCardGet        = errors.New("unable to get the card")
	tooShortPKEncoded = encodeMessage([]byte(tooShortPK))
)

//
// validate :: for an empty request body :: returns an error.
//
func TestValidateForAnEmptyRequest(t *testing.T) {

	validator := getCreateCardValidatorUnderTest(validatorDeps{})

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
// validate :: for an empty CSR stamps list :: returns an error.
//
func TestValidateForAnEmptyCSRStampList(t *testing.T) {

	scr, err := getCSRMessageAsJSON(
		publicKeyEncoded,
		validIdentity,
		validCardVersion,
		"",
		time.Now().UTC().Unix(),
	)

	assert.Nil(t, err)

	encoder, encodedCSR := presetEncoder(scr)
	encoder.On("DecodeString", publicKeyEncoded).Return(publicKeyBytes, nil)

	validator := getCreateCardValidatorUnderTest(validatorDeps{encoder: encoder})

	err = validator.Validate(mock.StartNoopSpan(), &api.CardCreateRequest{
		Headers: &api.Headers{
			UserID:        validIdentity,
			ApplicationID: emptyScopeID,
		},
		CSR: encodedCSR,
	}, new(model.CardDTO))

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRStampsListIsTooSmall, err)
}

//
// validate :: for an invalid previous card ID :: returns an error.
//
func TestValidateForAnInvalidPreviousCardID(t *testing.T) {

	previousCardID := validID
	scopeID := strings.Repeat("a", models.IDLength)
	scr, err := getCSRMessageAsJSON(
		publicKeyEncoded,
		validIdentity,
		validCardVersion,
		previousCardID,
		time.Now().UTC().Unix(),
	)

	assert.Nil(t, err)

	signatureBytes := []byte("signature")
	signatureEncoded := encodeMessage(signatureBytes)

	encoder, encodedCSR := presetEncoder(scr)
	encoder.On("DecodeString", publicKeyEncoded).Return(publicKeyBytes, nil)
	encoder.On("DecodeString", signatureEncoded).Return(signatureBytes, nil)

	crypto := new(mock.Crypto)
	crypto.On("ValidateVirgilCardSignature",
		[]byte(scr), []byte{}, publicKeyBytes, signatureBytes,
	).Return(nil)

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On(
		"DoesCardExistByPreviousIDAndScopeID",
		previousCardID,
		scopeID,
	).Return(true)

	validator := getCreateCardValidatorUnderTest(validatorDeps{
		encoder:        encoder,
		crypto:         crypto,
		cardRepository: cardRepositoryMock,
	})

	err = validator.Validate(mock.StartNoopSpan(), &api.CardCreateRequest{
		Headers: &api.Headers{
			UserID:        validIdentity,
			ApplicationID: scopeID,
		},
		CSR: encodedCSR,
		CSRStamps: []api.CSRStamp{{
			Signer:    model.SelfSignatureType,
			Signature: signatureEncoded,
		}},
	}, new(model.CardDTO))

	assert.Error(t, err)
	assert.Equal(t, api.ErrPreviousVirgilCardExistsAlready, err)
}

//
// validate :: for all valid parameters :: passes.
//
func TestValidateForAllValidParameters(t *testing.T) {

	scr, err := getCSRMessageAsJSON(
		publicKeyEncoded,
		validIdentity,
		validCardVersion,
		"",
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
	cardRepositoryMock.On("GetCardByID", cardID).Return(new(model.CardDTO), errCardGet)

	validator := getCreateCardValidatorUnderTest(validatorDeps{
		encoder:        encoder,
		crypto:         crypto,
		cardRepository: cardRepositoryMock,
	})

	err = validator.Validate(mock.StartNoopSpan(), &api.CardCreateRequest{
		Headers: &api.Headers{
			UserID:        validIdentity,
			ApplicationID: emptyScopeID,
		},
		CSR: encodedCSR,
		CSRStamps: []api.CSRStamp{{
			Signer:    model.SelfSignatureType,
			Signature: signatureEncoded,
		}},
	}, new(model.CardDTO))

	assert.Empty(t, err)
}

//
// validate :: for a card that has duplicates :: returns an error.
//
func TestValidateForACardThatHasDuplicates(t *testing.T) {

	scr, err := getCSRMessageAsJSON(
		publicKeyEncoded,
		validIdentity,
		validCardVersion,
		"",
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

	validator := getCreateCardValidatorUnderTest(validatorDeps{
		encoder:        encoder,
		crypto:         crypto,
		cardRepository: cardRepositoryMock,
	})

	err = validator.Validate(mock.StartNoopSpan(), &api.CardCreateRequest{
		Headers: &api.Headers{
			UserID:        validIdentity,
			ApplicationID: emptyScopeID,
		},
		CSR: encodedCSR,
		CSRStamps: []api.CSRStamp{{
			Signer:    model.SelfSignatureType,
			Signature: signatureEncoded,
		}},
	}, new(model.CardDTO))

	assert.Error(t, err)
	assert.Equal(t, api.ErrVirgilCardContentSnapshotIsNotUnique, err)
}

//
// validatePreviousCardID :: for an empty ID :: passes.
//
func TestValidatePreviousCardIDForAnEmptyID(t *testing.T) {

	emptyID := ""
	validator := getCreateCardValidatorUnderTest(validatorDeps{})

	err := validator.validatePreviousCardID(mock.StartNoopSpan(), emptyID, emptyString, emptyString)

	assert.Empty(t, err)
}

//
// validatePreviousCardID :: for a previous card ID that has already been previous :: returns an error.
//
func TestValidatePreviousCardIDForAnExistingPreviousCard(t *testing.T) {

	scopeID := validID
	alreadyClaimedPreviousCardID := validID

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On("DoesCardExistByPreviousIDAndScopeID",
		mock.Anything,
		mock.Anything,
	).Return(true)

	validator := getCreateCardValidatorUnderTest(validatorDeps{cardRepository: cardRepositoryMock})

	err := validator.validatePreviousCardID(mock.StartNoopSpan(), alreadyClaimedPreviousCardID, scopeID, "")

	assert.Error(t, err)
	assert.Equal(t, api.ErrPreviousVirgilCardExistsAlready, err)
}

//
// validatePreviousCardID :: for an empty ID :: returns an error.
//
func TestValidatePreviousCardIDForANotExistingCardID(t *testing.T) {

	nonExistingCardID := "Some non-existing card ID"
	scopeID := validID

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On("GetCardByID", nonExistingCardID).Return(new(model.CardDTO), errCardGet)
	cardRepositoryMock.On(
		"DoesCardExistByPreviousIDAndScopeID",
		mock.Anything,
		mock.Anything,
	).Return(false)

	validator := getCreateCardValidatorUnderTest(validatorDeps{cardRepository: cardRepositoryMock})

	err := validator.validatePreviousCardID(mock.StartNoopSpan(), nonExistingCardID, scopeID, "")

	assert.Error(t, err)
	assert.Equal(t, api.ErrPreviousVirgilCardDoesNotExist, helper.ExtractHTTPError(err))
}

//
// TestValidatePreviousCardIDForAPreviousCardIDRegisteredForAnotherScope
// validatePreviousCardID :: for a previous card ID registered for another scope :: returns an error.
//
func TestValidatePreviousCardIDForAPreviousCardIDRegisteredForAnotherScope(t *testing.T) {

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

	validator := getCreateCardValidatorUnderTest(validatorDeps{cardRepository: cardRepositoryMock})

	err := validator.validatePreviousCardID(mock.StartNoopSpan(), previousCardID, correctScopeID, "")

	assert.Error(t, err)
	assert.Equal(t, err, api.ErrPreviousVirgilCardIsRegisteredForAnotherScope)
}

//
// TestValidatePreviousCardIDForAPreviousCardIdentityNotEqualsCurrentOne
// validatePreviousCardID :: for a previous card ID which identity differs from current one :: returns an error.
//
func TestValidatePreviousCardIDForAPreviousCardIdentityNotEqualsCurrentOne(t *testing.T) {

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

	validator := getCreateCardValidatorUnderTest(validatorDeps{cardRepository: &cardRepositoryMock})

	err := validator.validatePreviousCardID(mock.StartNoopSpan(), previousCardID, scopeID, "identity")

	assert.Error(t, err)
	assert.Equal(t, api.ErrPreviousVirgilCardIdentityIsIncorrect, err)
}

//
// TestValidatePreviousCardIDWithAValidPreviousCardIDPasses: validatePreviousCardID :: for a valid previous card ID :: passes.
//
func TestValidatePreviousCardIDWithAValidPreviousCardIDPasses(t *testing.T) {

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

	validator := getCreateCardValidatorUnderTest(validatorDeps{cardRepository: &cardRepositoryMock})

	err := validator.validatePreviousCardID(mock.StartNoopSpan(), previousCardID, correctScopeID, validIdentity)

	assert.Empty(t, err)
}

//
// TestValidateVirgilCardDuplicatesForAnExistingCardID :: for an existing card ID :: returns an error.
//
func TestValidateVirgilCardDuplicatesForAnExistingCardID(t *testing.T) {

	existingCardID := "Existing card ID"

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On("GetCardByID", existingCardID).Return(&model.CardDTO{
		ID: existingCardID,
	}, nil)

	validator := getCreateCardValidatorUnderTest(validatorDeps{cardRepository: cardRepositoryMock})

	err := validator.validateVirgilCardDuplicates(mock.StartNoopSpan(), existingCardID)

	assert.Error(t, err)
	assert.Equal(t, err, api.ErrVirgilCardContentSnapshotIsNotUnique)
}

//
// TestValidateVirgilCardDuplicatesForANotExistingCardID :: for a not existing card ID :: returns an error.
//
func TestValidateVirgilCardDuplicatesForANotExistingCardID(t *testing.T) {

	notExistingCardID := "Not existing card ID"

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On(
		"GetCardByID",
		notExistingCardID,
		mock.Anything,
	).Return(new(model.CardDTO), errCardGet)

	validator := getCreateCardValidatorUnderTest(validatorDeps{cardRepository: cardRepositoryMock})

	err := validator.validateVirgilCardDuplicates(mock.StartNoopSpan(), notExistingCardID)

	assert.Empty(t, err)
}

//
// TestValidateWithTooShortPublicKey :: for an empty public key :: returns an error.
//
func TestValidateWithTooShortPublicKey(t *testing.T) {

	scr, err := getCSRMessageAsJSON(
		tooShortPKEncoded,
		validIdentity,
		validCardVersion,
		"",
		time.Now().UTC().Unix(),
	)

	assert.Nil(t, err)

	signatureBytes := []byte("signature")
	signatureEncoded := encodeMessage(signatureBytes)
	cardID := strings.Repeat("a", models.IDLength)

	encoder, encodedCSR := presetEncoder(scr)
	encoder.On("DecodeString", tooShortPKEncoded).Return([]byte(tooShortPK), nil)
	encoder.On("DecodeString", signatureEncoded).Return(signatureBytes, nil)

	crypto := new(mock.Crypto)
	crypto.On("ValidateVirgilCardSignature",
		[]byte(scr), []byte{}, []byte(tooShortPK), signatureBytes,
	).Return(nil)
	crypto.On("CalculateCardID", []byte(scr)).Return(cardID)

	cardRepositoryMock := new(mock.CardRepository)
	cardRepositoryMock.On("GetCardByID", cardID).Return(new(model.CardDTO), errCardGet)

	validator := getCreateCardValidatorUnderTest(validatorDeps{
		encoder:        encoder,
		crypto:         crypto,
		cardRepository: cardRepositoryMock,
	})

	err = validator.Validate(mock.StartNoopSpan(), &api.CardCreateRequest{
		Headers: &api.Headers{
			UserID:        validIdentity,
			ApplicationID: emptyScopeID,
		},
		CSR: encodedCSR,
		CSRStamps: []api.CSRStamp{{
			Signer:    model.SelfSignatureType,
			Signature: signatureEncoded,
		}},
	}, new(model.CardDTO))

	assert.Equal(t, err, api.ErrCSRPublicKeyIsTooShort)
}

//
// getCreateCardValidatorUnderTest returns validator instance.
//
func getCreateCardValidatorUnderTest(deps validatorDeps) *CreateCardValidator {

	if nil == deps.cardRepository {
		deps.cardRepository = &mock.CardRepository{}
	}

	if nil == deps.encoder {
		deps.encoder = &mock.Base64Encoder{}
	}

	if nil == deps.crypto {
		deps.crypto = &mock.Crypto{}
	}

	return NewCreateCardValidator(deps.cardRepository, deps.crypto, &CSRValidator{
		encoder: deps.encoder,
	}, &CSRStampsValidator{
		crypto:  deps.crypto,
		encoder: deps.encoder,
	})
}
