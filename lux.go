package lux

import (
	ctx "context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"

	"github.com/caddyserver/certmagic"
	"github.com/julienschmidt/httprouter"
	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/swagger"
	"golang.org/x/net/http2"
)

type Lux struct {
	logger      *zerolog.Logger
	server      *http.Server
	builtRouter *httprouter.Router
	swagger     *swagger.Swagger
	jwtConfig   *context.JWTConfig
	ctx         ctx.Context
}

func New(swaggerInfo *swagger.Info, logger *zerolog.Logger) *Lux {
	swg := new(swagger.Swagger)
	if swaggerInfo != nil {
		swg.Info = *swaggerInfo
	}
	swg.SwaggerVersion = "2.0"
	return &Lux{
		logger:      logger,
		server:      new(http.Server),
		builtRouter: httprouter.New(),
		swagger:     swg,
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

func (l *Lux) SetInfoEmail(email string) {
	l.swagger.Info.Contact.Email = email
}

func (l *Lux) SetInfoLicense(name, link string) {
	l.swagger.Info.License.Name = name
	l.swagger.Info.License.URL = link
}

func (l *Lux) SetJWTConfig(cfg *context.JWTConfig) {
	l.jwtConfig = cfg
}

func (l *Lux) ShowSwagger(path string) {
	swaggerjson, err := json.Marshal(l.swagger)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(".", "swagger", "dist", "swagger.json"))
	if err != nil {
		panic(err)
	}

	f.Write(swaggerjson)
	f.Close()

	l.logger.Warn().Str("path", path).Msg("Swagger is available")
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

func (l *Lux) buildServer(ctx ctx.Context, addr string) {
	l.server.Addr = addr
	l.server.Handler = l
	l.builtRouter = new(httprouter.Router)
	l.logger.Info().Str("addr", addr).Msg("Server is ready to serve")
}

func (l *Lux) ListenAndServe1(ctx ctx.Context, addr string) error {
	l.buildServer(ctx, addr)
	if err := l.server.ListenAndServe(); err != nil {
		l.logger.Fatal().Str("error", err.Error()).Msg("Listen and serve error")
		return err
	}
	return nil
}

func (l *Lux) ListenAndServe1TLS(ctx ctx.Context, addr string, certFile string, keyFile string) error {
	l.buildServer(ctx, addr)
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
	if err := certmagic.HTTPS(addr, l.builtRouter); err != nil {
		l.logger.Fatal().Str("error", err.Error()).Msg("Listen and serve Auto TLS error")
		return err
	}
	return nil
}

func (l *Lux) ListenAndServe2(ctx ctx.Context, addr string) error {
	l.buildServer(ctx, addr)
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
