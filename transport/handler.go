package transport

import (
	"net/http"

	kitHTTP "github.com/VirgilSecurity/virgil-services-core-kit/http"
	"github.com/VirgilSecurity/virgil-services-core-kit/http/response"
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"

	"github.com/VirgilSecurity/virgil-services-cards/src/app/controller"
	"github.com/VirgilSecurity/virgil-services-cards/src/events"
)

const (
	// nolint
	SuperseededCardIDHTTPHeader = "X-Virgil-Is-Superseeded"
)

//
// CardsHandler provides an abstraction on transport layer.
//
type CardsHandler struct {
	eventMeter      events.EventProvider
	cardsController controller.Provider
}

//
// NewCardsHandler return Cards handler instance.
//
func NewCardsHandler(keysController controller.Provider, eventMeter events.EventProvider) *CardsHandler {

	return &CardsHandler{
		eventMeter:      eventMeter,
		cardsController: keysController,
	}
}

//
// CardCreate handles POST /card endpoint.
//
func (h *CardsHandler) CardCreate(req *http.Request) response.Provider {

	span := tracer.SpanFromContext(req.Context())
	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentTransport,
		},
	)
	defer span.Finish()

	request, err := NewBaseRequest(req)
	if err != nil {
		return response.New(tracer.SetSpanErrorAndReturn(span, err))
	}

	card, err := h.cardsController.CardCreate(span, request)
	if err != nil {
		h.eventMeter.IncCardCreateError(request.AccountID, request.ApplicationID)
		return response.New(err)
	}

	if card.PreviousCardID == "" {
		h.eventMeter.IncCardCreateSuccess(request.AccountID, request.ApplicationID)
	} else {
		h.eventMeter.IncCardOverrideSuccess(request.AccountID, request.ApplicationID)
	}

	return response.New(card).SetStatus(kitHTTP.StatusCreated)
}

//
// CardGet handles GET /card/:card_id endpoint.
//
func (h *CardsHandler) CardGet(req *http.Request, cardID string) response.Provider {

	span := tracer.SpanFromContext(req.Context())
	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentTransport,
		},
	)
	defer span.Finish()

	request, err := NewGetRequest(req)
	if err != nil {
		return response.New(tracer.SetSpanErrorAndReturn(span, err))
	}

	card, err := h.cardsController.CardGet(span, request, cardID)
	if err != nil {
		h.eventMeter.IncCardGetError(request.AccountID, request.ApplicationID)
		return response.New(err)
	}
	h.eventMeter.IncCardGetSuccess(request.AccountID, request.ApplicationID)

	resp := response.New(card)
	if card.IsSuperseeded {
		resp.SetHeader(SuperseededCardIDHTTPHeader, "true")
	}

	return resp
}

//
// CardSearch handles POST /card/actions/search endpoint.
//
func (h *CardsHandler) CardSearch(req *http.Request) response.Provider {

	span := tracer.SpanFromContext(req.Context())
	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentTransport,
		},
	)
	defer span.Finish()

	request, err := NewCardSearchRequest(req)
	if err != nil {
		return response.New(tracer.SetSpanErrorAndReturn(span, err))
	}

	cards, err := h.cardsController.CardSearch(span, request)
	if err != nil {
		h.eventMeter.IncCardSearchError(request.AccountID, request.ApplicationID)
		return response.New(err)
	}
	h.eventMeter.IncCardSearchSuccess(request.AccountID, request.ApplicationID)

	return response.New(cards)
}

//
// CardDelete handles POST /card/actions/delete endpoint.
//
func (h *CardsHandler) CardDelete(req *http.Request) response.Provider {

	span := tracer.SpanFromContext(req.Context())
	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentTransport,
		},
	)
	defer span.Finish()

	request, err := NewBaseRequest(req)
	if err != nil {
		return response.New(tracer.SetSpanErrorAndReturn(span, err))
	}

	card, err := h.cardsController.CardDelete(span, request)
	if err != nil {
		h.eventMeter.IncChainDeleteError(request.AccountID, request.ApplicationID)
		return response.New(err)
	}
	h.eventMeter.IncChainDeleteSuccess(request.AccountID, request.ApplicationID)

	return response.New(card)
}
