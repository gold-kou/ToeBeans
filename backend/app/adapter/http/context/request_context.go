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
const tokenUserNameContextKey contextKey = "token_user_name"

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

func SetTokenUserName(parent context.Context, tokenUserName string) context.Context {
	return context.WithValue(parent, tokenUserNameContextKey, tokenUserName)
}

func GetTokenUserName(ctx context.Context) (tokenUserName string, e error) {
	v := ctx.Value(tokenUserNameContextKey)
	tokenUserName, ok := v.(string)
	if !ok {
		e = errors.New("token_user_name is unset")
	}
	return
}
