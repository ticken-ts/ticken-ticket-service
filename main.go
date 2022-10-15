package main

import (
	"ticken-ticket-service/app"
	"ticken-ticket-service/config"
	"ticken-ticket-service/env"
	"ticken-ticket-service/infra"
)

func main() {
	tickenEnv, err := env.Load()
	if err != nil {
		panic(err)
	}

	tickenConfig, err := config.Load(".")
	if err != nil {
		panic(err)
	}

	infraBuilder, err := infra.NewBuilder(tickenConfig)
	if err != nil {
		panic(err)
	}

	tickenTicketServer := app.New(infraBuilder, tickenConfig)
	if tickenEnv.IsDev() {
		tickenTicketServer.Populate()
		tickenTicketServer.EmitFakeJWT()
	}

	tickenTicketServer.Start()
}
