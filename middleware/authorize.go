package middleware

import (
	"net/http"

	"github.com/snowmerak/lux/v3/context"
)

type AuthChecker func(lc *context.LuxContext, authorizationHeader string, tokenCookies ...*http.Cookie) bool

type AuthMiddleware Request

func Auth(authChecker AuthChecker, tokenName ...string) AuthMiddleware {
	return func(ctx *context.LuxContext) (*context.LuxContext, int) {
		req := ctx.Request
		authorizationHeader := req.Header.Get("Authorization")
		cookies := []*http.Cookie(nil)
		for _, name := range tokenName {
			cookie, err := req.Cookie(name)
			if err == nil {
				cookies = append(cookies, cookie)
			}
		}
		if authChecker(ctx, authorizationHeader, cookies...) {
			return ctx, http.StatusOK
		}
		return ctx, http.StatusUnauthorized
	}
}
