package helpers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
)

type Signer struct {
	privateKey *rsa.PrivateKey
}

func NewSigner(privateKey string) (*Signer, error) {
	privKeyBytes := []byte(privateKey)
	key, err := ssh.ParseRawPrivateKey(privKeyBytes)
	if err != nil {
		return nil, err
	}

	signer := new(Signer)
	signer.privateKey = key.(*rsa.PrivateKey)
	return signer, nil
}

func (signer *Signer) Sign(data interface{}) ([]byte, error) {
	if signer.privateKey == nil {
		return nil, fmt.Errorf("private key is not set")
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	hashFunc := sha256.New()
	_, err = hashFunc.Write(bytes)
	if err != nil {
		panic(err)
	}
	msgHashSum := hashFunc.Sum(nil)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	signature, err := rsa.SignPSS(rand.Reader, signer.privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		panic(err)
	}

	return signature, nil
}
