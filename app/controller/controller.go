package controller

import (
	"time"

	"github.com/VirgilSecurity/virgil-services-core-kit/db/cassandra"
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"

	"github.com/VirgilSecurity/virgil-services-cards/src/api"
	"github.com/VirgilSecurity/virgil-services-cards/src/dao"
	"github.com/VirgilSecurity/virgil-services-cards/src/model"
)

//
// Provider provides an interface to work with Cards controller.
//
type Provider interface {
	//
	// CardCreate is a handler for POST /card request.
	//
	CardCreate(span tracer.Span, request *api.CardCreateRequest) (*model.CardDTO, error)

	//
	// CardGet is a handler for GET /card/:card_id request.
	//
	CardGet(span tracer.Span, request *api.CardBaseRequest, cardID string) (*model.CardDTO, error)

	//
	// CardSearch is a handler for POST /card/actions/search request.
	//
	CardSearch(span tracer.Span, request *api.CardSearchRequest) ([]*model.CardDTO, error)

	//
	// CardDelete is a handler for POST /card/actions/delete request.
	//
	CardDelete(span tracer.Span, request *api.CardDeleteRequest) (*model.CardDTO, error)
}

//
// Controller serves the Cards requests.
//
type Controller struct {
	cardSigner          model.CardSigner
	cardRepository      dao.CardRepositoryProvider
	createCardValidator CreateCardValidatorProvider
	searchCardValidator SearchCardValidatorProvider
	deleteCardValidator DeleteCardValidatorProvider
}

//
// New returns an instance of the Cards controller.
//
func New(
	cardSigner model.CardSigner,
	cardRepository dao.CardRepositoryProvider,
	createCardValidator CreateCardValidatorProvider,
	searchCardValidator SearchCardValidatorProvider,
	deleteCardValidator DeleteCardValidatorProvider,
) *Controller {

	return &Controller{
		cardSigner:          cardSigner,
		cardRepository:      cardRepository,
		createCardValidator: createCardValidator,
		searchCardValidator: searchCardValidator,
		deleteCardValidator: deleteCardValidator,
	}
}

//
// CardCreate is a handler for POST /card request.
//
func (h *Controller) CardCreate(span tracer.Span, request *api.CardCreateRequest) (*model.CardDTO, error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentController,
		},
	)
	defer span.Finish()

	virgilCard := model.NewCardDTO()
	if err := h.createCardValidator.Validate(span, request, virgilCard); nil != err {
		return nil, err
	}

	// TODO just for backward compatibility for JS SDK. Should be removed in future!.
	//  This check should be in transport/helper.go:63
	if request.UserID == "" {
		return nil, tracer.SetSpanErrorAndReturn(span, api.ErrIdentityHeaderNotSet)
	}

	if err := h.cardRepository.SetCardChainID(span, virgilCard); nil != err {
		return nil, api.ErrInternalError.WithMessage(
			"error trying set chainID for card(%s): %+v",
			virgilCard.GetID(), err,
		)
	}

	// check chain's deleted state in case card overwrite operation only:
	if "" != virgilCard.PreviousCardID {
		isDeleted, err := h.cardRepository.IsChainDeleted(
			span,
			request.UserID,
			request.ApplicationID,
			virgilCard.GetChainID(),
		)
		if nil != err {
			return nil, api.ErrInternalError.WithMessage(
				"error checking chain deleted state for chain(%s): %+v",
				virgilCard.GetChainID(), err,
			)
		}

		if isDeleted {
			return nil, tracer.SetSpanErrorAndReturn(span, api.ErrChainAlreadyDeleted)
		}
	}

	if err := h.cardSigner.SignCardByCardsService(span, virgilCard); nil != err {
		return nil, api.ErrInternalError.WithMessage(
			"sign card(%s) error: %+v",
			virgilCard.GetID(), err,
		)
	}

	if err := h.cardRepository.SaveCard(span, virgilCard); nil != err {
		return nil, api.ErrInternalError.WithMessage(
			"save card(%s) error: %+v",
			virgilCard.GetID(), err,
		)
	}

	return virgilCard, nil
}

