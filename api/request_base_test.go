package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/VirgilSecurity/virgil-services-cards/src/model"
)

//
// Testing constants.
//
const (
	encodedString = "b3JpZ2luYWwgc3RyaW5n"
)

//
// TestGetContentSnapshotForAnEmptyObject :: for an empty content snapshot :: returns an empty value.
//
func TestGetContentSnapshotForAnEmptyObject(t *testing.T) {

	req := CardBaseRequest{}

	s := req.GetCSR()

	assert.Empty(t, s)
}

//
// TestGetContentSnapshotForNotEmptyObject :: for a content snapshot with value :: returns set value.
//
func TestGetContentSnapshotForNotEmptyObject(t *testing.T) {

	req := CardBaseRequest{CSR: encodedString}

	s := req.GetCSR()

	assert.Equal(t, encodedString, s)
}

//
// TestGetSelfStampForAnEmptyObject :: for an empty object :: returns an error.
//
func TestGetSelfStampForAnEmptyObject(t *testing.T) {

	req := CardBaseRequest{}

	selfStamp, err := req.GetSelfStamp()

	assert.Error(t, err)
	assert.Empty(t, selfStamp)
	assert.Equal(t, ErrSelfCSRStampIsMissing, err)
}

//
// TestGetSelfStampForAnObjectWithASelfStamp :: for an object with a self stamp :: returns a self stamp.
//
func TestGetSelfStampForAnObjectWithASelfStamp(t *testing.T) {

	signer := model.SelfSignatureType
	req := CardBaseRequest{
		CSRStamps: []CSRStamp{
			{
				Signer: signer,
			},
		},
	}

	selfStamp, err := req.GetSelfStamp()

	assert.Empty(t, err)
	assert.Equal(t, signer, selfStamp.GetSigner())
}
