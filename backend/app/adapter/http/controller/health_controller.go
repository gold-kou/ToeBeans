package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"

	applicationLog "github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
)

func HealthController(w http.ResponseWriter, r *http.Request) {
	l, e := applicationLog.NewLogger()
	if e != nil {
		log.Panic(e)
	}
	l.LogHTTPAccess(r)

	switch r.URL.Path {
	case "/health/liveness":
		switch r.Method {
		case http.MethodGet:

			helper.ResponseSimpleSuccess(w)
		default:
			methods := []string{http.MethodGet}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	case "/health/readiness":
		switch r.Method {
		case http.MethodGet:
			err := getHealthReadiness()
			if err != nil {
				helper.ResponseInternalServerError(w, fmt.Sprintf("readiness error: %s", err.Error()))
			}
			helper.ResponseSimpleSuccess(w)
		default:
			methods := []string{http.MethodGet}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	default:
		helper.ResponseInternalServerError(w, errMsgControllerPath)
	}
}

func getHealthReadiness() error {
	db, err := mysql.NewDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}
