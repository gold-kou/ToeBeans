package controller

import (
	"log"
	"net/http"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"

	applicationLog "github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
)

func HealthController(w http.ResponseWriter, r *http.Request) {
	l, e := applicationLog.NewLogger()
	if e != nil {
		log.Panic(e)
	}
	l.LogHTTPAccess(r)

	switch r.Method {
	case http.MethodGet:
		helper.ResponseSimpleSuccess(w)
	default:
		methods := []string{http.MethodGet}
		helper.ResponseNotAllowedMethod(w, "not allowed method", methods)
	}
}
