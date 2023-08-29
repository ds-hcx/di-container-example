package dao

import (
	"encoding/base64"

	"github.com/gocql/gocql"

	"github.com/VirgilSecurity/virgil-services-core-kit/db/cassandra"
	"github.com/VirgilSecurity/virgil-services-core-kit/errors"
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"
	"github.com/VirgilSecurity/virgil-services-core-kit/uuid"

	"github.com/VirgilSecurity/virgil-services-cards/src/model"
)

//
// SignatureList is a custom signature list type.
// It represents a list of hash maps that contain signature entries.
//
type SignatureList []map[string]string

//
// CardRepositoryProvider is an interface to operate over Virgil Card database entries.
//
type CardRepositoryProvider interface {
	//
	// GetCardByID returns a card by its ID.
	//
	GetCardByID(span tracer.Span, ID string) (*model.CardDTO, error)

	//
	// DoesCardExistByPreviousIDAndScopeID returns true if the card exists by search criteria.
	//
	DoesCardExistByPreviousIDAndScopeID(span tracer.Span, previousCardID, applicationID string) (bool, error)

	//
	// SaveCard saves the card to the database.
	//
	SaveCard(tracer.Span, *model.CardDTO) error

	//
	// SearchCardsByIdentities returns the cards matching search criteria.
	//
	SearchCardsByIdentities(span tracer.Span, identities []string, scopeID string) ([]*model.CardDTO, error)

	//
	// SetCardChainID sets chainID property of card given.
	//
	SetCardChainID(span tracer.Span, card *model.CardDTO) error

	//
	// IsChainDeleted checks if chain is deleted.
	//
	IsChainDeleted(span tracer.Span, identity string, scopeID string, chainID string) (bool, error)

	//
	// SetChainDeleted sets 'Deleted' time of chainID.
	// Returns 'false' in case the chain has already been deleted.
	//
	SetChainDeleted(span tracer.Span, identity string, scopeID string, chainID string, unixSeconds int64) (bool, error)
}

//
// ErrorCassandraNotFound to wrap out standard error
//
type ErrorCassandraNotFound error

//
// CardRepository id the data access layer to operate over Virgil Card DB instances.
//
type CardRepository struct {
	session *gocql.Session
}

//
// NewCardRepository returns an instance of the CardRepository.
//
func NewCardRepository(connector cassandra.GoCQLSessionProvider) *CardRepository {
	return &CardRepository{session: connector.GetGoCQLSession()}
}

//
// GetCardByID returns a Virgil Card by its ID
// The method returns an error if no entry was found in the database.
//
func (d *CardRepository) GetCardByID(span tracer.Span, ID string) (*model.CardDTO, error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentDAO,
		},
	)
	defer span.Finish()

	var (
		pubKeyString string
		signatures   SignatureList
	)

	var card = new(model.CardDTO)
	if err := d.session.Query(qGetCardByID, ID).Scan(
		&card.ID,
		&card.ContentSnapshot,
		&pubKeyString,
		&card.Identity,
		&card.ApplicationID,
		&card.ChainID,
		&signatures,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, cassandra.ErrEntityNotFound
		}

		return nil, tracer.SetSpanErrorAndReturn(
			span,
			errors.WithMessage(err, `internal error for get card by ID (%s)`, ID),
		)
	}

	publicKey, err := base64.StdEncoding.DecodeString(pubKeyString)
	if nil != err {
		return nil, tracer.SetSpanErrorAndReturn(
			span,
			errors.WithMessage(err, `virgil card public key (%s) decode error`, pubKeyString),
		)
	}

	card.PublicKey = publicKey
	card.Signatures = wrapDBSignatureListToDTOs(signatures)

	return card, nil
}

//
// DoesCardExistByPreviousIDAndScopeID returns true if card exists by search criteria.
//
func (d *CardRepository) DoesCardExistByPreviousIDAndScopeID(
	span tracer.Span,
	prevID, applicationID string,
) (bool, error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentDAO,
		},
	)
	defer span.Finish()

	var id string
	if err := d.session.Query(qIsPreviousCardIsExists, prevID, applicationID).Scan(&id); nil != err {
		if err == gocql.ErrNotFound {
			return false, nil
		}

		return false, tracer.SetSpanErrorAndReturn(span, err)
	}

	if id != "" {
		return true, nil
	}

	return false, nil
}

