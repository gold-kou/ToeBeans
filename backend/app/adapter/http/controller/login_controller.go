package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"

	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"

	model "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
)

func LoginController(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		idToken, err := login(r)
		switch err := err.(type) {
		case nil:
			expiration := time.Now().Add(lib.TokenExpirationHour * time.Hour)
			cookie := &http.Cookie{
				Name:     helper.CookieIDToken,
				Value:    idToken,
				Expires:  expiration,
				HttpOnly: true,
				Secure:   false,
			}
			http.SetCookie(w, cookie)

			w.Header().Set(helper.HeaderKeyContentType, helper.HeaderValueApplicationJSON)
			w.WriteHeader(http.StatusOK)

			resp := model.Token{
				IdToken: idToken,
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				panic(err.Error())
			}
		case *helper.BadRequestError:
			helper.ResponseBadRequest(w, err.Error())
		case *helper.AuthorizationError:
			helper.ResponseUnauthorized(w, err.Error())
		case *helper.ForbiddenError:
			helper.ResponseForbidden(w, err.Error())
		case *helper.InternalServerError:
			helper.ResponseInternalServerError(w, err.Error())
		default:
			helper.ResponseInternalServerError(w, err.Error())
		}
	default:
		methods := []string{http.MethodPost}
		helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
	}
}

func login(r *http.Request) (idToken string, err error) {
	// get request parameter
	var reqLogin *model.RequestLogin
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return "", helper.NewBadRequestError(err.Error())
	}
	defer r.Body.Close()
	if err = json.Unmarshal(b, &reqLogin); err != nil {
		log.Println(err)
		return "", helper.NewBadRequestError(err.Error())
	}

	// validation check
	err = reqLogin.ValidateParam()
	if err != nil {
		log.Println(err)
		return "", helper.NewBadRequestError(err.Error())
	}

	// db connect
	db, err := mysql.NewDB()
	if err != nil {
		log.Println(err)
		return "", helper.NewInternalServerError(err.Error())
	}
	defer db.Close()
	tx := mysql.NewDBTransaction(db)

	// repository
	userRepo := repository.NewUserRepository(db)

	// UseCase
	l := usecase.NewLogin(r.Context(), tx, reqLogin, userRepo)
	if idToken, err = l.LoginUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData || err == usecase.ErrNotCorrectPassword {
			return "", helper.NewBadRequestError(errMsgWrongUserNameOrPassword)
		}
		if err == usecase.ErrNotVerifiedUser {
			return "", helper.NewForbiddenError(err.Error())
		}
		return "", helper.NewInternalServerError(err.Error())
	}
	return
}
