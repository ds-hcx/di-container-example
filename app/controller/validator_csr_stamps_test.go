package controller

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/virgil.v5/cryptoimpl"

	"github.com/VirgilSecurity/virgil-services-core-kit/errors"
	"github.com/VirgilSecurity/virgil-services-core-kit/models"
	"github.com/VirgilSecurity/virgil-services-core-kit/test/helper"

	"github.com/VirgilSecurity/virgil-services-cards/src/api"
	"github.com/VirgilSecurity/virgil-services-cards/src/model"
	"github.com/VirgilSecurity/virgil-services-cards/test/mock"
)

//
// validateCSRStamps :: for an empty CSR stamp list :: returns an error.
//
func TestValidateCSRStampsForAnEmptyStampList(t *testing.T) {

	var csrStampList []api.CSRStamp
	validator := NewCSRStampsValidator(&mock.Crypto{}, &mock.Base64Encoder{})

	err := validator.Validate(mock.StartNoopSpan(), csrStampList, emptyString, new(ParametersStore))

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRStampsListIsTooSmall, helper.ExtractHTTPError(err))
}

//
// validateCSRStamps :: for a too long CSR stamp list :: returns an error.
//
func TestValidateCSRStampsForATooLongStampList(t *testing.T) {

	var stampList []api.CSRStamp
	for i := 0; i < CSRStampsListMaxLength+1; i++ {
		stampList = append(stampList, api.CSRStamp{})
	}

	validator := NewCSRStampsValidator(&mock.Crypto{}, &mock.Base64Encoder{})

	err := validator.Validate(mock.StartNoopSpan(), stampList, emptyString, new(ParametersStore))

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRStampsListIsTooLarge, err)
}

//
// validateCSRStamps :: for a CSR stamp list without a self stamp :: returns an error.
//
func TestValidateCSRStampsForAStampListWithoutSelfStamp(t *testing.T) {

	encoder, _ := presetEncoder("")
	encoder.On("DecodeString", encodedString).Return([]byte(originalString), nil)
	validator := NewCSRStampsValidator(&mock.Crypto{}, encoder)
	stampList := []api.CSRStamp{{
		Signer:    model.ApplicationSignatureType,
		Signature: encodedString,
	}}

	err := validator.Validate(mock.StartNoopSpan(), stampList, emptyString, new(ParametersStore))

	assert.Error(t, err)
	assert.Equal(t, api.ErrSelfCSRStampIsMissing, err)
}

//
// validateCSRStamps :: for a too large CSR stamp list :: returns an error.
//
func TestValidateCSRStampsForATooLargeStampList(t *testing.T) {

	encoder, _ := presetEncoder("")
	encoder.On("DecodeString", encodedString).Return([]byte(originalString), nil)
	validator := NewCSRStampsValidator(&mock.Crypto{}, encoder)
	stampList := []api.CSRStamp{{
		Signer:    model.ApplicationSignatureType,
		Signature: encodedString,
	}}

	err := validator.Validate(mock.StartNoopSpan(), stampList, emptyString, new(ParametersStore))

	assert.Error(t, err)
	assert.Equal(t, api.ErrSelfCSRStampIsMissing, err)
}

