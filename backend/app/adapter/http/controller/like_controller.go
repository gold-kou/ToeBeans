package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gold-kou/ToeBeans/backend/app/lib"
	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func LikeController(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/like":
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
			case *helper.InternalServerError:
				helper.ResponseInternalServerError(w, err.Error())
			default:
				helper.ResponseInternalServerError(w, err.Error())
			}
		default:
			methods := []string{http.MethodPost}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	case strings.HasPrefix(r.URL.Path, "/like/"):
		switch r.Method {
		case http.MethodDelete:
			err := deleteLike(r)
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
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	default:
		helper.ResponseInternalServerError(w, errMsgControllerPath)
	}
}

func registerLike(r *http.Request) error {
	// authorization
	cookie, err := r.Cookie(helper.CookieIDToken)
	if err != nil {
		log.Println(err)
		return helper.NewAuthorizationError(err.Error())
	}
	tokenUserName, err := lib.VerifyToken(cookie.Value)
	if err != nil {
		return helper.NewAuthorizationError(err.Error())
	}

	// get request parameter
	var reqRegisterLike *modelHTTP.Like
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	defer r.Body.Close()
	if err := json.Unmarshal(b, &reqRegisterLike); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}

	// validation check
	err = reqRegisterLike.ValidateParam()
	if err != nil {
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
	u := usecase.NewRegisterLike(r.Context(), tx, tokenUserName, reqRegisterLike, userRepo, postingRepo, likeRepo, notificationRepo)
	if err = u.RegisterLikeUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrDuplicateData {
			return helper.NewBadRequestError("Whoops, you already liked the posting")
		}
		if err == usecase.ErrLikeYourSelf {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}

func deleteLike(r *http.Request) error {
	// authorization
	cookie, err := r.Cookie(helper.CookieIDToken)
	if err != nil {
		log.Println(err)
		return helper.NewAuthorizationError(err.Error())
	}
	tokenUserName, err := lib.VerifyToken(cookie.Value)
	if err != nil {
		return helper.NewAuthorizationError(err.Error())
	}

	// get request parameter
	vars := mux.Vars(r)
	paramPostingID, ok := vars["posting_id"]
	if !ok || paramPostingID == "0" {
		log.Println(err)
		return helper.NewBadRequestError("posting_id: cannot be blank")
	}
	postingID, err := strconv.Atoi(paramPostingID)
	if err != nil {
		return helper.NewInternalServerError(err.Error())
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
	u := usecase.NewDeleteLike(r.Context(), tx, tokenUserName, int64(postingID), userRepo, postingRepo, likeRepo)
	if err = u.DeleteLikeUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
