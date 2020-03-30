package rest

import (
	"github.com/di-collective/ditebak/backend/pkg/queryables"
	"github.com/di-collective/ditebak/backend/pkg/service"
)

// Config of REST API
type Config struct {
	Resource   string
	Queryables queryables.Collection
	Service    service.Service

	CreatePayload func() interface{}            // constructor of HTTP request payload for CREATE
	UpdatePayload func() interface{}            // constructor of HTTP request payload for UPDATE
	Convert       func(interface{}) interface{} // convert HTTP request payload to service payload
}
