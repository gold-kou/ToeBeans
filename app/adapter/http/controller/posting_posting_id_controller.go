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

func PostingPostingIDController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

	switch r.Method {
	case http.MethodDelete:
		err = deletePosting(r)
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
		helper.ResponseBadRequest(w, "not allowed method")
	}
}

func deletePosting(r *http.Request) error {
	// authorization
	userName, err := lib.VerifyHeaderToken(r)
	if err != nil {
		log.Println(err)
		return helper.NewAuthorizationError(err.Error())
	}

	// get request parameter
	vars := mux.Vars(r)
	postingID, ok := vars["posting_id"]
	if !ok {
		log.Println(err)
		return helper.NewBadRequestError("parameter posting_id is required")
	}
	id, err := strconv.Atoi(postingID)
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
	postingRepo := repository.NewPostingRepository(db)

	// UseCase
	u := usecase.NewDeletePosting(r.Context(), tx, uint64(id), userName, postingRepo)
	if err = u.DeletePostingUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
