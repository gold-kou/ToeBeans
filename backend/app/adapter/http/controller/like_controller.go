package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"

	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func LikeController(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/likes/"):
		switch r.Method {
		case http.MethodPost:
			err := registerLike(r)
			switch err := err.(type) {
			case nil:
				helper.ResponseSimpleSuccess(w)
			case *helper.BadRequestError:
				helper.ResponseBadRequest(w, err.Error())
			case *helper.AuthorizationError:
				helper.ResponseUnauthorized(w, err.Error())
			case *helper.ConflictError:
				helper.ResponseConflictError(w, err.Error())
			case *helper.InternalServerError:
				helper.ResponseInternalServerError(w, err.Error())
			default:
				helper.ResponseInternalServerError(w, err.Error())
			}
		case http.MethodDelete:
			err := deleteLike(r)
			switch err := err.(type) {
			case nil:
				helper.ResponseSimpleSuccess(w)
			case *helper.BadRequestError:
				helper.ResponseBadRequest(w, err.Error())
			case *helper.AuthorizationError:
				helper.ResponseUnauthorized(w, err.Error())
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

func registerLike(r *http.Request) error {
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		log.Println(err)
		return helper.NewInternalServerError(err.Error())
	}

	// get request parameter
	vars := mux.Vars(r)
	paramPostingID, _ := vars["posting_id"]
	postingID, err := strconv.Atoi(paramPostingID)
	if err != nil {
		return helper.NewInternalServerError(err.Error())
	}

	// validation check
	if err = validation.Validate(postingID, validation.Required); err != nil {
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
	postingRepo := repository.NewPostingRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// UseCase
	tokenUserID, err := context.GetTokenUserID(r.Context())
	if err != nil {
		log.Println(err)
		return helper.NewInternalServerError(err.Error())
	}
	u := usecase.NewRegisterLike(tx, tokenUserID, tokenUserName, postingID, userRepo, postingRepo, likeRepo, notificationRepo)
	if err = u.RegisterLikeUseCase(r.Context()); err != nil {
		log.Println(err)
		if err == usecase.ErrLikeYourPosting {
			return helper.NewConflictError(err.Error())
		} else if err == usecase.ErrAlreadyLiked {
			return helper.NewConflictError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}

func deleteLike(r *http.Request) error {
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		log.Println(err)
		return helper.NewInternalServerError(err.Error())
	}

	// get request parameter
	vars := mux.Vars(r)
	paramPostingID, _ := vars["posting_id"]
	postingID, err := strconv.Atoi(paramPostingID)
	if err != nil {
		return helper.NewInternalServerError(err.Error())
	}

	// validation check
	if err = validation.Validate(postingID, validation.Required); err != nil {
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
	postingRepo := repository.NewPostingRepository(db)
	likeRepo := repository.NewLikeRepository(db)

	// UseCase
	u := usecase.NewDeleteLike(tx, tokenUserName, int64(postingID), userRepo, postingRepo, likeRepo)
	if err = u.DeleteLikeUseCase(r.Context()); err != nil {
		log.Println(err)
		if err == usecase.ErrDeleteNotExistsLike {
			return helper.NewConflictError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
