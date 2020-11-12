package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gold-kou/ToeBeans/app/lib"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"
	applicationLog "github.com/gold-kou/ToeBeans/app/adapter/http/log"
	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

func CommentController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

	switch r.Method {
	case http.MethodPost:
		err = registerComment(r)
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
		helper.ResponseNotAllowedMethod(w, "not allowed method", methods)
	}
}

func registerComment(r *http.Request) error {
	// authorization
	userName, err := lib.VerifyHeaderToken(r)
	if err != nil {
		log.Println(err)
		return helper.NewAuthorizationError(err.Error())
	}

	// get request parameter
	var reqRegisterComment *modelHTTP.Comment
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	defer r.Body.Close()
	if err := json.Unmarshal(b, &reqRegisterComment); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}

	// validation check
	err = reqRegisterComment.ValidateParam()
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
	postingRepo := repository.NewPostingRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// UseCase
	u := usecase.NewRegisterComment(r.Context(), tx, userName, reqRegisterComment, postingRepo, commentRepo, notificationRepo)
	if err = u.RegisterCommentUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		if err == repository.ErrDuplicateData {
			return helper.NewBadRequestError("Whoops, you already commented that")
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
