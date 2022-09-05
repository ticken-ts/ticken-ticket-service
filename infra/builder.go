package infra

// TODO
// * Handle more than one service type using config file
// * Log errors. This include passing a logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ticken-ticket-service/infra/db"
	"ticken-ticket-service/utils"
)

type builder struct {
	tickenConfig *utils.TickenConfig
}

func NewBuilder(tickenConfig *utils.TickenConfig) (*builder, error) {
	if tickenConfig == nil {
		return nil, fmt.Errorf("configuration is mandatory")
	}

	builder := new(builder)
	builder.tickenConfig = tickenConfig

	return builder, nil
}

func (builder *builder) BuildDb() Db {
	switch builder.tickenConfig.Config.Database.Driver {
	case utils.MongoDriver:
		return buildMongoDb(builder.tickenConfig.Env.MongoUri)
	default:
		panic(fmt.Errorf("database driver %s not implemented", builder.tickenConfig.Config.Database.Driver))
	}
}

func (builder *builder) BuildRouter() Router {
	return gin.Default()
}

func buildMongoDb(uri string) Db {
	mongoDb := db.NewMongoDb()
	err := mongoDb.Connect(uri)
	if err != nil {
		panic(err)
	}
	return mongoDb
}
