package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	applicationLog "github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func PostingsController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

	switch r.Method {
	case http.MethodGet:
		postings, err := getPostings(r)
		switch err := err.(type) {
		case nil:
			var httpPostings []modelHTTP.ResponseGetPosting
			var resp modelHTTP.ResponseGetPostings
			if len(postings) >= 1 {
				for _, p := range postings {
					httpPosting := modelHTTP.ResponseGetPosting{
						PostingId:  p.ID,
						UserName:   p.UserName,
						UploadedAt: p.CreatedAt,
						Title:      p.Title,
						ImageUrl:   p.ImageURL,
						Liked:      p.LikedCount,
					}
					httpPostings = append(httpPostings, httpPosting)
				}
				resp = modelHTTP.ResponseGetPostings{
					Postings: httpPostings,
				}
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
		case *helper.InternalServerError:
			helper.ResponseInternalServerError(w, err.Error())
		default:
			helper.ResponseInternalServerError(w, err.Error())
		}
	default:
		methods := []string{http.MethodGet}
		helper.ResponseNotAllowedMethod(w, "not allowed method", methods)
	}
}

func getPostings(r *http.Request) (postings []model.Posting, err error) {
	// authorization
	tokenUserName, err := lib.VerifyHeaderToken(r)
	if err != nil {
		log.Println(err)
		return nil, helper.NewAuthorizationError(err.Error())
	}

	// get request parameter
	sinceAt := r.URL.Query().Get("since_at")
	if sinceAt == "" {
		log.Println(err)
		return nil, helper.NewBadRequestError("since_at: cannot be blank.")
	}
	sinceAtFormatted, err := time.Parse(time.RFC3339, sinceAt)
	if err != nil {
		log.Println(err)
		return nil, helper.NewBadRequestError(err.Error())
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		log.Println(err)
		return nil, helper.NewBadRequestError("limit: cannot be blank.")
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Println(err)
		return nil, helper.NewBadRequestError(err.Error())
	}

	userName := r.URL.Query().Get("user_name")
	if err := validation.Validate(userName, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
		log.Println(err)
		return nil, helper.NewBadRequestError(err.Error())
	}

	// db connect
	db, err := mysql.NewDB()
	if err != nil {
		log.Println(err)
		return nil, helper.NewInternalServerError(err.Error())
	}
	defer db.Close()
	tx := mysql.NewDBTransaction(db)

	// repository
	userRepo := repository.NewUserRepository(db)
	postingRepo := repository.NewPostingRepository(db)

	// UseCase
	u := usecase.NewGetPostings(r.Context(), tx, tokenUserName, sinceAtFormatted, int8(limitInt), userName, userRepo, postingRepo)
	if postings, err = u.GetPostingsUseCase(); err != nil {
		log.Println(err)
		if err == usecase.ErrDecodeImage {
			return nil, helper.NewBadRequestError(err.Error())
		}
		if err == repository.ErrNotExistsData {
			return nil, helper.NewBadRequestError(err.Error())
		}
		return nil, helper.NewInternalServerError(err.Error())
	}
	return
}
