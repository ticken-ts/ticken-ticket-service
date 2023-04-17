package fakes

import (
	"fmt"
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

		_, err := loader.serviceProvider.GetUserManager().RegisterUser(
			attendant.Email,
			password,
			attendant.Firstname,
			attendant.Lastname,
			attendant.WalletPrivKey,
		)

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
