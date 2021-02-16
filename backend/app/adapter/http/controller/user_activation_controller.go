package controller

import (
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	applicationLog "github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func UserActivationController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

	switch r.Method {
	case http.MethodPut:
		err = activateUser(r)
		switch err := err.(type) {
		case nil:
			helper.ResponseSimpleSuccess(w)
		case *helper.BadRequestError:
			helper.ResponseBadRequest(w, err.Error())
		case *helper.AuthorizationError:
			helper.ResponseUnauthorized(w, err.Error())
		case *helper.InternalServerError:
			helper.ResponseInternalServerError(w, err.Error())
		default:
			helper.ResponseInternalServerError(w, err.Error())
		}
	default:
		methods := []string{http.MethodPut}
		helper.ResponseNotAllowedMethod(w, "not allowed method", methods)
	}
}

func activateUser(r *http.Request) (err error) {
	// get path parameter
	vars := mux.Vars(r)
	userName, ok := vars["user_name"]
	if !ok || userName == "" {
		log.Println(err)
		return helper.NewBadRequestError("user_name cannot be blank")
	}
	activationKey, ok := vars["activation_key"]
	if !ok || activationKey == "" {
		log.Println(err)
		return helper.NewBadRequestError("activation_key cannot be blank")
	}

	// validation check
	if err = validation.Validate(userName, validation.Required, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	if err = validation.Validate(activationKey, validation.Required, is.UUID); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}

	// db connect
	db, err := mysql.NewDB()
	if err != nil {
		log.Println(err)
		return helper.NewInternalServerError(err.Error())
	}
	defer db.Close()
	tx := mysql.NewDBTransaction(db)

	// repository
	userRepo := repository.NewUserRepository(db)

	// UseCase
	u := usecase.NewUserActivation(r.Context(), tx, userName, activationKey, userRepo)
	if err = u.UserActivationUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrUserActivationNotFound {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return
}
