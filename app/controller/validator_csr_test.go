package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/VirgilSecurity/virgil-services-core-kit/models"
	"github.com/VirgilSecurity/virgil-services-core-kit/test/helper"

	"github.com/VirgilSecurity/virgil-services-cards/src/api"
	"github.com/VirgilSecurity/virgil-services-cards/src/dao"
	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto"
	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/encoder"
	"github.com/VirgilSecurity/virgil-services-cards/src/model"
	"github.com/VirgilSecurity/virgil-services-cards/test/mock"
)

//
// Testing constants.
//
const (
	emptyString          = ""
	emptyScopeID         = emptyString
	emptyRequestIdentity = emptyString
	validCardVersion     = model.CardVersion5
	notJSONMessage       = "some message not a JSON message"
	emptyJSONMessage     = "{}"
	validIdentity        = "Valid Identity"
	validPublicKey       = "valid and correct public key "
	encodedPublicKey     = "dmFsaWQgYW5kIGNvcnJlY3QgcHVibGljIGtleSA="
	validID              = "e680bef87ba75d331b0a02bfa6a20f02eb5c5ba9bc96fc61ca595404b10026f4"
	invalidBase64String  = "!@#$%^&*()"
)

//
// Testing variables.
//
var (
	publicKeyBytes   = []byte(validPublicKey)
	publicKeyEncoded = encodeMessage(publicKeyBytes)
)

//
// validateCSR :: for an empty CSR :: returns an error.
//
func TestValidateCSRWithEmptyCSRMessage(t *testing.T) {

	validator := getCSRValidatorUnderTest(validatorDeps{})

	decodedCSRParams, err := validator.Validate(mock.StartNoopSpan(), emptyString, &api.CSR{}, emptyRequestIdentity)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRIsEmpty, err)
	assert.Empty(t, decodedCSRParams)
}

//
// validateCSR :: for an incorrectly encoded message :: returns an error.
//
func TestValidateCSRWithMalformedCSRMessage(t *testing.T) {

	encoder := new(mock.Base64Encoder)
	decodeErr := errors.New("decode error")
	encoder.On("DecodeString", invalidBase64String).Return([]byte{}, decodeErr)
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedCSRParams, err := validator.Validate(
		mock.StartNoopSpan(),
		invalidBase64String,
		&api.CSR{},
		emptyRequestIdentity,
	)

	assert.Error(t, err)
	assert.Equal(t, api.ErrContentSnapshotIsNotABase64EncodedString, helper.ExtractHTTPError(err))
	assert.Empty(t, decodedCSRParams)
}

//
// validateCSR :: for a not JSON message :: returns an error.
//
func TestValidateCSRWithANotJSONCSRMessage(t *testing.T) {

	encoder, encodedMsg := presetEncoder(notJSONMessage)
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedCSRParams, err := validator.Validate(
		mock.StartNoopSpan(),
		encodedMsg,
		&api.CSR{},
		emptyRequestIdentity,
	)

	assert.Error(t, err)
	assert.Equal(t, api.ErrContentSnapshotIsNotAJSONMessage, helper.ExtractHTTPError(err))
	assert.Empty(t, decodedCSRParams)
}

//
// validateCSR :: for an incorrect timestamp parameter :: returns a proper error.
//
func TestValidateCSRWithInvalidCreatedAtFieldType(t *testing.T) {

	incorrectMessage := fmt.Sprintf(`{"%s": "145"}`, api.CreatedAtFieldName)
	encoder, encodedMsg := presetEncoder(incorrectMessage)
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedCSRParams, err := validator.Validate(mock.StartNoopSpan(), encodedMsg, &api.CSR{}, validIdentity)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRCreationTimeIsIncorrect, helper.ExtractHTTPError(err))
	assert.Empty(t, decodedCSRParams)
}

//
// validateCSR :: for an empty JSON message :: returns an error.
//
func TestValidateCSRForAnEmptyJSONMessage(t *testing.T) {

	encoder, encodedMsg := presetEncoder(emptyJSONMessage)
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedCSRParams, err := validator.Validate(mock.StartNoopSpan(), encodedMsg, &api.CSR{}, emptyRequestIdentity)

	assert.Error(t, err)
	assert.Empty(t, decodedCSRParams)
	assert.Equal(t, api.ErrCSRIdentityIsEmpty, err)
}

//
// validateCSR :: for an invalid public key :: returns an error.
//
func TestValidateCSRForAnInvalidPublicKey(t *testing.T) {

	csrMsg, err := getCSRMessageAsJSON(invalidBase64String, "", "", "", 0)

	assert.NoError(t, err)

	encoder, encodedMsg := presetEncoder(csrMsg)
	errPublicKeyDecode := errors.New("public key decode error")
	encoder.On("DecodeString", invalidBase64String).Return([]byte{}, errPublicKeyDecode)
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedCSRParams, err := validator.Validate(mock.StartNoopSpan(), encodedMsg, &api.CSR{}, emptyRequestIdentity)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRPublicKeyDecoding, helper.ExtractHTTPError(err))
	assert.Empty(t, decodedCSRParams)
}

