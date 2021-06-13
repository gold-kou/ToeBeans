package middleware

import (
	"net/http"
	"time"

	requestContext "github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
)

type StatusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *StatusResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.status = statusCode
}

type LoggingMiddleware struct {
	logger *log.Logger
}

func NewLoggingMiddleware(logger *log.Logger) LoggingMiddleware {
	return LoggingMiddleware{logger: logger}
}

func (mw LoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessLog := &log.AccessLog{
			Method:        r.Method,
			Host:          r.Host,
			Path:          r.URL.Path,
			Query:         r.URL.RawQuery,
			RequestSize:   r.ContentLength,
			UserAgent:     r.UserAgent(),
			RemoteAddr:    r.RemoteAddr,
			XForwardedFor: r.Header.Get(helper.HeaderKeyXForwardedFor),
			Referer:       r.Referer(),
			Protocol:      r.Proto,
		}
		ctx := requestContext.SetAccessLog(r.Context(), accessLog)
		sw := &StatusResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(sw, r.WithContext(ctx))

		accessLog.Status = sw.status
		requestedAt, e := requestContext.GetRequestedAt(ctx)
		if e != nil {
			panic(e)
		} else {
			accessLog.Latency = time.Since(requestedAt)
		}
		mw.logger.LogHTTPAccess(accessLog)
	})
}
