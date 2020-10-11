package controller

import (
	"log"
	"net/http"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"

	applicationLog "github.com/gold-kou/ToeBeans/app/adapter/http/log"
)

func HealthController(w http.ResponseWriter, r *http.Request) {
	l, e := applicationLog.NewLogger()
	if e != nil {
		log.Panic(e)
	}
	l.LogHTTPAccess(r)

	switch r.Method {
	case "GET":
		helper.ResponseSimpleSuccess(w)
	default:
		helper.ResponseBadRequest(w, "not allowed method")
	}
}
