package middleware

import (
	"net/http"

	"github.com/snowmerak/lux/v3/context"

	"github.com/snowmerak/lux/v3/util"
)

var AccessControl = accessControl{}

type accessControl struct{}

type AllowStaticIpsMiddleware Request

func (ac accessControl) AllowStaticIPs(ips ...string) AllowStaticIpsMiddleware {
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

type BlockStaticIPsMiddleware Request

func (ac accessControl) BlockStaticIPs(ips ...string) BlockStaticIPsMiddleware {
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

type AllowDynamicIPsMiddleware Request

func (ac accessControl) AllowDynamicIPs(checker func(remoteIP string) bool) AllowDynamicIPsMiddleware {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		remoteIP := util.GetIP(req.RemoteAddr)
		if checker(remoteIP) {
			return ctx, http.StatusOK
		}
		return ctx, http.StatusForbidden
	}
}

type BlockDynamicIPsMiddleware Request

func (ac accessControl) BlockDynamicIPs(checker func(remoteIP string) bool) BlockDynamicIPsMiddleware {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		remoteIP := util.GetIP(req.RemoteAddr)
		if checker(remoteIP) {
			return ctx, http.StatusForbidden
		}
		return ctx, http.StatusOK
	}
}

type AllowStaticPortsMiddleware Request

func (ac accessControl) AllowStaticPorts(ports ...string) AllowStaticPortsMiddleware {
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

type BlockStaticPortsMiddleware Request

func (ac accessControl) BlockStaticPorts(ports ...string) BlockStaticPortsMiddleware {
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

type AllowDynamicPortsMiddleware Request

func (ac accessControl) AllowDynamicPorts(checker func(remotePort string) bool) AllowDynamicPortsMiddleware {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		remotePort := util.GetPort(req.RemoteAddr)
		if checker(remotePort) {
			return ctx, http.StatusOK
		}
		return ctx, http.StatusForbidden
	}
}

type BlockDynamicPortsMiddleware Request

func (ac accessControl) BlockDynamicPorts(checker func(remotePort string) bool) BlockDynamicPortsMiddleware {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		remotePort := util.GetPort(req.RemoteAddr)
		if checker(remotePort) {
			return ctx, http.StatusForbidden
		}
		return ctx, http.StatusOK
	}
}
