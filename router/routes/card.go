package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	kitHTTP "github.com/VirgilSecurity/virgil-services-core-kit/http"
	"github.com/VirgilSecurity/virgil-services-core-kit/http/response"
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"

	"github.com/VirgilSecurity/virgil-services-cards/src/middleware"
	"github.com/VirgilSecurity/virgil-services-cards/src/transport"
)

const (

	//
	// CardIDPlaceholder represents the Card ID placeholder.
	//
	CardIDPlaceholder = "{card_id}"

	//
	// RoutePrefix base Cards service routing prefix.
	//
	RoutePrefix = "/card"

	//
	// RouteCardCreate POST /card route.
	//
	RouteCardCreate = RoutePrefix

	//
	// RouteCardGet GET /card/{card_id} route.
	//
	RouteCardGet = RoutePrefix + "/" + CardIDPlaceholder

	//
	// RouteCardSearch POST /card/actions/search route.
	//
	RouteCardSearch = RoutePrefix + "/actions/search"

	//
	// RouteCardDelete POST /card/actions/delete route.
	//
	RouteCardDelete = RoutePrefix + "/actions/delete"
)

//
// InitCardsRouteList makes an initialization of Cards routes.
//
func InitCardsRouteList(t tracer.Tracer, r kitHTTP.RouterProvider, h *transport.CardsHandler) {

	r.Post(RouteCardCreate, func(req *http.Request) response.Provider {
		return middleware.WithTracer(t, req, func(req *http.Request) response.Provider {
			return h.CardCreate(req)
		})
	})

	r.Get(RouteCardGet, func(req *http.Request) response.Provider {
		return middleware.WithTracer(t, req, func(req *http.Request) response.Provider {
			return h.CardGet(req, mux.Vars(req)["card_id"])
		})
	})

	r.Post(RouteCardSearch, func(req *http.Request) response.Provider {
		return middleware.WithTracer(t, req, func(req *http.Request) response.Provider {
			return h.CardSearch(req)
		})
	})

	r.Post(RouteCardDelete, func(req *http.Request) response.Provider {
		return middleware.WithTracer(t, req, func(req *http.Request) response.Provider {
			return h.CardDelete(req)
		})
	})
}
