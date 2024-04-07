package controller

import (
	"net/http"

	"github.com/snowmerak/lux/v3/context"
	"github.com/snowmerak/lux/v3/middleware"
)

type Method string

const (
	GET     = Method(http.MethodGet)
	POST    = Method(http.MethodPost)
	PUT     = Method(http.MethodPut)
	PATCH   = Method(http.MethodPatch)
	DELETE  = Method(http.MethodDelete)
	OPTIONS = Method(http.MethodOptions)
	HEAD    = Method(http.MethodHead)
	CONNECT = Method(http.MethodConnect)
	TRACE   = Method(http.MethodTrace)
	ANY     = Method("*")
)

type Handler func(ctx *context.LuxContext) error

type Controller struct {
	RequestMiddlewares  []middleware.Request
	Handler             Handler
	ResponseMiddlewares []middleware.Response
}

func (c *Controller) Serve(lc *context.LuxContext) error {
	if err := middleware.ApplyRequests(lc, c.RequestMiddlewares); err != nil {
		return err
	}

	if err := c.Handler(lc); err != nil {
		return err
	}

	if err := middleware.ApplyResponses(lc, c.ResponseMiddlewares); err != nil {
		return err
	}

	return nil
}
