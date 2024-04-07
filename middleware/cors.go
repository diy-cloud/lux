package middleware

import (
	"strings"

	"github.com/snowmerak/lux/context"
)

var SetAllow = setAllow{}

type setAllow struct{}

func (sa setAllow) Headers(headers ...string) Response {
	return func(l *context.LuxContext) (*context.LuxContext, error) {
		l.Response.Headers.Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
		return l, nil
	}
}

func (sa setAllow) Methods(methods ...string) Response {
	return func(l *context.LuxContext) (*context.LuxContext, error) {
		l.Response.Headers.Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
		return l, nil
	}
}

func (sa setAllow) Origins(origins ...string) Response {
	return func(l *context.LuxContext) (*context.LuxContext, error) {
		l.Response.Headers.Set("Access-Control-Allow-Origin", strings.Join(origins, ","))
		return l, nil
	}
}

func (sa setAllow) Credentials() Response {
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
