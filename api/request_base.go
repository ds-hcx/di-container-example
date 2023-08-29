package api

//
// Headers describes common values for most types of requests.
//
type Headers struct {
	UserID        string `json:"-"`
	ApplicationID string `json:"-"`
	AccountID     string `json:"-"`
}

//
// CardBaseRequest is a base structure for Virgil Card request which contains signed card object
// that comes to HTTP handler.
//
type CardBaseRequest struct {
	*Headers
	// CSR is a base64-encoded JSON message holding CSR-message (Card Signing Request).
	CSR string `json:"content_snapshot"`

	// CSRStamps is a collection of stamps objects that prove CSR validity.
	CSRStamps []CSRStamp `json:"signatures"`
}

//
// GetCSR returns content snapshot value.
//
func (r *CardBaseRequest) GetCSR() string {
	return r.CSR
}

//
// GetSelfStamp returns a self-signed CSR stamp.
//
func (r *CardBaseRequest) GetSelfStamp() (*CSRStamp, error) {
	for _, entry := range r.CSRStamps {
		if entry.IsSelf() {
			return &entry, nil
		}
	}

	return nil, ErrSelfCSRStampIsMissing
}
