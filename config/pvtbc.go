package config

type PvtbcConfig struct {
	MspID              string `mapstructure:"msp_id"`
	PeerEndpoint       string `mapstructure:"peer_endpoint"`
	GatewayPeer        string `mapstructure:"gateway_peer"`
	CertificatePath    string `mapstructure:"certificate_path"`
	PrivateKeyPath     string `mapstructure:"private_key_path"`
	TLSCertificatePath string `mapstructure:"tls_certificate_path"`
}