//
// validateCSRStamps :: for an application CSR stamp not in Scope ID :: passes.
//
func TestValidateCSRStampsForAnApplicationCSRStampNotInScopeID(t *testing.T) {

	publicKeyBytes := []byte(originalString)
	encodedPublicKey := encodedString
	emptyIdentity := emptyString
	validVersion := model.CardVersion5
	emptyPreviousCardID := emptyString
	createdAt := time.Now().UTC().Unix()
	csr, err := getCSRMessageAsJSON(encodedPublicKey, emptyIdentity, validVersion, emptyPreviousCardID, createdAt)

	assert.NoError(t, err)

	signatureBytes := []byte("Signature")
	encodedSignature := encodeMessage(signatureBytes)
	csrParameters := ParametersStore{
		publicKey: publicKeyBytes,
		csr:       []byte(csr),
	}

	crypto := new(mock.Crypto)
	crypto.On("ValidateVirgilCardSignature",
		csrParameters.csr, []byte{}, publicKeyBytes, signatureBytes,
	).Return(nil)
	encoder, _ := presetEncoder(string(signatureBytes))
	validator := NewCSRStampsValidator(crypto, encoder)
	signatureList := []api.CSRStamp{
		{
			Signature: encodedSignature,
			Signer:    model.SelfSignatureType,
		},
		{
			Signature: encodedSignature,
			Signer:    model.ApplicationSignatureType,
		},
	}

	err = validator.Validate(mock.StartNoopSpan(), signatureList, emptyString, &csrParameters)

	assert.NoError(t, err)
}

//
// validateCSRStamps :: for an application CSR stamp in Scope ID :: passes.
//
func TestValidateCSRStampsForAnApplicationCSRStampInScopeID(t *testing.T) {

	publicKeyBytes := []byte(originalString)
	encodedPublicKey := encodedString
	emptyIdentity := emptyString
	validVersion := model.CardVersion5
	emptyPreviousCardID := emptyString
	createdAt := time.Now().UTC().Unix()
	appID := strings.Repeat("b", models.IDLength)
	scopeWithAppID := appID
	csr, err := getCSRMessageAsJSON(encodedPublicKey, emptyIdentity, validVersion, emptyPreviousCardID, createdAt)

	assert.NoError(t, err)

	signatureBytes := []byte("Signature")
	encodedSignature := encodeMessage(signatureBytes)
	csrParameters := ParametersStore{
		publicKey: publicKeyBytes,
		csr:       []byte(csr),
	}

	crypto := new(mock.Crypto)
	crypto.On("ValidateVirgilCardSignature",
		csrParameters.csr, []byte{}, publicKeyBytes, signatureBytes,
	).Return(nil)
	encoder, _ := presetEncoder(string(signatureBytes))
	validator := NewCSRStampsValidator(crypto, encoder)
	signatureList := []api.CSRStamp{
		{
			Signature: encodedSignature,
			Signer:    model.SelfSignatureType,
		},
		{
			Signature: encodedSignature,
			Signer:    model.ApplicationSignatureType,
		},
	}

	err = validator.Validate(mock.StartNoopSpan(), signatureList, scopeWithAppID, &csrParameters)

	assert.NoError(t, err)
}

//
// validateCSRStamps :: for two self-signed CSR stamps :: returns an error.
//
func TestValidateCSRStampsForTwoSelfSignedCSRStamps(t *testing.T) {

	publicKeyBytes := []byte(originalString)
	encodedPublicKey := encodedString
	emptyIdentity := emptyString
	validVersion := model.CardVersion5
	emptyPreviousCardID := emptyString
	createdAt := time.Now().UTC().Unix()
	csr, err := getCSRMessageAsJSON(encodedPublicKey, emptyIdentity, validVersion, emptyPreviousCardID, createdAt)

	assert.NoError(t, err)

	signatureBytes := []byte("Signature")
	encodedSignature := encodeMessage(signatureBytes)
	csrParameters := ParametersStore{
		publicKey: publicKeyBytes,
		csr:       []byte(csr),
	}

	crypto := new(mock.Crypto)
	crypto.On("ValidateVirgilCardSignature",
		csrParameters.csr, []byte{}, publicKeyBytes, signatureBytes,
	).Return(nil)
	encoder, _ := presetEncoder(string(signatureBytes))
	validator := NewCSRStampsValidator(crypto, encoder)
	signatureList := []api.CSRStamp{
		{
			Signature: encodedSignature,
			Signer:    model.SelfSignatureType,
		},
		{
			Signature: encodedSignature,
			Signer:    model.SelfSignatureType,
		},
	}

	err = validator.Validate(mock.StartNoopSpan(), signatureList, emptyString, &csrParameters)

	assert.Error(t, err)
	assert.Equal(t, api.ErrSelfCSRStampMustBeUnique, err)
}

