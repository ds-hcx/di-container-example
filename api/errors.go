package api

import (
	"fmt"

	"github.com/VirgilSecurity/virgil-services-core-kit/errors"
)

//
// Core errors.
//
var (
	// Bulk internal errors.
	ErrInternalError = errors.NewHTTP500Error(
		10000,
		"Request serving internal error. Try again later.",
	)
	ErrNotFound = errors.NewHTTP404Error(
		10001,
		"Requested card entity not found.",
	)
	ErrApplicationIDHeaderIsNotSet = errors.NewHTTP400Error(
		20310,
		"Request scope application is not set.",
	)
	ErrIdentityHeaderNotSet = errors.NewHTTP400Error(
		20311,
		"Request identity is not set.",
	)
	// Application errors.
	ErrRequestParsing = errors.NewHTTP400Error(
		30001,
		"Request body parsing error. Invalid JSON, field name or field type.",
	)
)

//
// Base Card Request common parsing and validation errors.
//
var (
	// Common card request errors.
	ErrCSRIsEmpty = errors.NewHTTP400Error(
		40001,
		"Content snapshot is empty.",
	)
	ErrContentSnapshotIsNotABase64EncodedString = errors.NewHTTP400Error(
		40002,
		"Content snapshot is not a base64-encoded string.",
	)
	ErrContentSnapshotIsNotAJSONMessage = errors.NewHTTP400Error(
		40003,
		"Content snapshot is not a base64-encoded JSON message.",
	)
	ErrCSRStampSignatureIsMissing = errors.NewHTTP400Error(
		40006,
		"Signature is missing in signature entries.",
	)
	ErrCSRStampSignerIsIncorrect = errors.NewHTTP400Error(
		40007,
		"Signer is missing or is incorrect in one of signature entries.",
	)
	ErrCSRStampSignerIsEmpty = errors.NewHTTP400Error(
		40035,
		"Signer is empty in one of signature entries.",
	)
	ErrCSRStampSignerIsTooLong = errors.NewHTTP400Error(
		40036,
		"Signer is too long in one of signature entries. It mustn't exceed 1024 characters.",
	)
	ErrSelfCSRStampIsMissing = errors.NewHTTP400Error(
		40008,
		"Self signature is missing for the Virgil Card.",
	)
	ErrSelfCSRStampMustBeUnique = errors.NewHTTP400Error(
		40031,
		"Self signature must be only one.",
	)
	ErrCSRVersionIsIncorrect = errors.NewHTTP400Error(
		40011,
		"Virgil Card version must be 5.0.",
	)
	ErrCSRPublicKeyDecoding = errors.NewHTTP400Error(
		40012,
		"Public key is not a base64-encoded string.",
	)
	ErrCSRPreviousCardIDIsIncorrect = errors.NewHTTP400Error(
		40014,
		"Previous Virgil Card ID is not a valid ID.",
	)
	ErrPreviousVirgilCardDoesNotExist = errors.NewHTTP400Error(
		40015,
		"Previous Virgil Card ID does not exist.",
	)
	ErrPreviousVirgilCardIsRegisteredForAnotherScope = errors.NewHTTP400Error(
		40016,
		"Previous Virgil Card ID is registered for another application.",
	)
	ErrPreviousVirgilCardIdentityIsIncorrect = errors.NewHTTP400Error(
		40032,
		"Previous Virgil Card identity doesn't match current Virgil Card one.",
	)
	ErrPreviousVirgilCardExistsAlready = errors.NewHTTP400Error(
		40037,
		"Previous Virgil Card exists already.",
	)
	ErrCSRIdentityIsIncorrect = errors.NewHTTP400Error(
		40017,
		"Identity is incorrect. It mustn't exceed 1024 bytes.",
	)
	ErrCSRIdentityDoesNotMatchRequestIdentity = errors.NewHTTP400Error(
		40034,
		"Identity is incorrect. It must match request identity provided in JWT.",
	)
	ErrCSRIdentityIsEmpty = errors.NewHTTP400Error(
		40033,
		"Identity is empty.",
	)
	ErrCSRCreationTimeIsIncorrect = errors.NewHTTP400Error(
		40018,
		"Creation time is incorrect.",
	)
	ErrCSRStampExtraSnapshotDecoding = errors.NewHTTP400Error(
		40019,
		"Extra content snapshot decoding error.",
	)
	ErrSignatureDecoding = errors.NewHTTP400Error(
		40020,
		"Signature is not a base64-encoded string.",
	)
	ErrSignatureVerificationFailed = errors.NewHTTP400Error(
		40021,
		"Signature verification failed for one of signature entries.",
	)
	ErrCSRStampsListIsTooSmall = errors.NewHTTP400Error(
		40022,
		"Signature list must contain at least self signature.",
	)
	ErrCSRStampsListIsTooLarge = errors.NewHTTP400Error(
		40023,
		"Signature list must contain at most eight entries.",
	)
	ErrExtraContentSnapshotIsTooLong = errors.NewHTTP400Error(
		40026,
		"Extra snapshot is too long for one of signature entries. It must not exceed 1024 bytes.",
	)
	ErrVirgilCardContentSnapshotIsNotUnique = errors.NewHTTP400Error(
		40027,
		"Virgil card content snapshot is not unique.",
	)
	ErrCSRPublicKeyIsTooLong = errors.NewHTTP400Error(
		40029,
		"Public key exceeds 4096 bytes.",
	)
)

//
// Card Create errors.
//
var (
	ErrCSRPublicKeyIsTooShort = errors.NewHTTP400Error(
		40030,
		"Public key is less than 16 bytes.",
	)
)

//
// Card Get handler errors.
//
var (
	ErrVirgilCardApplicationIDIsNotInTheAuthApplicationList = errors.NewHTTP403Error(
		40100,
		"Trying to get the Virgil Card that is scoped for another application.",
	)
)

//
// Card Search handler errors.
//
var (
	ErrIdentitySearchTermCannotBeEmpty = errors.NewHTTP400Error(
		40200,
		"Identity search parameter cannot be empty.",
	)
	ErrIdentitySearchCountIsLimited = func(searchIdentitiesLimit int) errors.HTTPError {
		return errors.NewHTTP400Error(
			40300,
			fmt.Sprintf("Identities to search amount limited to %d.", searchIdentitiesLimit),
		)
	}
)

//
// Card Delete handler errors.
//
var (
	// Delete Card specific errors
	ErrCSRPublicKeyMustBeEmpty = errors.NewHTTP400Error(
		40410,
		"Public key must be empty.",
	)
	ErrCSRPrevCardIDMustNotBeEmpty = errors.NewHTTP400Error(
		40420,
		"Empty card ID to delete (previous Card ID).",
	)
	ErrChainAlreadyDeleted = errors.NewHTTP403Error(
		40310,
		"Deleted card can not be deleted.",
	)
)