//
// TestValidateCSRPublicKeyForATooLongPublicKey :: for a too short public key :: returns an error
//
func TestValidateCSRPublicKeyForATooLongPublicKey(t *testing.T) {

	tooLongPublicKey := make([]byte, PublicKeyMaxLength+1)
	encoder, encodedPublicKey := presetEncoder(string(tooLongPublicKey))
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedKey, err := validator.decodeCSRPublicKeyAndCheckMaxLength(encodedPublicKey)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRPublicKeyIsTooLong, err)
	assert.Empty(t, decodedKey)
}

//
// validateCSR :: for a too long identity :: returns an error.
//
func TestValidateCSRForATooLongIdentity(t *testing.T) {

	tooLongIdentity := strings.Repeat("a", IdentityMaxLength+1)
	csrMsg, err := getCSRMessageAsJSON(encodedPublicKey, tooLongIdentity, "", "", 0)

	assert.NoError(t, err)

	encoder, encodedMsg := presetEncoder(csrMsg)
	encoder.On("DecodeString", encodedPublicKey).Return(publicKeyBytes, nil)
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedCSRParams, err := validator.Validate(mock.StartNoopSpan(), encodedMsg, &api.CSR{}, emptyRequestIdentity)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRIdentityIsIncorrect, err)
	assert.Empty(t, decodedCSRParams)
}

//
// validateCSR :: for an invalid previous card ID :: returns an error.
//
func TestValidateCSRForAnInvalidPreviousCardID(t *testing.T) {

	requestIdentity := validIdentity
	invalidPreviousCardID := "invalid card ID"
	csrMsg, err := getCSRMessageAsJSON(encodedPublicKey, validIdentity, "", invalidPreviousCardID, 0)

	assert.NoError(t, err)

	encoder, encodedMsg := presetEncoder(csrMsg)
	encoder.On("DecodeString", encodedPublicKey).Return(publicKeyBytes, nil)
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedCSRParams, err := validator.Validate(mock.StartNoopSpan(), encodedMsg, &api.CSR{}, requestIdentity)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRPreviousCardIDIsIncorrect, err)
	assert.Empty(t, decodedCSRParams)
}

//
// validateCSR :: for an invalid create_at :: returns an error.
//
func TestValidateCSRWithAnInvalidCreateAtValue(t *testing.T) {

	requestIdentity := validIdentity
	csrMsg, err := getCSRMessageAsJSON(encodedPublicKey, validIdentity, "", "", 0)

	assert.NoError(t, err)

	encoder, encodedMsg := presetEncoder(csrMsg)
	encoder.On("DecodeString", encodedPublicKey).Return(publicKeyBytes, nil)
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedCSRParams, err := validator.Validate(mock.StartNoopSpan(), encodedMsg, &api.CSR{}, requestIdentity)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRCreationTimeIsIncorrect, err)
	assert.Empty(t, decodedCSRParams)
}

//
// validateCSR :: for an invalid version :: returns an error.
//
func TestValidateCSRForAnInvalidVersion(t *testing.T) {

	now := time.Now().UTC().Unix()
	invalidVersion := "v3"
	requestIdentity := validIdentity
	csrMsg, err := getCSRMessageAsJSON(encodedPublicKey, validIdentity, invalidVersion, "", now)

	assert.NoError(t, err)

	encoder, encodedMsg := presetEncoder(csrMsg)
	encoder.On("DecodeString", encodedPublicKey).Return(publicKeyBytes, nil)
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedCSRParams, err := validator.Validate(mock.StartNoopSpan(), encodedMsg, &api.CSR{}, requestIdentity)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRVersionIsIncorrect, err)
	assert.Empty(t, decodedCSRParams)
}

//
// TestValidateCSRForAnEmptyIdentity :: for an empty identity :: returns an error.
//
func TestValidateCSRForAnEmptyIdentity(t *testing.T) {

	csrMsg, err := getCSRMessageAsJSON(encodedPublicKey, "", validCardVersion, "", time.Now().UTC().Unix())

	assert.NoError(t, err)

	encoder, encodedMsg := presetEncoder(csrMsg)
	encoder.On("DecodeString", encodedPublicKey).Return(publicKeyBytes, nil)
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedCSRParams, err := validator.Validate(mock.StartNoopSpan(), encodedMsg, &api.CSR{}, emptyRequestIdentity)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRIdentityIsEmpty, err)
	assert.Empty(t, decodedCSRParams)

}

//
// TestValidateCSRPublicKeyForAValidPublicKey :: for a valid public key :: passes.
//
func TestValidateCSRPublicKeyForAValidPublicKey(t *testing.T) {
	encoder, encodedMsg := presetEncoder(string(publicKeyBytes))
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedKey, err := validator.decodeCSRPublicKeyAndCheckMaxLength(encodedMsg)

	assert.NoError(t, err)
	assert.Equal(t, publicKeyBytes, decodedKey)
}