//
// validateCSRStamps :: for two app-signed CSR stamps :: passes if al least one of them is within ScopeIP.
//
func TestValidateSignatureList_WithTwoApplicationSignatures_ReturnsAnError(t *testing.T) {

	publicKeyBytes := []byte(originalString)
	encodedPublicKey := encodedString
	emptyIdentity := emptyString
	validVersion := model.CardVersion5
	emptyPreviousCardID := emptyString
	createdAt := time.Now().UTC().Unix()
	appID2 := strings.Repeat("c", models.IDLength)
	scopeIDWithSecondApp := appID2
	csr, err := getCSRMessageAsJSON(encodedPublicKey, emptyIdentity, validVersion, emptyPreviousCardID, createdAt)

	assert.NoError(t, err)

	signatureBytes := []byte("Signature")
	encodedSignature := encodeMessage(signatureBytes)
	csrParameters := ParametersStore{
		publicKey: publicKeyBytes,
		csr:       []byte(csr),
	}

	crypto := new(mock.Crypto)
	crypto.On("ValidateVirgilCardSignature",
		csrParameters.csr, []byte{}, publicKeyBytes, signatureBytes,
	).Return(nil)
	encoder, _ := presetEncoder(string(signatureBytes))
	validator := NewCSRStampsValidator(crypto, encoder)
	signatureList := []api.CSRStamp{
		{
			Signature: encodedSignature,
			Signer:    model.SelfSignatureType,
		},
		{
			Signature: encodedSignature,
			Signer:    model.ApplicationSignatureType,
		},
		{
			Signature: encodedSignature,
			Signer:    model.ApplicationSignatureType,
		},
	}

	err = validator.Validate(mock.StartNoopSpan(), signatureList, scopeIDWithSecondApp, &csrParameters)

	assert.NoError(t, err)
}

//
// validateCSRStamp :: for an invalid extra snapshot :: returns an error.
//
func TestValidateCSRStampForAnInvalidSnapshot(t *testing.T) {

	encoder := presetEncoderWithError(originalString, errors.New("error"))
	validator := NewCSRStampsValidator(&mock.Crypto{}, encoder)
	csrStamp := api.CSRStamp{Snapshot: encodedString}

	err := validator.validateCSRStamp(mock.StartNoopSpan(), csrStamp, new(ParametersStore))

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRStampExtraSnapshotDecoding, helper.ExtractHTTPError(err))
}

//
// validateCSRStamp :: for an invalid signature :: returns an error.
//
func TestValidateCSRStampForAnInvalidSignature(t *testing.T) {

	encoder, _ := presetEncoder(originalString)
	validator := NewCSRStampsValidator(&mock.Crypto{}, encoder)
	csrStamp := api.CSRStamp{
		Snapshot: encodedString,
		Signer:   model.ApplicationSignatureType,
	}

	err := validator.validateCSRStamp(mock.StartNoopSpan(), csrStamp, new(ParametersStore))

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRStampSignatureIsMissing, err)
}

//
// validateCSRStamp :: for an invalid signer type :: returns an error.
//
func TestValidateCSRStampForAnInvalidSignerType(t *testing.T) {

	encoder, _ := presetEncoder(originalString)
	validator := NewCSRStampsValidator(&mock.Crypto{}, encoder)
	csrStamp := api.CSRStamp{
		Snapshot:  encodedString,
		Signature: encodedString,
	}

	err := validator.validateCSRStamp(mock.StartNoopSpan(), csrStamp, new(ParametersStore))

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRStampSignerIsEmpty, err)
}

