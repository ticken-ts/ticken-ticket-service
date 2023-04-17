package fakes

import (
	"encoding/pem"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

type SeedAttendant struct {
	Email         string `json:"email"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	WalletPrivKey string `json:"wallet_private_key"`
}

func (loader *Loader) seedAttendant(toSeed []*SeedAttendant) []error {
	var seedErrors = make([]error, 0)

	if loader.repoProvider.GetUserRepository().Count() > 0 {
		return seedErrors
	}

	for _, attendant := range toSeed {
		// lower firstname + dot (".") + lastname
		// example: Facundo Torraca -> facundo.torraca
		password := strings.ToLower(strings.Join([]string{attendant.Firstname, attendant.Lastname}, "."))

		pemEncodedPrivKey, err := walletPrivKeyFromHEXToPEM(attendant.WalletPrivKey)
		if err == nil {
			_, err = loader.serviceProvider.GetUserManager().RegisterUser(
				attendant.Email,
				password,
				attendant.Firstname,
				attendant.Lastname,
				pemEncodedPrivKey,
			)
		}

		if err != nil {
			seedErrors = append(
				seedErrors,
				fmt.Errorf("failed to seed organizer with name %s %s: %s",
					attendant.Firstname, attendant.Lastname, err.Error()),
			)
		}
	}

	return seedErrors
}

// walletPrivKeyFromHEXToPEM passes the wallet private key
// from hex format to PEM format. Hex format is the format
// that is displayed by ganache, and PEM format is the
// formar adopted by Ticken by convention
func walletPrivKeyFromHEXToPEM(hexPrivKey string) (string, error) {
	privKey, err := crypto.HexToECDSA(hexPrivKey)
	if err != nil {
		return "", err
	}

	pemEncoded := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: crypto.FromECDSA(privKey),
		})

	return string(pemEncoded), nil
}
