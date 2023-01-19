package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func loadPrivateKey(RSAPrivateKey string) (*rsa.PrivateKey, error) {
	privPem, _ := pem.Decode([]byte(RSAPrivateKey))
	if privPem.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("prviate key is not from type RSA")
	}

	var privateParsedKye interface{}
	privateParsedKye, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse RSA public key, generating a temp one: %s", err.Error())
	}

	privateKey, ok := privateParsedKye.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("unable to parse RSA private key")
	}

	return privateKey, nil
}

func loadPublicKey(RSAPublicKey string) (*rsa.PublicKey, error) {
	pubPem, _ := pem.Decode([]byte(RSAPublicKey))
	if pubPem.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("public key is not from type RSA")
	}

	var publicParsedKey interface{}
	publicParsedKey, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse RSA public key, generating a temp one: %s", err.Error())
	}

	publicKey, ok := publicParsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unable to parse RSA private key")
	}

	return publicKey, nil
}

func LoadRSA(RSAPrivateKey string, RSAPublicKey string) (*rsa.PrivateKey, error) {
	priv, err := loadPrivateKey(RSAPrivateKey)
	if err != nil {
		return nil, err
	}

	pub, err := loadPublicKey(RSAPublicKey)
	if err != nil {
		return nil, err
	}

	priv.PublicKey = *pub
	return priv, nil
}
