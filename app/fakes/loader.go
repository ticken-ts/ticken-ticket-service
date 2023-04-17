package fakes

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"os"
	"reflect"
	"strings"
	"ticken-ticket-service/config"
	"ticken-ticket-service/env"
	"ticken-ticket-service/log"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/services"
	"ticken-ticket-service/utils"
)

const Filename = "fakes.json"

type Loader struct {
	repoProvider    repos.IProvider
	serviceProvider services.IProvider
	config          *config.Config
}

func NewFakeLoader(repoProvider repos.IProvider, serviceProvider services.IProvider, config *config.Config) *Loader {
	return &Loader{
		config:          config,
		repoProvider:    repoProvider,
		serviceProvider: serviceProvider,
	}
}

func (loader *Loader) Populate() error {
	if env.TickenEnv.IsProd() || !utils.FileExists(Filename) {
		return nil
	}

	seedContent := make(map[string]json.RawMessage)

	seedRawContent, err := os.ReadFile(Filename)
	if err != nil {
		log.TickenLogger.Panic().Msg(fmt.Sprintf("failed to read seed file: %s", err.Error()))
	}

	if err := json.Unmarshal(seedRawContent, &seedContent); err != nil {
		log.TickenLogger.Panic().Msg(fmt.Sprintf("failed to unmarshal seed file: %s", err.Error()))
	}

	for _, modelName := range []string{"attendant"} {
		log.TickenLogger.Info().Msg(
			fmt.Sprintf("%s: %s",
				color.GreenString("seeding model: "),
				color.New(color.FgBlue, color.Bold).Sprintf(modelName)),
		)

		switch modelName {

		case strings.ToLower(reflect.TypeOf(models.Attendant{}).Name()):
			attendantsToSeed := make([]*SeedAttendant, 0)

			if err := json.Unmarshal(seedContent[modelName], &attendantsToSeed); err != nil {
				log.TickenLogger.Error().Msg(fmt.Sprintf("failed to unmarshal attendant values: %s", err.Error()))
				continue
			}

			seedErrors := loader.seedAttendant(attendantsToSeed)
			if seedErrors != nil && len(seedErrors) > 0 {
				for _, seedError := range seedErrors {
					log.TickenLogger.Error().Msg(seedError.Error())
				}
				continue
			}
		}
	}
	return nil
}
