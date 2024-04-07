package controller

import (
	"net/http"

	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/middleware"
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
	requestMiddlewares  []middleware.Request
	handler             Handler
	responseMiddlewares []middleware.Response
}

func (c *Controller) Serve(lc *context.LuxContext) error {
	if err := middleware.ApplyRequests(lc, c.requestMiddlewares); err != nil {
		return err
	}

	if err := c.handler(lc); err != nil {
		return err
	}

	if err := middleware.ApplyResponses(lc, c.responseMiddlewares); err != nil {
		return err
	}

	return nil
}
