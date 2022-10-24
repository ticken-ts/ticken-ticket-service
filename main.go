package main

import (
	"ticken-ticket-service/app"
	"ticken-ticket-service/config"
	"ticken-ticket-service/env"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/log"
)

func main() {
	tickenEnv, err := env.Load()
	if err != nil {
		panic(err)
	}

	log.InitGlobalLogger()

	tickenConfig, err := config.Load(tickenEnv.ConfigFilePath, tickenEnv.ConfigFileName)
	if err != nil {
		log.TickenLogger.Log().Err(err)
		panic(err)
	}

	infraBuilder, err := infra.NewBuilder(tickenConfig)
	if err != nil {
		log.TickenLogger.Log().Err(err)
		panic(err)
	}

	tickenTicketServer := app.New(infraBuilder, tickenConfig)
	if tickenEnv.IsDev() {
		tickenTicketServer.Populate()
		tickenTicketServer.EmitFakeJWT()
	}

	tickenTicketServer.Start()
}
