package middleware

import (
	"fmt"
	"net/http"

	"github.com/snowmerak/lux/context"
)

type Request func(*context.LuxContext) (*context.LuxContext, int)
type Response func(*context.LuxContext) (*context.LuxContext, error)

func ApplyRequests(ctx *context.LuxContext, middlewares []Request) string {
	for _, m := range middlewares {
		if m == nil {
			continue
		}

		_, code := m(ctx)
		if 400 <= code && code < 600 {
			ctx.Response.WriteHeader(code)
			return fmt.Sprintf("Middleware Request Reading %s: %s from %s", ctx.Request.URL.Path, http.StatusText(code), ctx.Request.RemoteAddr)
		}
	}
	return ""
}

func ApplyResponses(ctx *context.LuxContext, middlewares []Response) string {
	for _, m := range middlewares {
		if m == nil {
			continue
		}

		_, err := m(ctx)
		if err != nil {
			return fmt.Sprintf("Middleware Response Writing %s: %s from %s", ctx.Request.URL.Path, err.Error(), ctx.Request.RemoteAddr)
		}
	}
	return ""
}
