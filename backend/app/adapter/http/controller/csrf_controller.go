package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gorilla/csrf"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	modelHttp "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
)

func CSRFTokenController(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
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
	default:
		methods := []string{http.MethodGet}
		helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
	}
}
