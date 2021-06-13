package middleware

import (
	"net/http"
	"time"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"
)

func CurrentTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.SetRequestedAt(r.Context(), time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
