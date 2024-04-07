package controller

import (
	"net"

	"github.com/snowmerak/lux/v3/context"
)

type SocketHandler func(ctx *context.WSContext) error

type SocketController struct {
	Handler SocketHandler
}

func (c *SocketController) Serve(conn net.Conn) error {
	ctx := &context.WSContext{Conn: conn}
	return c.Handler(ctx)
}
