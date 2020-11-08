package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/app/lib"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"
	applicationLog "github.com/gold-kou/ToeBeans/app/adapter/http/log"
	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/application/usecase"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

func FollowUserNameController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

	switch r.Method {
	case http.MethodDelete:
		err = deleteFollow(r)
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
		methods := []string{http.MethodDelete}
		helper.ResponseNotAllowedMethod(w, "not allowed method", methods)
	}
}

func deleteFollow(r *http.Request) error {
	// authorization
	userName, err := lib.VerifyHeaderToken(r)
	if err != nil {
		log.Println(err)
		return helper.NewAuthorizationError(err.Error())
	}

	// get request parameter
	vars := mux.Vars(r)
	followedUserName, ok := vars["followed_user_name"]
	if !ok {
		log.Println(err)
		return helper.NewBadRequestError("parameter followed_user_name is required")
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
	followRepo := repository.NewFollowRepository(db)

	// UseCase
	u := usecase.NewDeleteFollow(r.Context(), tx, userName, followedUserName, userRepo, followRepo)
	if err = u.DeleteFollowUseCase(); err != nil {
		log.Println(err)
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
