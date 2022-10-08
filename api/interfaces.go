package api

import "ticken-ticket-service/infra"

type Controller interface {
	Setup(router infra.Router)
}

type Middleware interface {
	Setup(router infra.Router)
}
