package events

import (
	"strconv"

	"github.com/VirgilSecurity/virgil-services-core-kit/metrics"
)

//
// ServiceEvent describes business metric entity with information about current event.
//
type ServiceEvent struct {
	tag           metrics.EventTag
	actionID      int
	actionType    int
	accountID     string
	applicationID string
}

//
// NewServiceEvent is a constructor for Event.
//
func NewServiceEvent(actionID int, actionType int, accountID string, applicationID string) *ServiceEvent {
	return &ServiceEvent{
		tag:           metrics.ServiceTag,
		actionID:      actionID,
		actionType:    actionType,
		accountID:     accountID,
		applicationID: applicationID,
	}
}

//
// ID function generates unique identifier for Event based on Event data.
//
func (e *ServiceEvent) ID() uint64 {
	id := uint64(14695981039346656037)
	for _, c := range []byte(
		string(e.tag) + strconv.Itoa(e.actionID) + strconv.Itoa(e.actionType) + e.accountID + e.applicationID,
	) {
		id *= 1099511628211
		id ^= uint64(c)
	}

	return id
}

//
// Fields returns list available fields in struct.
//
func (e *ServiceEvent) Fields() metrics.MetricData {
	return metrics.MetricData{
		"event_tag":           e.tag,
		"event_action_id":     e.actionID,
		"event_action_type":   e.actionType,
		"user_account_id":     e.accountID,
		"user_application_id": e.applicationID,
	}
}
