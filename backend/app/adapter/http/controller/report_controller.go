package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func ReportController(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/reports/users/"):
		switch r.Method {
		case http.MethodPost:
			err := submitUserReport(r)
			switch err := err.(type) {
			case nil:
				helper.ResponseSimpleSuccess(w)
			case *helper.BadRequestError:
				helper.ResponseBadRequest(w, err.Error())
			case *helper.AuthorizationError:
				helper.ResponseUnauthorized(w, err.Error())
			case *helper.NotFoundError:
				helper.ResponseNotFound(w, err.Error())
			case *helper.InternalServerError:
				helper.ResponseInternalServerError(w, err.Error())
			default:
				helper.ResponseInternalServerError(w, err.Error())
			}
		default:
			methods := []string{http.MethodPost}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	case strings.HasPrefix(r.URL.Path, "/reports/postings/"):
		switch r.Method {
		case http.MethodPost:
			err := submitPostingReport(r)
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
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	default:
		helper.ResponseInternalServerError(w, errMsgControllerPath)
	}
}

func submitUserReport(r *http.Request) error {
	// get request parameter
	vars := mux.Vars(r)
	userName, _ := vars["user_name"]
	var reqSubmitUserReport *modelHTTP.RequestSubmitUserReport
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	defer r.Body.Close()
	if err = json.Unmarshal(b, &reqSubmitUserReport); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}

	// validation check
	if err = validation.Validate(userName, validation.Required); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	err = reqSubmitUserReport.ValidateParam()
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
	userReportRepo := repository.NewUserReportRepository(db)

	// UseCase
	u := usecase.NewSubmitUserReport(tx, userName, reqSubmitUserReport, userRepo, userReportRepo)
	if err = u.SubmitUserReportUseCase(r.Context()); err != nil {
		log.Println(err)
		if err == usecase.ErrNotExitsUser {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}

func submitPostingReport(r *http.Request) error {
	// get request parameter
	vars := mux.Vars(r)
	postingIDStr, _ := vars["posting_id"]
	postingID, err := strconv.Atoi(postingIDStr)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	var reqSubmitPostingReport *modelHTTP.RequestSubmitPostingReport
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	defer r.Body.Close()
	if err = json.Unmarshal(b, &reqSubmitPostingReport); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}

	// validation check
	if err = validation.Validate(postingID, validation.Required); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	err = reqSubmitPostingReport.ValidateParam()
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
	postingReportRepo := repository.NewPostingReportRepository(db)

	// UseCase
	u := usecase.NewSubmitPostingReport(tx, postingID, reqSubmitPostingReport, postingRepo, postingReportRepo)
	if err = u.SubmitPostingReportUseCase(r.Context()); err != nil {
		log.Println(err)
		if err == usecase.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
