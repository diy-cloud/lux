package middleware

import (
	"fmt"
	"net/http"

	"github.com/snowmerak/lux/context"
)

type Request func(*context.LuxContext) (*context.LuxContext, int)
type Response func(*context.LuxContext) (*context.LuxContext, error)

func ApplyRequests(ctx *context.LuxContext, middlewares []Request) error {
	for _, m := range middlewares {
		if m == nil {
			continue
		}

		_, code := m(ctx)
		if 400 <= code && code < 600 {
			ctx.Response.WriteHeader(code)
			return fmt.Errorf("middleware request reading %s: %s from %s", ctx.Request.URL.Path, http.StatusText(code), ctx.Request.RemoteAddr)
		}
	}
	return nil
}

func ApplyResponses(ctx *context.LuxContext, middlewares []Response) error {
	for _, m := range middlewares {
		if m == nil {
			continue
		}

		_, err := m(ctx)
		if err != nil {
			return fmt.Errorf("middleware response writing %s: %s from %s", ctx.Request.URL.Path, err.Error(), ctx.Request.RemoteAddr)
		}
	}
	return nil
}
