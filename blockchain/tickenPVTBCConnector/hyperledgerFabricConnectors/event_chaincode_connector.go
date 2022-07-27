package hyperledgerFabricConnectors

import (
	"container/list"
	"encoding/json"
	"google.golang.org/grpc"
	"ticken-ticket-service/models/event"
	"time"
)

const (
	EVENT_CC_GET_FUNCTION = "Get"
)

type perBCEvent struct {
	EventID  string    `json:"event_id"`
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
	Sections list.List `json:"sections"`
}

type EventChaincodeConnector interface {
	Connect(grpcConn *grpc.ClientConn, channel string) error
}

type eventChaincodeConnector struct {
	hyperledgerFabricBaseConnector BaseConnector
}

func NewEventChaincodeConnector() EventChaincodeConnector {
	tts := new(eventChaincodeConnector)
	return tts
}

func (c *eventChaincodeConnector) Connect(grpcConn *grpc.ClientConn, channel string) error {
	c.hyperledgerFabricBaseConnector = NewBaseConnector(mspID, certPath, keyPath)
	err := c.hyperledgerFabricBaseConnector.Connect(grpcConn, channel, chaincode)
	if err != nil {
		return err
	}
	return nil
}

func (c *eventChaincodeConnector) GetEvent(eventID string) (*event.Event, error) {
	eventData, err := c.hyperledgerFabricBaseConnector.Query(EVENT_CC_GET_FUNCTION, eventID)
	if err != nil {
		return nil, err
	}

	payload := new(perBCEvent)

	err = json.Unmarshal(eventData, &payload)
	if err != nil {
		return nil, err
	}

	event := event.Event{
		EventID: payload.EventID,
	}

	return &event, nil
}
