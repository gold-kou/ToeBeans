package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
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
		case http.MethodGet:
			exists, err := getFollowState(r)
			switch err := err.(type) {
			case nil:
				resp := modelHTTP.ResponseGetFollowState{
					IsFollow: exists,
				}
				w.Header().Set(helper.HeaderKeyContentType, helper.HeaderValueApplicationJSON)
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					panic(err.Error())
				}
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
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		return helper.NewAuthorizationError(err.Error())
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
	u := usecase.NewRegisterFollow(tx, tokenUserName, followedUserName, userRepo, followRepo, notificationRepo)
	if err = u.RegisterFollowUseCase(r.Context()); err != nil {
		log.Println(err)
		if err == usecase.ErrAlreadyFollowed {
			return helper.NewConflictError("Whoops, you already followed the user")
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}

// return follow or not
func getFollowState(r *http.Request) (bool, error) {
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		return false, helper.NewAuthorizationError(err.Error())
	}

	// get request parameter
	vars := mux.Vars(r)
	followedUserName, _ := vars["followed_user_name"]

	// validation check
	if err = validation.Validate(followedUserName, validation.Required, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
		log.Println(err)
		return false, helper.NewBadRequestError(err.Error())
	}

	// db connect
	db, err := mysql.NewDB()
	if err != nil {
		log.Println(err)
		return false, helper.NewInternalServerError(err.Error())
	}
	defer db.Close()
	tx := mysql.NewDBTransaction(db)

	// repository
	userRepo := repository.NewUserRepository(db)
	followRepo := repository.NewFollowRepository(db)

	// UseCase
	u := usecase.NewGetFollowState(tx, tokenUserName, followedUserName, userRepo, followRepo)
	if err = u.GetFollowStateUseCase(r.Context()); err != nil {
		if err == usecase.ErrNotExistsData {
			return false, nil
		}
		log.Println(err)
		return false, helper.NewInternalServerError(err.Error())
	}
	return true, nil
}

func deleteFollow(r *http.Request) error {
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		return helper.NewAuthorizationError(err.Error())
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
	u := usecase.NewDeleteFollow(tx, tokenUserName, followedUserName, userRepo, followRepo)
	if err = u.DeleteFollowUseCase(r.Context()); err != nil {
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
