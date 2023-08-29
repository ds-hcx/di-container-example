package api

//
// CardSearchRequest is a search Virgil Card request object.
// it could contain:
// 1. "identity" string field which contains single identity e-mail.
// or
// 2. "identities" strings array to handle multi-identities search.
//
type CardSearchRequest struct {
	*Headers
	Identity   string   `json:"identity"`
	Identities []string `json:"identities"`
}

//
// GetIdentities returns the identities slice.
//
func (r *CardSearchRequest) GetIdentities() []string {
	var identities []string

	if r.Identity != "" {
		identities = append(identities, r.Identity)
	} else if len(r.Identities) > 0 {
		identities = append(identities, r.Identities...)
	}

	return identities
}
