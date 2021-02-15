package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/repository"

	"github.com/gold-kou/ToeBeans/app/application/usecase"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"

	applicationLog "github.com/gold-kou/ToeBeans/app/adapter/http/log"
	model "github.com/gold-kou/ToeBeans/app/domain/model/http"
)

func PasswordResetEmailController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

	switch r.Method {
	case http.MethodPost:
		err := passwordResetEmail(r)
		switch err := err.(type) {
		case nil:
			helper.ResponseSimpleSuccess(w)
		case *helper.BadRequestError:
			helper.ResponseBadRequest(w, err.Error())
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

func passwordResetEmail(r *http.Request) error {
	// get request parameter
	var reqPasswordResetEmail *model.Email
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	defer r.Body.Close()
	if err = json.Unmarshal(b, &reqPasswordResetEmail); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}

	// validation check
	err = reqPasswordResetEmail.ValidateParam()
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

	// UseCase
	re := usecase.NewPasswordResetEmail(r.Context(), tx, reqPasswordResetEmail, userRepo)
	if err = re.PasswordResetEmailUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		if err == usecase.ErrOverPasswordResetCount {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
