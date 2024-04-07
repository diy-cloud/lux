package middleware

import (
	"strings"

	"github.com/snowmerak/lux/v3/context"
)

var SetAllow = setAllow{}

type setAllow struct{}

type AllowHeadersMiddleware Response

func (sa setAllow) Headers(headers ...string) AllowHeadersMiddleware {
	return func(l *context.LuxContext) (*context.LuxContext, error) {
		l.Response.Headers.Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
		return l, nil
	}
}

type AllowMethodsMiddleware Response

func (sa setAllow) Methods(methods ...string) AllowMethodsMiddleware {
	return func(l *context.LuxContext) (*context.LuxContext, error) {
		l.Response.Headers.Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
		return l, nil
	}
}

type AllowOriginsMiddleware Response

func (sa setAllow) Origins(origins ...string) AllowOriginsMiddleware {
	return func(l *context.LuxContext) (*context.LuxContext, error) {
		l.Response.Headers.Set("Access-Control-Allow-Origin", strings.Join(origins, ","))
		return l, nil
	}
}

type AllowCredentialsMiddleware Response

func (sa setAllow) Credentials() AllowCredentialsMiddleware {
	return func(l *context.LuxContext) (*context.LuxContext, error) {
		l.Response.Headers.Set("Access-Control-Allow-Credentials", "true")
		return l, nil
	}
}

func (sa setAllow) CORS() Response {
	return func(l *context.LuxContext) (*context.LuxContext, error) {
		l.Response.Headers.Set("Access-Control-Allow-Headers", "*")
		l.Response.Headers.Set("Access-Control-Allow-Origin", "*")
		l.Response.Headers.Set("Access-Control-Allow-Methods", "*")
		l.Response.Headers.Set("Access-Control-Allow-Credentials", "true")
		return l, nil
	}
}