//
// SaveCard persists the Virgil Card DTO to the database.
//
func (d *CardRepository) SaveCard(span tracer.Span, card *model.CardDTO) error {

	// Create a root span, because action is complicated and contains several database queries below.
	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentDAO,
		},
	)
	defer span.Finish()

	batchSave := d.session.NewBatch(gocql.UnloggedBatch)
	batchSave.SetConsistency(gocql.LocalQuorum) //just in case

	batchSave.Query(qCreateCardInCardPKTable,
		card.GetID(),
		card.GetIdentity(),
		card.GetEncodedPublicKey(),
		card.GetContentSnapshot(),
		card.GetVersion(),
		card.GetApplicationID(),
		card.GetPreviousCardID(),
		wrapSignatureDTOsToDBSignatureList(card.Signatures),
		card.GetCreatedAt().Unix(),
		card.GetChainID(),
	)

	batchSave.Query(qCreateCardInIdentityPKTable,
		card.GetID(),
		card.GetIdentity(),
		card.GetEncodedPublicKey(),
		card.GetContentSnapshot(),
		card.GetVersion(),
		card.GetApplicationID(),
		card.GetPreviousCardID(),
		wrapSignatureDTOsToDBSignatureList(card.Signatures),
		card.GetCreatedAt().Unix(),
		card.GetChainID(),
	)

	// Insert/Update chain's IDs:
	if "" == card.PreviousCardID {
		batchSave.Query(getCreateChainIdentitiesQuery(card.GetID()),
			card.GetCreatedAt().Unix(),
			card.GetIdentity(),
			card.GetApplicationID(),
			card.GetChainID(),
		)
	} else {
		batchSave.Query(getUpdateChainIdentitiesQuery(card.GetID()),
			card.GetIdentity(),
			card.GetApplicationID(),
			card.GetChainID(),
		)
	}

	if "" != card.GetPreviousCardID() {
		batchSave.Query(qUpdatePreviousCardIDs, card.GetPreviousCardID(), card.GetApplicationID())
	}

	if err := d.session.ExecuteBatch(batchSave); err != nil {
		return tracer.SetSpanErrorAndReturn(span, errors.WithMessage(
			err,
			"unable to save a new card (%s)", card.ID,
		))
	}

	return nil
}

//
// SearchCardsByIdentities returns Virgil Cards by their identities. index - identity, application, card_id.
//
func (d *CardRepository) SearchCardsByIdentities(
	span tracer.Span,
	identities []string,
	scopeID string,
) ([]*model.CardDTO, error) {

	// Create a root span, because action is complicated and contains several database queries below.
	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentDAO,
		},
	)
	defer span.Finish()

	var (
		// chains search
		deletedAt    int
		chainCardIDs []string
		// cards content
		contentSnapshot string
		signatures      SignatureList

		cardIDs = make([]string, 0)
		cards   = make([]*model.CardDTO, 0)
	)

	// Mark query with separate Span.
	{
		span := span.Tracer().StartSpan(
			tracer.GetCallerInfo(),
			tracer.ChildOf(span.Context()),
			tracer.Tags{
				tracer.TagComponent: tracer.ComponentDAO,
			},
		)
		defer span.Finish()

		iterCardIDs := d.session.Query(getSearchCardIDsByMultipleIdentitiesQuery(identities), scopeID).Iter()
		for iterCardIDs.Scan(&deletedAt, &chainCardIDs) {
			if 0 >= deletedAt {
				cardIDs = append(cardIDs, chainCardIDs...)
			}
		}
		if err := iterCardIDs.Close(); nil != err {
			return nil, tracer.SetSpanErrorAndReturn(span, errors.WithMessage(
				err,
				"error selecting from chain table in SearchCardsByIdentities for identities (%v) and appID (%s)",
				identities,
				scopeID,
			))
		}
	}

	if 0 == len(cardIDs) {
		return []*model.CardDTO{}, nil
	}

	// Mark query with separate Span.
	{
		span := span.Tracer().StartSpan(
			tracer.GetCallerInfo(),
			tracer.ChildOf(span.Context()),
			tracer.Tags{
				tracer.TagComponent: tracer.ComponentDAO,
			},
		)
		defer span.Finish()

		cardsIterator := d.session.Query(getSearchByCardsIDs(cardIDs)).Iter()
		for cardsIterator.Scan(&contentSnapshot, &signatures) {
			cards = append(cards, &model.CardDTO{
				ContentSnapshot: contentSnapshot,
				Signatures:      wrapDBSignatureListToDTOs(signatures),
			})
		}
		if err := cardsIterator.Close(); nil != err {
			return nil, tracer.SetSpanErrorAndReturn(span, errors.WithMessage(
				err,
				"error selecting cards in SearchCardsByIdentities for cardIDs (%v)", cardIDs,
			))
		}
	}

	return cards, nil
}