//
// CardGet is a handler for GET /card/:card_id request.
//
func (h *Controller) CardGet(
	span tracer.Span,
	request *api.CardBaseRequest,
	cardID string,
) (*model.CardDTO, error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentController,
		},
	)
	defer span.Finish()

	card, err := h.cardRepository.GetCardByID(span, cardID)
	if nil != err {
		if err == cassandra.ErrEntityNotFound {
			return nil, api.ErrNotFound
		}

		return nil, api.ErrInternalError.WithMessage(
			"internal error for get card from the database by its ID(%s): %+v",
			cardID, err,
		)
	}

	if !card.DoesScopeMatch(request.ApplicationID) {
		return nil, tracer.SetSpanErrorAndReturn(
			span,
			api.ErrVirgilCardApplicationIDIsNotInTheAuthApplicationList,
		)
	}

	card.IsSuperseeded, err = h.cardRepository.DoesCardExistByPreviousIDAndScopeID(span, card.ID, card.ApplicationID)
	if err != nil {
		return nil, err
	}

	return card, nil
}

//
// CardSearch is a handler for POST /card/actions/search request.
//
func (h *Controller) CardSearch(span tracer.Span, request *api.CardSearchRequest) ([]*model.CardDTO, error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentController,
		},
	)
	defer span.Finish()

	if err := h.searchCardValidator.Validate(span, request); nil != err {
		return nil, err
	}

	virgilCards, err := h.cardRepository.SearchCardsByIdentities(
		span,
		request.GetIdentities(),
		request.ApplicationID,
	)
	if nil != err {
		return nil, api.ErrInternalError.WithMessage(
			"error searching cards by identities: %+v",
			err,
		)
	}

	return virgilCards, nil
}

//
// CardDelete is a handler for POST /card/actions/delete request.
//
func (h *Controller) CardDelete(span tracer.Span, request *api.CardDeleteRequest) (*model.CardDTO, error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentController,
		},
	)
	defer span.Finish()

	// TODO just for backward compatibility for JS SDK. Should be removed in future!.
	//  This check should be in transport/helper.go:63
	if request.UserID == "" {
		return nil, tracer.SetSpanErrorAndReturn(span, api.ErrIdentityHeaderNotSet)
	}

	virgilCard := model.NewCardDTO()
	if err := h.deleteCardValidator.Validate(span, request, virgilCard); nil != err {
		return nil, err
	}

	if err := h.cardRepository.SetCardChainID(span, virgilCard); nil != err {
		return nil, api.ErrInternalError.WithMessage(
			"error trying set chainID for card(%s): %+v",
			virgilCard.GetID(), err,
		)
	}

	if err := h.cardSigner.SignCardByCardsService(span, virgilCard); nil != err {
		return nil, api.ErrInternalError.WithMessage(
			"sign card(%s) error: %+v",
			virgilCard.GetID(), err,
		)
	}

	isDeletedNow, err := h.cardRepository.SetChainDeleted(
		span,
		request.UserID,
		request.ApplicationID,
		virgilCard.GetChainID(),
		time.Now().Unix(),
	)
	if nil != err {
		return nil, api.ErrInternalError.WithMessage(
			"error trying to mark chain(%s) as deleted: %+v",
			virgilCard.GetChainID(), err,
		)
	}

	if !isDeletedNow {
		return nil, tracer.SetSpanErrorAndReturn(span, api.ErrChainAlreadyDeleted)
	}

	if err := h.cardRepository.SaveCard(span, virgilCard); nil != err {
		return nil, api.ErrInternalError.WithMessage(
			"save card(%s) error: %+v",
			virgilCard.GetID(), err,
		)
	}

	return virgilCard, nil
}
