package middleware

import (
	"net/http"

	"github.com/snowmerak/lux/context"

	"github.com/snowmerak/lux/util"
)

var AccessControl = accessControl{}

type accessControl struct{}

func (ac accessControl) AllowStaticIPs(ips ...string) Request {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		remoteIP := util.GetIP(req.RemoteAddr)
		for _, ip := range ips {
			if remoteIP == ip {
				return ctx, http.StatusOK
			}
		}
		return ctx, http.StatusForbidden
	}
}

func (ac accessControl) BlockStaticIPs(ips ...string) Request {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		remoteIP := util.GetIP(req.RemoteAddr)
		for _, ip := range ips {
			if remoteIP == ip {
				return ctx, http.StatusForbidden
			}
		}
		return ctx, http.StatusOK
	}
}

func (ac accessControl) AllowDynamicIPs(checker func(remoteIP string) bool) Request {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		remoteIP := util.GetIP(req.RemoteAddr)
		if checker(remoteIP) {
			return ctx, http.StatusOK
		}
		return ctx, http.StatusForbidden
	}
}

func (ac accessControl) BlockDynamicIPs(checker func(remoteIP string) bool) Request {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		remoteIP := util.GetIP(req.RemoteAddr)
		if checker(remoteIP) {
			return ctx, http.StatusForbidden
		}
		return ctx, http.StatusOK
	}
}

func (ac accessControl) AllowStaticPorts(ports ...string) Request {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		remotePort := util.GetPort(req.RemoteAddr)
		for _, port := range ports {
			if remotePort == port {
				return ctx, http.StatusOK
			}
		}
		return ctx, http.StatusForbidden
	}
}

func (ac accessControl) BlockStaticPorts(ports ...string) Request {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		remotePort := util.GetPort(req.RemoteAddr)
		for _, port := range ports {
			if remotePort == port {
				return ctx, http.StatusForbidden
			}
		}
		return ctx, http.StatusOK
	}
}

func (ac accessControl) AllowDynamicPorts(checker func(remotePort string) bool) Request {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		remotePort := util.GetPort(req.RemoteAddr)
		if checker(remotePort) {
			return ctx, http.StatusOK
		}
		return ctx, http.StatusForbidden
	}
}

func (ac accessControl) BlockDynamicPorts(checker func(remotePort string) bool) Request {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		remotePort := util.GetPort(req.RemoteAddr)
		if checker(remotePort) {
			return ctx, http.StatusForbidden
		}
		return ctx, http.StatusOK
	}
}