//
// validateCSR :: for all valid parameters :: passes.
//
func TestValidateCSRForAllValidParameters(t *testing.T) {

	now := time.Now().UTC().Unix()
	requestIdentity := validIdentity
	csrMsg, err := getCSRMessageAsJSON(encodedPublicKey, validIdentity, validCardVersion, "", now)

	assert.NoError(t, err)

	encoder, encodedMsg := presetEncoder(csrMsg)
	encoder.On("DecodeString", encodedPublicKey).Return(publicKeyBytes, nil)
	validator := getCSRValidatorUnderTest(validatorDeps{encoder: encoder})

	decodedCSRParams, err := validator.Validate(mock.StartNoopSpan(), encodedMsg, &api.CSR{}, requestIdentity)

	assert.Empty(t, err)
	assert.Equal(t, []byte(csrMsg), decodedCSRParams.csr)
	assert.Equal(t, publicKeyBytes, decodedCSRParams.publicKey)
}

//
// validateCSRIdentity :: for a valid not empty identity :: passes.
//
func TestValidateCSRIdentityForAValidNotEmptyValue(t *testing.T) {

	validator := getCSRValidatorUnderTest(validatorDeps{})
	requestIdentity := validIdentity

	err := validator.validateCSRIdentity(validIdentity, requestIdentity)

	assert.NoError(t, err)
}

//
// validateCSRIdentity :: for a valid not empty identity other than request one :: returns an error.
//
func TestValidateCSRIdentityForAValueOtherThanRequestOne(t *testing.T) {

	validator := getCSRValidatorUnderTest(validatorDeps{})
	requestIdentity := "Some other identity"

	err := validator.validateCSRIdentity(validIdentity, requestIdentity)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRIdentityDoesNotMatchRequestIdentity, err)
}

//
// validateCSRIdentity :: for an empty identity :: returns an error.
//
func TestValidateCSRIdentityForAnEmptyValue(t *testing.T) {

	validator := getCSRValidatorUnderTest(validatorDeps{})

	err := validator.validateCSRIdentity("", emptyRequestIdentity)

	assert.Error(t, err)
	assert.Equal(t, api.ErrCSRIdentityIsEmpty, err)
}

//
// validateCSRPreviousCardID :: for an empty value :: passes.
//
func TestValidateCSRPreviousCardIDForAnEmptyValue(t *testing.T) {

	validator := getCSRValidatorUnderTest(validatorDeps{})

	err := validator.validateCSRPreviousCardID("")

	assert.NoError(t, err)
}

//
// validateCSRPreviousCardID :: for a valid value :: passes.
//
func TestValidateCSRPreviousCardIDForAValidValue(t *testing.T) {

	validPreviousCardID := strings.Repeat("a", models.IDLength)
	validator := getCSRValidatorUnderTest(validatorDeps{})

	err := validator.validateCSRPreviousCardID(validPreviousCardID)

	assert.NoError(t, err)
}

//
// validateCSRCreatedAt :: for a valid value :: passes.
//
func TestValidateCSRCreatedAtForAValidValue(t *testing.T) {

	validCreationTimestamp := time.Now().UTC().Unix()
	validator := getCSRValidatorUnderTest(validatorDeps{})

	err := validator.validateCSRCreatedAt(validCreationTimestamp)

	assert.NoError(t, err)
}

//
// validateCSRVersion :: for a valid value :: passes.
//
func TestValidateCSRVersionForAValidValue(t *testing.T) {

	validator := getCSRValidatorUnderTest(validatorDeps{})

	err := validator.validateCSRVersion(validCardVersion)

	assert.NoError(t, err)
}

//
// getCSRValidatorUnderTest returns validator instance.
//
func getCSRValidatorUnderTest(deps validatorDeps) *CSRValidator {

	return NewCSRValidator(deps.encoder)
}

//
// presetEncoder returns an encoder mock with preset valid response.
//
func presetEncoder(message string) (*mock.Base64Encoder, string) {

	msg := []byte(message)
	encodedMsg := encodeMessage(msg)
	encoder := new(mock.Base64Encoder)
	encoder.On("DecodeString", encodedMsg).Return(msg, nil)

	return encoder, encodedMsg
}

//
// presetEncoderWithError returns an encoder mock with preset error response.
//
func presetEncoderWithError(message string, err error) *mock.Base64Encoder {

	msg := []byte(message)
	encodedMsg := encodeMessage(msg)
	encoder := new(mock.Base64Encoder)
	encoder.On("DecodeString", encodedMsg).Return([]byte{}, err)

	return encoder
}

//
// encodeMessage performs a base64-encoding.
//
func encodeMessage(content []byte) string {

	return encoder.NewBase64().EncodeToString(content)
}

//
// getEncodedCSRMessage returns a serialized JSON message.
//
func getCSRMessageAsJSON(publicKey, identity, version, previousCardID string, createdAt int64) (string, error) {

	jsonMessage, err := json.Marshal(api.CSR{
		PublicKey:      publicKey,
		Identity:       identity,
		PreviousCardID: previousCardID,
		Version:        version,
		CreatedAt:      createdAt,
	})

	return string(jsonMessage), err
}

//
// validatorDeps is a validator dependencies structure.
//
type validatorDeps struct {
	crypto         crypto.Provider
	encoder        encoder.Provider
	cardRepository dao.CardRepositoryProvider
}