//
// validateCSRStamp :: for all valid parameters :: passes.
//
func TestValidateCSRStampForAllValidParameters(t *testing.T) {

	encoder, _ := presetEncoder(originalString)
	validator := NewCSRStampsValidator(&mock.Crypto{}, encoder)
	csrStamp := api.CSRStamp{
		Snapshot:  encodedString,
		Signature: encodedString,
		Signer:    model.ApplicationSignatureType,
	}

	err := validator.validateCSRStamp(mock.StartNoopSpan(), csrStamp, new(ParametersStore))

	assert.NoError(t, err)
}

//
// validateCSRStampSignature :: for an empty value :: returns an error.
//
func TestValidateCSRStampSignatureForAnEmptyValue(t *testing.T) {

	validator := NewCSRStampsValidator(&mock.Crypto{}, &mock.Base64Encoder{})

	err := validator.validateCSRStampSignature(
		mock.StartNoopSpan(),
		emptyString,
		emptyString,
		false,
		new(ParametersStore),
	)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRStampSignatureIsMissing, err)
}

//
// validateCSRStampSignature :: for an incorrectly encoded value :: returns an error.
//
func TestValidateCSRStampSignatureForAnIncorrectlyEncodedValue(t *testing.T) {

	encoder := presetEncoderWithError(originalString, errors.New("error"))
	validator := NewCSRStampsValidator(&mock.Crypto{}, encoder)

	err := validator.validateCSRStampSignature(
		mock.StartNoopSpan(),
		encodedString,
		emptyString,
		false,
		new(ParametersStore),
	)

	assert.Error(t, err)
	assert.Equal(t, api.ErrSignatureDecoding, helper.ExtractHTTPError(err))
}

//
// validateCSRStampSignature :: for a correctly encoded value and not a self-signature :: passes.
//
func TestValidateCSRStampSignatureForACorrectlyEncodedValueAndNotASelfSignature(t *testing.T) {

	encoder, _ := presetEncoder(originalString)
	validator := NewCSRStampsValidator(&mock.Crypto{}, encoder)

	err := validator.validateCSRStampSignature(
		mock.StartNoopSpan(),
		encodedString,
		emptyString,
		false,
		new(ParametersStore),
	)

	assert.NoError(t, err)
}

//
// validateCSRStampSignature :: for an invalid signature with an empty snapshot :: returns an error.
//
func TestValidateCSRStampSignatureForAnIncorrectSignatureAndSelfSignatureWithoutSnapshot(t *testing.T) {

	invalidSignature := originalString
	invalidEncodedSignature := encodedString
	publicKeyBytes := []byte("public key")
	csrBytes := []byte("csr bytes")
	csrParams := ParametersStore{
		csr:       csrBytes,
		publicKey: publicKeyBytes,
	}

	crypto := new(mock.Crypto)
	crypto.On("ValidateVirgilCardSignature",
		csrBytes, []byte{}, publicKeyBytes, []byte(originalString),
	).Return(errors.New("error"))
	encoder, _ := presetEncoder(invalidSignature)
	validator := NewCSRStampsValidator(crypto, encoder)

	err := validator.validateCSRStampSignature(
		mock.StartNoopSpan(),
		invalidEncodedSignature,
		emptyString,
		true,
		&csrParams,
	)

	assert.Error(t, err)
	assert.Equal(t, api.ErrSignatureVerificationFailed, helper.ExtractHTTPError(err))
}

