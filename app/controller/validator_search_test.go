package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/VirgilSecurity/virgil-services-cards/src/api"
	"github.com/VirgilSecurity/virgil-services-cards/test/mock"
)

//
// TestValidateSearchRequestForSingleIdentity :: for valid Search request for singe identity :: passes.
//
func TestValidateSearchRequestForSingleIdentity(t *testing.T) {

	validator := NewSearchCardValidator()

	err := validator.Validate(mock.StartNoopSpan(), &api.CardSearchRequest{
		Identity: "identity",
	})

	assert.NoError(t, err)
}

//
// TestValidateSearchRequestForMultipleIdentities :: for valid Search request for multiple identities :: passes.
//
func TestValidateSearchRequestForMultipleIdentities(t *testing.T) {

	validator := NewSearchCardValidator()

	err := validator.Validate(mock.StartNoopSpan(), &api.CardSearchRequest{
		Identities: []string{"identity"},
	})

	assert.NoError(t, err)
}

//
// TestValidateSearchRequestWithEmptyIdentities :: for empty Search request :: returns an error.
//
func TestValidateSearchRequestWithEmptyIdentities(t *testing.T) {

	validator := NewSearchCardValidator()

	err := validator.Validate(mock.StartNoopSpan(), &api.CardSearchRequest{})

	assert.Error(t, err)
	assert.Equal(t, api.ErrIdentitySearchTermCannotBeEmpty, err)
}

//
// TestValidateSearchRequestWithWrongCountOfIdentities :: for wrong count fo Search identities :: returns an error.
//
func TestValidateSearchRequestWithWrongCountOfIdentities(t *testing.T) {

	validator := NewSearchCardValidator()

	err := validator.Validate(mock.StartNoopSpan(), &api.CardSearchRequest{
		Identities: func() []string {
			var identities []string
			for i := 0; i <= searchIdentitiesLimit; i++ {
				identities = append(identities, "identity")
			}
			return identities
		}(),
	})

	assert.Error(t, err)
	assert.Equal(t, api.ErrIdentitySearchCountIsLimited(searchIdentitiesLimit), err)
}
