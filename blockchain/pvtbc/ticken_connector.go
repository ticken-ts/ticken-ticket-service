package pvtbc

const (
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

type tickenConnector struct {
	EventChaincodeConnector
	TicketChaincodeConnector

	channel     string
	isConnected bool
	hfc         HyperledgerFabricConnector
}

func NewConnector() (TickenConnector, error) {
	tickenConnector := new(tickenConnector)

	hfc := NewHyperledgerFabricConnector(mspID, certPath, keyPath)
	err := hfc.Connect()
	if err != nil {
		return nil, err
	}

	tickenConnector.hfc = hfc

	return tickenConnector, nil
}

func (tickenConnector *tickenConnector) Connect(channel string) error {
	tickenConnector.isConnected = false

	eventCCConnector, err := NewEventChaincodeConnector(tickenConnector.hfc, channel)
	if err != nil {
		return err
	}

	ticketCCConnector, err := NewTicketChaincodeConnector(tickenConnector.hfc, channel)
	if err != nil {
		return err
	}

	tickenConnector.EventChaincodeConnector = eventCCConnector
	tickenConnector.TicketChaincodeConnector = ticketCCConnector
	tickenConnector.isConnected = true

	return nil
}

func (tickenConnector *tickenConnector) IsConnected() bool {
	return tickenConnector.isConnected
}