//
// validateCSRStampSignature :: for an invalid signature with not empty snapshot :: returns an error.
//
func TestValidateCSRStampSignatureForAnIncorrectSignatureAndSelfSignatureWithSnapshot(t *testing.T) {

	invalidSignature := originalString
	invalidEncodedSignature := encodedString
	publicKeyBytes := []byte("public key")
	csrBytes := []byte("csr bytes")
	extraSnapshotBytes := []byte("snapshot bytes")
	encodedExtraSnapshot := encodeMessage(extraSnapshotBytes)
	csrParams := ParametersStore{
		csr:       csrBytes,
		publicKey: publicKeyBytes,
	}

	crypto := new(mock.Crypto)
	crypto.On("ValidateVirgilCardSignature",
		csrBytes, extraSnapshotBytes, publicKeyBytes, []byte(originalString),
	).Return(errors.New("error"))
	encoder, _ := presetEncoder(invalidSignature)
	encoder.On("DecodeString", encodedExtraSnapshot).Return(extraSnapshotBytes, nil)
	validator := NewCSRStampsValidator(crypto, encoder)

	err := validator.validateCSRStampSignature(
		mock.StartNoopSpan(),
		invalidEncodedSignature,
		encodedExtraSnapshot,
		true,
		&csrParams,
	)

	assert.Error(t, err)
	assert.Equal(t, api.ErrSignatureVerificationFailed, helper.ExtractHTTPError(err))
}

//
// validateCSRStampSignature :: for all valid parameters :: passes.
//
func TestValidateCSRStampSignatureForAllValidParameters(t *testing.T) {

	validSignature := originalString
	validEncodedSignature := encodedString
	publicKeyBytes := []byte("public key")
	csrBytes := []byte("csr bytes")
	extraSnapshotBytes := []byte("snapshot bytes")
	encodedExtraSnapshot := encodeMessage(extraSnapshotBytes)
	csrParams := ParametersStore{
		csr:       csrBytes,
		publicKey: publicKeyBytes,
	}

	crypto := new(mock.Crypto)
	crypto.On("ValidateVirgilCardSignature",
		csrBytes, extraSnapshotBytes, publicKeyBytes, []byte(originalString),
	).Return(nil)
	encoder, _ := presetEncoder(validSignature)
	encoder.On("DecodeString", encodedExtraSnapshot).Return(extraSnapshotBytes, nil)
	validator := NewCSRStampsValidator(crypto, encoder)

	err := validator.validateCSRStampSignature(
		mock.StartNoopSpan(),
		validEncodedSignature,
		encodedExtraSnapshot,
		true,
		&csrParams,
	)

	assert.NoError(t, err)
	assert.Equal(t, extraSnapshotBytes, csrParams.extraSnapshot)
}

//
// validateCSRStampSignature :: for all valid parameters :: passes.
//
func TestValidateCSRStampSignatureForAllValidParametersAndWithoutMocks(t *testing.T) {

	signatureBytes := []byte("signature")
	encodedSignature := encodeMessage(signatureBytes)
	csrBytes := []byte("csr bytes")
	extraCSRBytes := []byte("snapshot bytes")
	keyPair, err := cryptoimpl.NewKeypair()

	assert.Nil(t, err)

	publicKeyBytes, err := keyPair.PublicKey().Encode()

	assert.Nil(t, err)

	encodedExtraCSR := encodeMessage(extraCSRBytes)
	crypto := new(mock.Crypto)
	crypto.On(
		"ValidateVirgilCardSignature",
		csrBytes,
		extraCSRBytes,
		publicKeyBytes,
		signatureBytes,
	).Return(nil)
	csrParams := ParametersStore{
		csr:       csrBytes,
		publicKey: publicKeyBytes,
	}

	encoder := mock.Base64Encoder{}
	encoder.On("DecodeString", encodedSignature).Return(signatureBytes, nil)
	encoder.On("DecodeString", encodedExtraCSR).Return(extraCSRBytes, nil)

	validator := NewCSRStampsValidator(crypto, &encoder)

	err = validator.validateCSRStampSignature(
		mock.StartNoopSpan(),
		encodedSignature,
		encodedExtraCSR,
		true,
		&csrParams,
	)

	assert.NoError(t, err)
	assert.Equal(t, extraCSRBytes, csrParams.extraSnapshot)
}

