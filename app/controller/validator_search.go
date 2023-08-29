package controller

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"

	"github.com/VirgilSecurity/virgil-services-cards/src/api"
)

//
// Validator performs card search request validations.
//
const (
	searchIdentitiesLimit = 50
)

//
// SearchCardValidatorProvider provider and interface to wot with card search validator.
//
type SearchCardValidatorProvider interface {
	//
	// Validate validates Virgil Card search request.
	//
	Validate(span tracer.Span, cardSearchRequest *api.CardSearchRequest) (err error)
}

//
// SearchCardValidator performs card search request validations.
//
type SearchCardValidator struct{}

//
// NewSearchCardValidator creates new instance of Search Card validator.
//
func NewSearchCardValidator() *SearchCardValidator {

	return &SearchCardValidator{}
}

//
// Validate performs a validation of the search Virgil Card request object.
//
func (v *SearchCardValidator) Validate(span tracer.Span, cardSearchRequest *api.CardSearchRequest) error {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentValidator,
		},
	)
	defer span.Finish()

	identitiesCount := len(cardSearchRequest.GetIdentities())
	if identitiesCount == 0 {
		return tracer.SetSpanErrorAndReturn(span, api.ErrIdentitySearchTermCannotBeEmpty)
	} else if searchIdentitiesLimit < identitiesCount {
		return tracer.SetSpanErrorAndReturn(span, api.ErrIdentitySearchCountIsLimited(searchIdentitiesLimit))
	}

	return nil
}
