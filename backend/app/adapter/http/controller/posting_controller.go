package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func PostingController(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/postings":
		switch r.Method {
		case http.MethodPost:
			err := registerPosting(r)
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
		case http.MethodGet:
			postings, likedCounts, likes, err := getPostings(r)
			switch err := err.(type) {
			case nil:
				var httpPostings []modelHTTP.ResponseGetPosting
				var resp modelHTTP.ResponseGetPostings
				for i, p := range postings {
					httpPosting := modelHTTP.ResponseGetPosting{
						PostingId:  p.ID,
						UserName:   p.UserName,
						UploadedAt: p.CreatedAt,
						Title:      p.Title,
						ImageUrl:   p.ImageURL,
						LikedCount: likedCounts[i],
						Liked:      false,
					}
					for _, l := range likes {
						if p.ID == l.PostingID {
							httpPosting.Liked = true
						}
					}
					httpPostings = append(httpPostings, httpPosting)
				}
				resp = modelHTTP.ResponseGetPostings{
					Postings: httpPostings,
				}
				w.Header().Set(helper.HeaderKeyContentType, helper.HeaderValueApplicationJSON)
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					log.Println(err.Error())
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
			methods := []string{http.MethodPost, http.MethodGet}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	case strings.HasPrefix(r.URL.Path, "/postings/"):
		switch r.Method {
		case http.MethodDelete:
			err := deletePosting(r)
			switch err := err.(type) {
			case nil:
				helper.ResponseSimpleSuccess(w)
			case *helper.BadRequestError:
				helper.ResponseBadRequest(w, err.Error())
			case *helper.AuthorizationError:
				helper.ResponseUnauthorized(w, err.Error())
			case *helper.ForbiddenError:
				helper.ResponseForbidden(w, err.Error())
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

func registerPosting(r *http.Request) error {
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		log.Println(err)
		return helper.NewInternalServerError(err.Error())
	}

	// get request parameter
	var reqRegisterPosting *modelHTTP.RequestRegisterPosting
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	defer r.Body.Close()
	if err := json.Unmarshal(b, &reqRegisterPosting); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}

	// validation check
	err = reqRegisterPosting.ValidateParam()
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	if strings.Contains(reqRegisterPosting.Title, "_") {
		return helper.NewBadRequestError("title: must not contain _.")
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

	// UseCase
	u := usecase.NewRegisterPosting(r.Context(), tx, tokenUserName, reqRegisterPosting, userRepo, postingRepo)
	if err = u.RegisterPostingUseCase(); err != nil {
		log.Println(err)
		if err == usecase.ErrDecodeImage || err == usecase.ErrNotCatImage {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}

func getPostings(r *http.Request) (postings []model.Posting, likedCounts []int64, likes []model.Like, err error) {
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		log.Println(err)
		err = helper.NewInternalServerError(err.Error())
		return
	}

	// get request parameter
	sinceAt := r.URL.Query().Get("since_at")
	if sinceAt == "" {
		log.Println(err)
		err = helper.NewBadRequestError("since_at: cannot be blank.")
		return
	}

	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 例えば 2020-01-01T00:00:00+09:00 でリクエストされても 2020-01-01T00:00:00 09:00 と変換されてしまうため、それをreplaceしている。
	sinceAtFormatted, err := time.ParseInLocation("2006-01-02T15:04:05+09:00", strings.Replace(sinceAt, " ", "+", 1), jst)
	if err != nil {
		log.Println(err)
		err = helper.NewBadRequestError(err.Error())
		return
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		log.Println(err)
		err = helper.NewBadRequestError("limit: cannot be blank.")
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Println(err)
		err = helper.NewBadRequestError(err.Error())
		return
	}

	// here user means selected user to see user profile
	userName := r.URL.Query().Get("user_name")
	if err = validation.Validate(userName, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
		log.Println(err)
		err = helper.NewBadRequestError(err.Error())
		return
	}

	// db connect
	db, err := mysql.NewDB()
	if err != nil {
		log.Println(err)
		err = helper.NewInternalServerError(err.Error())
		return
	}
	defer db.Close()
	tx := mysql.NewDBTransaction(db)

	// repository
	userRepo := repository.NewUserRepository(db)
	postingRepo := repository.NewPostingRepository(db)
	likeRepo := repository.NewLikeRepository(db)

	// UseCase
	u := usecase.NewGetPostings(r.Context(), tx, tokenUserName, sinceAtFormatted, int8(limitInt), userName, userRepo, postingRepo, likeRepo)
	if postings, likedCounts, likes, err = u.GetPostingsUseCase(); err != nil {
		log.Println(err)
		if err == usecase.ErrDecodeImage {
			err = helper.NewBadRequestError(err.Error())
			return
		}
		if err == usecase.ErrNotExistsData {
			err = helper.NewBadRequestError(err.Error())
			return
		}
		err = helper.NewInternalServerError(err.Error())
		return
	}
	return
}

func deletePosting(r *http.Request) error {
	// not allowed to guest user
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		log.Println(err)
		return helper.NewInternalServerError(err.Error())
	}
	if tokenUserName == helper.GuestUserName {
		log.Println(errMsgGuestUserForbidden)
		return helper.NewForbiddenError(errMsgGuestUserForbidden)
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

	// UseCase
	u := usecase.NewDeletePosting(r.Context(), tx, int64(postingID), tokenUserName, userRepo, postingRepo)
	if err = u.DeletePostingUseCase(); err != nil {
		log.Println(err)
		if err == usecase.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
