package transport

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	kitHTTP "github.com/VirgilSecurity/virgil-services-core-kit/http"

	"github.com/VirgilSecurity/virgil-services-cards/src/api"
)

//
// NewHeaders constructs Headers structure.
//
func NewHeaders(req *http.Request) (*api.Headers, error) {

	h := api.Headers{
		UserID:        req.Header.Get(kitHTTP.HeaderUserID),
		AccountID:     req.Header.Get(kitHTTP.HeaderAccountID),
		ApplicationID: req.Header.Get(kitHTTP.HeaderApplicationID),
	}

	if h.ApplicationID == "" {
		return nil, api.ErrApplicationIDHeaderIsNotSet
	}

	// if h.UserID == "" {
	// 	return nil, api.ErrIdentityHeaderNotSet
	// }

	return &h, nil
}

//
// NewGetRequest constructs CardBaseRequest structure.
//
func NewGetRequest(req *http.Request) (*api.CardBaseRequest, error) {

	h, err := NewHeaders(req)
	if err != nil {
		return nil, err
	}

	request := api.CardBaseRequest{
		Headers: h,
	}

	return &request, nil
}

//
// NewBaseRequest constructs CardBaseRequest structure.
//
func NewBaseRequest(req *http.Request) (*api.CardBaseRequest, error) {

	h, err := NewHeaders(req)
	if err != nil {
		return nil, err
	}

	request := api.CardBaseRequest{
		Headers: h,
	}

	if err = unmarshal(req.Body, &request); err != nil {
		return nil, err
	}

	return &request, nil
}

//
// NewCardSearchRequest constructs CardSearchRequest structure.
//
func NewCardSearchRequest(req *http.Request) (*api.CardSearchRequest, error) {

	h, err := NewHeaders(req)
	if err != nil {
		return nil, err
	}

	request := api.CardSearchRequest{
		Headers: h,
	}

	if err := unmarshal(req.Body, &request); err != nil {
		return nil, err
	}

	return &request, nil
}

//
// unmarshal makes unmarshal request body according request structure.
//
func unmarshal(req io.Reader, obj interface{}) error {

	body, err := ioutil.ReadAll(req)
	if err != nil {
		return api.ErrRequestParsing
	}

	if err = json.Unmarshal(body, obj); err != nil {
		// TODO(DF) separate different unmarshalling errors (add new api.Err) OR correct api.ErrRequestParsing error message
		// // JSON fields/types/structure errors will be validated later
		// _, ok := err.(*json.SyntaxError)
		// if ok {
		return api.ErrRequestParsing
		// }
	}

	return nil
}
