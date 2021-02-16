package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	applicationLog "github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
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
		case *helper.ForbiddenError:
			helper.ResponseForbidden(w, err.Error())
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
	tokenUserName, err := lib.VerifyHeaderToken(r)
	if err != nil {
		log.Println(err)
		return helper.NewAuthorizationError(err.Error())
	}
	if tokenUserName == lib.GuestUserName {
		log.Println(errMsgGuestUserForbidden)
		return helper.NewForbiddenError(errMsgGuestUserForbidden)
	}

	// get request parameter
	vars := mux.Vars(r)
	followedUserName, ok := vars["followed_user_name"]
	if !ok || followedUserName == "" {
		log.Println(err)
		return helper.NewBadRequestError("followed_user_name: cannot be blank.")
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
	u := usecase.NewDeleteFollow(r.Context(), tx, tokenUserName, followedUserName, userRepo, followRepo)
	if err = u.DeleteFollowUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
