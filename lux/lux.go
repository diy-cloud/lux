package lux

import (
	ctx "context"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/caddyserver/certmagic"
	"github.com/julienschmidt/httprouter"
	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/controller"
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

func (l *Lux) SetLogger(logger *zerolog.Logger) {
	l.logger = logger
}

func (l *Lux) SetReadHeaderTimeout(duration time.Duration) {
	l.server.ReadHeaderTimeout = duration
}

func (l *Lux) SetReadTimeout(duration time.Duration) {
	l.server.ReadTimeout = duration
}

func (l *Lux) SetWriteTimeout(duration time.Duration) {
	l.server.WriteTimeout = duration
}

func (l *Lux) SetIdleTimeout(duration time.Duration) {
	l.server.IdleTimeout = duration
}

func (l *Lux) SetMaxHeaderBytes(n int) {
	l.server.MaxHeaderBytes = n
}

func (l *Lux) SetJWTConfig(cfg *context.JWTConfig) {
	l.jwtConfig = cfg
}

func (l *Lux) AddController(route string, method controller.Method, controller controller.Controller) {
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

func (l *Lux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	luxCtx := new(context.LuxContext)
	luxCtx.Request = r
	luxCtx.Response = context.NewResponse()
	luxCtx.JWTConfig = l.jwtConfig
	luxCtx.Logger = l.logger
	luxCtx.Context = l.ctx
	luxCtx.RequestContext = r.Context()
	defer func() {
		for key, values := range luxCtx.Response.Headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(luxCtx.Response.StatusCode)
		w.Write(luxCtx.Response.Body)
	}()
	l.builtRouter.ServeHTTP(luxCtx.Response, luxCtx.Request)
}

func (l *Lux) buildServer(_ ctx.Context, addr string) {
	l.server.Addr = addr
	l.server.Handler = l
	l.builtRouter = new(httprouter.Router)
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
