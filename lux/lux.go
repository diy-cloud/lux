package lux

import (
	ctx "context"
	"errors"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/caddyserver/certmagic"
	"github.com/julienschmidt/httprouter"
	"github.com/snowmerak/lux/v3/context"
	"github.com/snowmerak/lux/v3/controller"
	"golang.org/x/net/http2"
)

type Lux struct {
	logger      *zerolog.Logger
	server      *http.Server
	builtRouter *httprouter.Router
	jwtConfig   *context.JWTConfig
	ctx         ctx.Context
}

func New() *Lux {
	return &Lux{
		logger:      &log.Logger,
		server:      new(http.Server),
		builtRouter: httprouter.New(),
	}
}

func SetLogger(l *Lux, logger *zerolog.Logger) {
	l.logger = logger
}

func SetReadHeaderTimeout(l *Lux, duration time.Duration) {
	l.server.ReadHeaderTimeout = duration
}

func SetReadTimeout(l *Lux, duration time.Duration) {
	l.server.ReadTimeout = duration
}

func SetWriteTimeout(l *Lux, duration time.Duration) {
	l.server.WriteTimeout = duration
}

func SetIdleTimeout(l *Lux, duration time.Duration) {
	l.server.IdleTimeout = duration
}

func SetMaxHeaderBytes(l *Lux, n int) {
	l.server.MaxHeaderBytes = n
}

func SetJWTConfig(l *Lux, cfg *context.JWTConfig) {
	l.jwtConfig = cfg
}

func (l *Lux) AddRestController(route string, method controller.Method, controller controller.RestController) {
	l.builtRouter.Handle(string(method), route, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		luxCtx := new(context.LuxContext)
		luxCtx.Request = r
		luxCtx.Response = context.NewResponse()
		luxCtx.RouteParams = p
		luxCtx.Context = l.ctx
		luxCtx.RequestContext = r.Context()
		luxCtx.Logger = l.logger

		if err := controller.Serve(luxCtx); err != nil {
			l.logger.Error().Str("error", err.Error()).Msg("Controller error")
		}

		for key, values := range luxCtx.Response.Headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(luxCtx.Response.StatusCode)
		w.Write(luxCtx.Response.Body)
	})
}

func (l *Lux) AddSocketController(route string, controller controller.SocketController) {
	l.builtRouter.Handle(http.MethodGet, route, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			l.logger.Error().Str("error", err.Error()).Msg("Socket upgrade error")
			return
		}

		if err := controller.Serve(conn); err != nil && errors.Is(err, wsutil.ClosedError{}) {
			l.logger.Error().Str("error", err.Error()).Msg("Socket controller error")
		}
	})
}

func (l *Lux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.builtRouter.ServeHTTP(w, r)
}

func (l *Lux) buildServer(_ ctx.Context, addr string) {
	l.server.Addr = addr
	l.server.Handler = l.builtRouter
	l.logger.Info().Str("addr", addr).Msg("Server is ready to serve")
}

func (l *Lux) ListenAndServe1(ctx ctx.Context, addr string) error {
	l.buildServer(ctx, addr)
	l.ctx = ctx
	if err := l.server.ListenAndServe(); err != nil {
		l.logger.Fatal().Str("error", err.Error()).Msg("Listen and serve error")
		return err
	}
	return nil
}

func (l *Lux) ListenAndServe1TLS(ctx ctx.Context, addr string, certFile string, keyFile string) error {
	l.buildServer(ctx, addr)
	l.ctx = ctx
	if err := l.server.ListenAndServeTLS(certFile, keyFile); err != nil {
		l.logger.Fatal().Str("error", err.Error()).Msg("Listen and serve TLS error")
		return err
	}
	return nil
}

func (l *Lux) ListenAndServe1AutoTLS(ctx ctx.Context, addr []string) error {
	if len(addr) == 0 {
		addr = []string{"localhost:443"}
	}
	l.buildServer(ctx, addr[0])
	l.ctx = ctx
	if err := certmagic.HTTPS(addr, l.builtRouter); err != nil {
		l.logger.Fatal().Str("error", err.Error()).Msg("Listen and serve Auto TLS error")
		return err
	}
	return nil
}

func (l *Lux) ListenAndServe2(ctx ctx.Context, addr string) error {
	l.buildServer(ctx, addr)
	l.ctx = ctx
	if err := http2.ConfigureServer(l.server, nil); err != nil {
		l.logger.Fatal().Str("error", err.Error()).Msg("Http2 configuration error")
		return err
	}
	if err := l.server.ListenAndServe(); err != nil {
		l.logger.Fatal().Str("error", err.Error()).Msg("Listen and serve http2 error")
		return err
	}
	return nil
}

func (l *Lux) ListenAndServe2TLS(ctx ctx.Context, addr string, certFile string, keyFile string) error {
	l.buildServer(ctx, addr)
	l.ctx = ctx
	if err := http2.ConfigureServer(l.server, nil); err != nil {
		l.logger.Fatal().Str("error", err.Error()).Msg("Http2 configuration error")
		return err
	}
	if err := l.server.ListenAndServeTLS(certFile, keyFile); err != nil {
		l.logger.Fatal().Str("error", err.Error()).Msg("Listen and serve http2 TLS error")
		return err
	}
	return nil
}

func (l *Lux) ListenAndServe2AutoTLS(ctx ctx.Context, addr []string) error {
	if len(addr) == 0 {
		addr = []string{"localhost:443"}
	}
	l.buildServer(ctx, addr[0])
	l.ctx = ctx
	if err := http2.ConfigureServer(l.server, nil); err != nil {
		l.logger.Fatal().Str("error", err.Error()).Msg("Http2 configuration error")
		return err
	}
	if err := certmagic.HTTPS(addr, l.builtRouter); err != nil {
		l.logger.Fatal().Str("error", err.Error()).Msg("Listen and serve http2 Auto TLS error")
		return err
	}
	return nil
}