//
// SetCardChainID sets chainID property of card given.
//
func (d *CardRepository) SetCardChainID(span tracer.Span, card *model.CardDTO) error {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentDAO,
		},
	)
	defer span.Finish()

	if "" == card.GetPreviousCardID() {
		// new card create
		chainID, err := uuid.NewV4()
		if nil != err {
			return tracer.SetSpanErrorAndReturn(
				span,
				errors.WithMessage(err, "error creating card's (%s) new chainID", card.GetID()),
			)

		}
		card.SetChainID(chainID.String())
		return nil
	}

	chainID := ""
	if err := d.session.Query(
		qGetChainIDByCardID,
		card.GetPreviousCardID(),
	).Scan(&chainID); nil != err {
		return tracer.SetSpanErrorAndReturn(span, errors.WithMessage(
			err,
			"error retrieving previous card's (%s) chainID", card.GetPreviousCardID(),
		))
	}

	if "" == chainID {
		return tracer.SetSpanErrorAndReturn(
			span,
			errors.New("previous card's (%s) chainID is empty", card.GetPreviousCardID()),
		)
	}

	card.SetChainID(chainID)

	return nil
}

//
// IsChainDeleted checks if chain is deleted.
//
func (d *CardRepository) IsChainDeleted(span tracer.Span, identity, scopeID, chainID string) (bool, error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentDAO,
		},
	)
	defer span.Finish()

	var deletedAt int64

	deletedAtIterator := d.session.Query(qGetChainDeletedAt, identity, scopeID, chainID).Iter()
	deletedAtIterator.Scan(&deletedAt)

	if err := deletedAtIterator.Close(); nil != err {
		return false, tracer.SetSpanErrorAndReturn(
			span,
			errors.WithMessage(err, "error retrieving chain's(%s) deleted_at", chainID),
		)
	}

	return 0 < deletedAt, nil
}

//
// SetChainDeleted sets 'Deleted' time of chainID.
// Returns 'false' in case the chain has already been deleted.
//
func (d *CardRepository) SetChainDeleted(
	span tracer.Span,
	identity, scopeID, chainID string,
	unixSeconds int64,
) (bool, error) {

	span = span.Tracer().StartSpan(
		tracer.GetCallerInfo(),
		tracer.ChildOf(span.Context()),
		tracer.Tags{
			tracer.TagComponent: tracer.ComponentDAO,
		},
	)
	defer span.Finish()

	var oldDeletedAt int64
	applied, err := d.session.Query(
		qSetChainDeletedAt,
		unixSeconds,
		identity,
		scopeID,
		chainID,
	).ScanCAS(&oldDeletedAt)
	if nil != err {
		return false, tracer.SetSpanErrorAndReturn(
			span,
			errors.WithMessage(err, "error updating chain's(%s) deleted_at", chainID),
		)
	}

	return applied, nil
}

//
// wrapDBSignaturesToDTOs wraps database signatures data into the DTOs.
//
func wrapDBSignatureListToDTOs(signatures SignatureList) []*model.CardSignatureDTO {

	var signs []*model.CardSignatureDTO
	for _, signature := range signatures {
		signs = append(signs, &model.CardSignatureDTO{
			Signer:    signature["signer"],
			Signature: signature["signature"],
			Snapshot:  signature["snapshot"],
		})
	}

	return signs
}

//
// wrapSignatureDTOsToDBSignatureList wraps signature DTOs to the database map records.
//
func wrapSignatureDTOsToDBSignatureList(signatureDTOs []*model.CardSignatureDTO) SignatureList {

	var signatures SignatureList
	for _, signature := range signatureDTOs {
		signatures = append(signatures, map[string]string{
			"signer":    signature.GetSigner(),
			"signature": signature.GetSignature(),
			"snapshot":  signature.GetExtraContent(),
		})
	}

	return signatures
}
