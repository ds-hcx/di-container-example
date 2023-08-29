package dao

import (
	"fmt"
	"strings"
)

//
// Database constants.
//
const (
	CollectionCardWithIDPrimary       = "card_by_card_id"
	CollectionCardWithIdentityPrimary = "card_by_identity"
	CollectionCardPreviousIDs         = "card_previous_ids"
	CollectionCardChain               = "card_chain"

	InsertFormatFullCardInfo = `
	INSERT INTO %s (
		id,
		identity,
		public_key,
		content_snapshot,
		version,
		application_id,
		previous_card_id,
		signatures,
		created_at_timestamp,
		chain_id
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

//
// Database query statements.
//
// #nosec
var (
	// Select chain ID by card's ID.
	qGetChainIDByCardID = fmt.Sprintf(`
	SELECT
		chain_id
	FROM %s
	WHERE id = ? 
	`, CollectionCardWithIDPrimary)

	// Select chain ID by card's ID.
	qGetChainDeletedAt = fmt.Sprintf(`
	SELECT deleted_at 
	FROM %s
	WHERE identity = ? AND application_id = ? AND chain_id = ?
	`, CollectionCardChain)

	// Mark chain as deleted.
	qSetChainDeletedAt = fmt.Sprintf(`
	UPDATE %s 
		SET deleted_at = ? 
	WHERE identity = ? AND application_id = ? AND chain_id = ? 
		IF deleted_at = 0 
	`, CollectionCardChain)

	// Select card by its ID query.
	qGetCardByID = fmt.Sprintf(`
	SELECT
		id,
		content_snapshot,
		public_key,
		identity,
		application_id,
		chain_id,
		signatures
	FROM %s
	WHERE id = ?
	`, CollectionCardWithIDPrimary)

	// Get the count of cards by their previous ID and application ID query.
	qIsPreviousCardIsExists = fmt.Sprintf(`
	SELECT
		previous_card_id
	FROM %s
	WHERE previous_card_id = ? AND application_id = ? LIMIT 1
	`, CollectionCardPreviousIDs)

	qCreateCardInIdentityPKTable = fmt.Sprintf(`
	INSERT INTO %s (
		id,
		identity,
		public_key,
		content_snapshot,
		version,
		application_id,
		previous_card_id,
		signatures,
		created_at_timestamp,
		chain_id
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, CollectionCardWithIdentityPrimary)

	qCreateCardInCardPKTable = fmt.Sprintf(InsertFormatFullCardInfo, CollectionCardWithIDPrimary)

	// Update card previous IDs query.
	qUpdatePreviousCardIDs = fmt.Sprintf(`
	INSERT INTO
		%s
	(previous_card_id, application_id)
	VALUES (?, ?);
	`, CollectionCardPreviousIDs)
)

//
// getCreateChainIdentitiesQuery returns a create cards chain identities query.
//
func getCreateChainIdentitiesQuery(id string) string {

	return fmt.Sprintf(`
	UPDATE
		%s
	SET
		ids = ids + {'%s'},
		created_at_timestamp = ?,
		deleted_at = 0
	WHERE
		identity = ? AND application_id = ? AND chain_id = ?
	`, CollectionCardChain,
		strings.Replace(id, "'", "''", -1),
	)
}

//
// getUpdateChainIdentitiesQuery returns an update cards chain identities query.
//
func getUpdateChainIdentitiesQuery(id string) string {

	return fmt.Sprintf(`
	UPDATE
		%s
	SET
		ids = ids + {'%s'}
	WHERE
		identity = ? AND application_id = ? AND chain_id = ? `,
		CollectionCardChain,
		strings.Replace(id, "'", "''", -1))
}

//
// joinCqlEscapedStrings returns a result of joining strings with single quote character inside escaped
// and elements are separated by given separator.
//
func joinCqlEscapedStrings(literals []string, separator string) string {
	joined := ""
	for i, stringElement := range literals {
		if i > 0 {
			joined += separator
		}
		joined += strings.Replace(stringElement, "'", "''", -1)
	}
	return joined
}

//
// getSearchByCardsIDs returns a query to search cards of set of identities.
//
func getSearchByCardsIDs(ids []string) string {

	// Select cards for ID list
	qSelectCardsByIDs := fmt.Sprintf(`
	SELECT
		content_snapshot,
		signatures
	FROM
		%s
	WHERE`,
		CollectionCardWithIDPrimary)

	qSelectCardsByIDs += " id IN ('" + joinCqlEscapedStrings(ids, "', '") + "')"
	return qSelectCardsByIDs
}

//
// getSearchCardIDsByMultipleIdentitiesQuery returns a query to search cards of set of identities.
//
func getSearchCardIDsByMultipleIdentitiesQuery(identities []string) string {

	// Select NOT deleted chains by identities list
	qSelectCardIDsForSetOfIdentitiesAndAppID := fmt.Sprintf(`
	SELECT
		deleted_at, 
		ids
	FROM
		%s
	WHERE`,
		CollectionCardChain)

	qSelectCardIDsForSetOfIdentitiesAndAppID += " application_id = ?"
	qSelectCardIDsForSetOfIdentitiesAndAppID += " AND identity IN ('" + joinCqlEscapedStrings(identities, "', '") + "')"
	return qSelectCardIDsForSetOfIdentitiesAndAppID
}
