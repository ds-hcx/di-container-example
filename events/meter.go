package events

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/log"
	"github.com/VirgilSecurity/virgil-services-core-kit/metrics"
)

//
// EventProvider provides an interface to work with the service business metrics.
//
type EventProvider interface {
	//
	// IncCardCreateSuccess increments Card create success event.
	//
	IncCardCreateSuccess(accountID, applicationID string)

	//
	// IncCardOverrideSuccess increments Card override event.
	//
	IncCardOverrideSuccess(accountID, applicationID string)

	//
	// IncCardCreateError increments Card create error event.
	//
	IncCardCreateError(accountID, applicationID string)

	//
	// IncCardGetSuccess increments Card get success event.
	//
	IncCardGetSuccess(accountID, applicationID string)

	//
	// IncCardGetError increments Card get error event.
	//
	IncCardGetError(accountID, applicationID string)

	//
	// IncCardSearchSuccess increments Card search success event.
	//
	IncCardSearchSuccess(accountID, applicationID string)

	//
	// IncCardSearchError increments Card search error event.
	//
	IncCardSearchError(accountID, applicationID string)

	//
	// IncChainDeleteSuccess increments Card delete success event.
	//
	IncChainDeleteSuccess(accountID, applicationID string)

	//
	// IncChainDeleteError increments Card delete error event.
	//
	IncChainDeleteError(accountID, applicationID string)
}

//
// EventMeter struct represents Metrics consumer which will push them to the metrics client.
//
type EventMeter struct {
	meter  metrics.EventAdder
	logger log.Logger
}

//
// NewEventMeter initializes EventMeter client.
//
func NewEventMeter(meter metrics.EventAdder, logger log.Logger) *EventMeter {
	return &EventMeter{
		meter:  meter,
		logger: logger,
	}
}

//
// IncCardCreateSuccess increments Card create success event.
//
func (m EventMeter) IncCardCreateSuccess(accountID, applicationID string) {
	m.pushServiceEvent(metrics.CardCreateSuccess, accountID, applicationID)
}

//
// IncCardOverrideSuccess increments Card override event.
//
func (m EventMeter) IncCardOverrideSuccess(accountID, applicationID string) {
	m.pushServiceEvent(metrics.CardOverrideSuccess, accountID, applicationID)
}

//
// IncCardCreateError increments Card create error event.
//
func (m EventMeter) IncCardCreateError(accountID, applicationID string) {
	m.pushServiceEvent(metrics.CardCreateError, accountID, applicationID)
}

//
// IncCardGetSuccess increments Card get success event.
//
func (m EventMeter) IncCardGetSuccess(accountID, applicationID string) {
	m.pushServiceEvent(metrics.CardGetSuccess, accountID, applicationID)
}

//
// IncCardGetError increments Card get error event.
//
func (m EventMeter) IncCardGetError(accountID, applicationID string) {
	m.pushServiceEvent(metrics.CardGetError, accountID, applicationID)
}

//
// IncCardSearchSuccess increments Card search success event.
//
func (m EventMeter) IncCardSearchSuccess(accountID, applicationID string) {
	m.pushServiceEvent(metrics.CardSearchSuccess, accountID, applicationID)
}

//
// IncCardSearchError increments Card search error event.
//
func (m EventMeter) IncCardSearchError(accountID, applicationID string) {
	m.pushServiceEvent(metrics.CardSearchError, accountID, applicationID)
}

//
// IncChainDeleteSuccess increments Card delete success event.
//
func (m EventMeter) IncChainDeleteSuccess(accountID, applicationID string) {
	m.pushServiceEvent(metrics.ChainDeleteSuccess, accountID, applicationID)
}

//
// IncChainDeleteError increments Card delete error event.
//
func (m EventMeter) IncChainDeleteError(accountID, applicationID string) {
	m.pushServiceEvent(metrics.ChainDeleteError, accountID, applicationID)
}

//
// pushServiceEvent makes a push of service event to the old ES storage and to the Click House.
//
func (m *EventMeter) pushServiceEvent(actionID int, accountID string, applicationID string) {

	// TODO should be deprecated later.
	// Push service event to the old ES storage.
	if err := m.meter.Add(metrics.NewEvent(
		actionID,
		metrics.TypeAPI,
		accountID,
		applicationID,
	)); err != nil {
		m.logger.Error("%v", err)
	}

	// Push service event to the new ClickHouse storage.
	if err := m.meter.Add(NewServiceEvent(
		actionID,
		metrics.TypeAPI,
		accountID,
		applicationID,
	)); err != nil {
		m.logger.Error("%v", err)
	}
}
