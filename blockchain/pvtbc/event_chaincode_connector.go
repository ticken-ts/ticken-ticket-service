package pvtbc

import (
	"container/list"
	"encoding/json"
	"ticken-ticket-service/models"
	"time"
)

const (
	EventCCGetFunction   = "Get"
	TickenEventChaincode = "ticken-event"
)

type perBCEvent struct {
	EventID  string    `json:"event_id"`
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
	Sections list.List `json:"sections"`
}

type eventChaincodeConnector struct {
	baseCC ChaincodeConnector
}

func NewEventChaincodeConnector(hfc HyperledgerFabricConnector, channelName string) (EventChaincodeConnector, error) {
	cc, err := NewChaincodeConnector(hfc, channelName, TickenEventChaincode)
	if err != nil {
		return nil, err
	}

	eventCC := new(eventChaincodeConnector)
	eventCC.baseCC = cc

	return eventCC, nil
}

func (eventCC *eventChaincodeConnector) GetEvent(eventID string) (*models.Event, error) {
	eventData, err := eventCC.baseCC.Query(EventCCGetFunction, eventID)
	if err != nil {
		return nil, err
	}

	payload := new(perBCEvent)

	err = json.Unmarshal(eventData, &payload)
	if err != nil {
		return nil, err
	}

	event := models.Event{
		EventID: payload.EventID,
	}

	return &event, nil
}
