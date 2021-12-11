package controller

import (
	"log"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"

	"github.com/gold-kou/ToeBeans/backend/app/lib"
	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func FollowController(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/follows/"):
		switch r.Method {
		case http.MethodPost:
			err := registerFollow(r)
			switch err := err.(type) {
			case nil:
				helper.ResponseSimpleSuccess(w)
			case *helper.BadRequestError:
				helper.ResponseBadRequest(w, err.Error())
			case *helper.AuthorizationError:
				helper.ResponseUnauthorized(w, err.Error())
			case *helper.ForbiddenError:
				helper.ResponseForbidden(w, err.Error())
			case *helper.ConflictError:
				helper.ResponseConflictError(w, err.Error())
			case *helper.InternalServerError:
				helper.ResponseInternalServerError(w, err.Error())
			default:
				helper.ResponseInternalServerError(w, err.Error())
			}
		case http.MethodDelete:
			err := deleteFollow(r)
			switch err := err.(type) {
			case nil:
				helper.ResponseSimpleSuccess(w)
			case *helper.BadRequestError:
				helper.ResponseBadRequest(w, err.Error())
			case *helper.AuthorizationError:
				helper.ResponseUnauthorized(w, err.Error())
			case *helper.ForbiddenError:
				helper.ResponseForbidden(w, err.Error())
			case *helper.ConflictError:
				helper.ResponseConflictError(w, err.Error())
			case *helper.InternalServerError:
				helper.ResponseInternalServerError(w, err.Error())
			default:
				helper.ResponseInternalServerError(w, err.Error())
			}
		default:
			methods := []string{http.MethodPost, http.MethodDelete}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	default:
		helper.ResponseInternalServerError(w, errMsgControllerPath)
	}
}

func registerFollow(r *http.Request) error {
	// not allowed to guest user
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if tokenUserName == lib.GuestUserName {
		log.Println(errMsgGuestUserForbidden)
		return helper.NewForbiddenError(errMsgGuestUserForbidden)
	}

	// get request parameter
	vars := mux.Vars(r)
	followedUserName, _ := vars["followed_user_name"]

	// validation check
	if err = validation.Validate(followedUserName, validation.Required, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
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
	followRepo := repository.NewFollowRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// UseCase
	u := usecase.NewRegisterFollow(r.Context(), tx, tokenUserName, followedUserName, userRepo, followRepo, notificationRepo)
	if err = u.RegisterFollowUseCase(); err != nil {
		log.Println(err)
		if err == usecase.ErrAlreadyFollowed {
			return helper.NewConflictError("Whoops, you already followed the user")
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}

func deleteFollow(r *http.Request) error {
	// not allowed to guest user
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		return helper.NewAuthorizationError(err.Error())
	}
	if tokenUserName == lib.GuestUserName {
		log.Println(errMsgGuestUserForbidden)
		return helper.NewForbiddenError(errMsgGuestUserForbidden)
	}

	// get request parameter
	vars := mux.Vars(r)
	followedUserName, _ := vars["followed_user_name"]

	// validation check
	if err = validation.Validate(followedUserName, validation.Required, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
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
	followRepo := repository.NewFollowRepository(db)

	// UseCase
	u := usecase.NewDeleteFollow(r.Context(), tx, tokenUserName, followedUserName, userRepo, followRepo)
	if err = u.DeleteFollowUseCase(); err != nil {
		log.Println(err)
		if err == usecase.ErrNotExitsUser {
			return helper.NewConflictError(err.Error())
		} else if err == usecase.ErrDeleteNotExistsFollow {
			return helper.NewConflictError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
