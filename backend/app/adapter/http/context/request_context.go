package context

import (
	"context"
	"errors"
	"time"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
)

type contextKey string

const requestedAtContextKey contextKey = "requested_at"
const accessLogContextKey contextKey = "access_log"

func SetRequestedAt(parent context.Context, requestedAt time.Time) context.Context {
	return context.WithValue(parent, requestedAtContextKey, requestedAt)
}

func GetRequestedAt(ctx context.Context) (requestedAt time.Time, e error) {
	v := ctx.Value(requestedAtContextKey)
	requestedAt, ok := v.(time.Time)
	if !ok {
		e = errors.New("requested_at is unset")
	}
	return
}

func SetAccessLog(parent context.Context, accessLog *log.AccessLog) context.Context {
	return context.WithValue(parent, accessLogContextKey, accessLog)
}

func GetAccessLog(ctx context.Context) (accessLog *log.AccessLog, e error) {
	v := ctx.Value(accessLogContextKey)
	accessLog, ok := v.(*log.AccessLog)
	if !ok {
		e = errors.New("access_log is unset")
	}
	return
}
