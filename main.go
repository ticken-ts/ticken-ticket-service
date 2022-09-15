package main

import (
	"ticken-ticket-service/app"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/utils"
)

func main() {
	tickenConfig, err := utils.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	builder, err := infra.NewBuilder(tickenConfig)
	if err != nil {
		panic(err)
	}

	db := builder.BuildDb()
	router := builder.BuildRouter()

	tickenTicketServer := app.New(router, db, tickenConfig)
	if tickenConfig.IsDev() {
		tickenTicketServer.Populate()
	}

	tickenTicketServer.Start()
}
