package middleware

import (
	"net/http"
	"strings"

	"golang.org/x/net/context"

	httpContext "github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/lib"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context
		var cookie *http.Cookie
		var tokenUserName string
		var err error

		// ignore patterns
		ignoreReqs := map[string]string{"/csrf-token": http.MethodGet, "/user": http.MethodPost, "/login": http.MethodPost, "/user-activation/": http.MethodGet, "/password-reset-email": http.MethodPost, "/password-reset": http.MethodPost, "/health/liveness": http.MethodGet, "/health/readiness": http.MethodGet}
		for path, method := range ignoreReqs {
			// MEMO: /user-activation/{user_name}/{activation_key} を考慮してHasPrefixを使う
			if strings.HasPrefix(r.URL.Path, path) && r.Method == method {
				goto next
			}
		}

		// verify
		cookie, err = r.Cookie(helper.CookieIDToken)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tokenUserName, err = lib.VerifyToken(cookie.Value)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx = httpContext.SetTokenUserName(r.Context(), tokenUserName)

	next:
		if ctx == nil {
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
