package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gorilla/csrf"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	applicationLog "github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
	modelHttp "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
)

func CSRFTokenController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

	switch r.Method {
	case http.MethodGet:
		switch err := err.(type) {
		case nil:
			w.WriteHeader(http.StatusOK)
			w.Header().Set(helper.HeaderKeyContentType, helper.HeaderValueApplicationJSON)

			token := csrf.Token(r)
			csrf.MaxAge(3600 * lib.TokenExpirationHour)
			resp := modelHttp.ResponseGetCsrfToken{
				CsrfToken: token,
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				panic(err.Error())
			}
		case *helper.InternalServerError:
			helper.ResponseInternalServerError(w, err.Error())
		default:
			helper.ResponseInternalServerError(w, err.Error())
		}
	default:
		methods := []string{http.MethodGet}
		helper.ResponseNotAllowedMethod(w, "not allowed method", methods)
	}
}