//
//  validateCSRSigner :: with a too long signer type :: returns an error.
//
func TestValidateCSRStampTypeWithAToLongSignerType(t *testing.T) {

	validator := NewCSRStampsValidator(&mock.Crypto{}, &mock.Base64Encoder{})

	err := validator.validateCSRSigner(strings.Repeat("a", CSRStampSignerMaxLength+1))

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRStampSignerIsTooLong, err)
}

//
//  validateCSRSigner :: without a signature type :: returns an error.
//
func TestValidateCSRStampTypeWithAnEmptySignatureType(t *testing.T) {

	validator := NewCSRStampsValidator(&mock.Crypto{}, &mock.Base64Encoder{})

	err := validator.validateCSRSigner("")

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRStampSignerIsEmpty, err)
}

//
//  validateCSRSigner :: for a custom signer type :: passes.
//
func TestValidateCSRStampTypeWithACustomSignerType(t *testing.T) {

	validator := NewCSRStampsValidator(&mock.Crypto{}, &mock.Base64Encoder{})

	err := validator.validateCSRSigner("INCORRECT TYPE")

	assert.NoError(t, err)
}

//
// validateCSRSigner :: for an application signature type :: passes.
//
func TestValidateCSRStampTypeWithAnApplicationSignatureType(t *testing.T) {

	validator := NewCSRStampsValidator(&mock.Crypto{}, &mock.Base64Encoder{})

	err := validator.validateCSRSigner(model.ApplicationSignatureType)

	assert.NoError(t, err)
}

//
// validateCSRSigner :: for a self signature type :: passes.
//
func TestValidateCSRStampTypeWithASelfSignatureType(t *testing.T) {

	validator := NewCSRStampsValidator(&mock.Crypto{}, &mock.Base64Encoder{})

	err := validator.validateCSRSigner(model.SelfSignatureType)

	assert.NoError(t, err)
}

//
// validateCSRSigner :: for a Virgil signature type :: returns an error.
//
func TestValidateCSRStampTypeWithAVirgilSignatureType(t *testing.T) {

	validator := NewCSRStampsValidator(&mock.Crypto{}, &mock.Base64Encoder{})

	err := validator.validateCSRSigner(model.VirgilSignatureType)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRStampSignerIsIncorrect, err)
}

//
//  validateCSRStampSnapshot :: for incorrect extra snapshot :: returns an error.
//
func TestValidateCSRStampSnapshotForIncorrectExtraSnapshot(t *testing.T) {

	decodeErr := errors.New("decoding error")
	encoder := presetEncoderWithError(originalString, decodeErr)
	validator := NewCSRStampsValidator(&mock.Crypto{}, encoder)

	err := validator.validateCSRStampSnapshot(encodedString)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRStampExtraSnapshotDecoding, helper.ExtractHTTPError(err))
}

//
//  validateCSRStampSnapshot :: for a too long extra snapshot :: returns an error.
//
func TestValidateCSRStampSnapshotForATooLongExtraSnapshot(t *testing.T) {

	tooLongExtraSnapshot := make([]byte, CSRStampSnapshotMaxLength+1)
	encoder, tooLongEncodedExtraSnapshot := presetEncoder(string(tooLongExtraSnapshot))
	validator := NewCSRStampsValidator(&mock.Crypto{}, encoder)

	err := validator.validateCSRStampSnapshot(tooLongEncodedExtraSnapshot)

	assert.Error(t, err)
	assert.Equal(t, api.ErrExtraContentSnapshotIsTooLong, err)
}

//
//  validateCSRStampSnapshot :: for a valid extra snapshot :: passes.
//
func TestValidateCSRStampSnapshotForAValidExtraSnapshot(t *testing.T) {

	encoder, _ := presetEncoder(originalString)
	validator := NewCSRStampsValidator(&mock.Crypto{}, encoder)

	err := validator.validateCSRStampSnapshot(encodedString)

	assert.NoError(t, err)
}
