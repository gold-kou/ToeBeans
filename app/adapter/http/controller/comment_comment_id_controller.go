package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/app/lib"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"
	applicationLog "github.com/gold-kou/ToeBeans/app/adapter/http/log"
	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/application/usecase"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

func CommentCommentIDController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

	switch r.Method {
	case http.MethodDelete:
		err = deleteComment(r)
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

func deleteComment(r *http.Request) error {
	// authorization
	tokenUserName, err := lib.VerifyHeaderToken(r)
	if err != nil {
		log.Println(err)
		return helper.NewAuthorizationError(err.Error())
	}

	// get request parameter
	vars := mux.Vars(r)
	commentID, ok := vars["comment_id"]
	if !ok || commentID == "0" {
		log.Println(err)
		return helper.NewBadRequestError("parameter comment_id is required")
	}
	id, err := strconv.Atoi(commentID)
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
	commentRepo := repository.NewCommentRepository(db)

	// UseCase
	u := usecase.NewDeleteComment(r.Context(), tx, tokenUserName, int64(id), userRepo, commentRepo)
	if err = u.DeleteCommentUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
