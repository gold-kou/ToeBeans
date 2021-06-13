package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func PostingController(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/posting":
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
		default:
			methods := []string{http.MethodPost}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	case r.URL.Path == "/postings":
		switch r.Method {
		case http.MethodGet:
			postings, likes, err := getPostings(r)
			switch err := err.(type) {
			case nil:
				var httpPostings []modelHTTP.ResponseGetPosting
				var resp modelHTTP.ResponseGetPostings
				for _, p := range postings {
					httpPosting := modelHTTP.ResponseGetPosting{
						PostingId:  p.ID,
						UserName:   p.UserName,
						UploadedAt: p.CreatedAt,
						Title:      p.Title,
						ImageUrl:   p.ImageURL,
						LikedCount: p.LikedCount,
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
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	case strings.HasPrefix(r.URL.Path, "/posting/"):
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

func getPostings(r *http.Request) (postings []model.Posting, likes []model.Like, err error) {
	// authorization
	cookie, err := r.Cookie(helper.CookieIDToken)
	if err != nil {
		log.Println(err)
		return nil, nil, helper.NewAuthorizationError(err.Error())
	}
	tokenUserName, err := lib.VerifyToken(cookie.Value)
	if err != nil {
		return nil, nil, helper.NewAuthorizationError(err.Error())
	}

	// get request parameter
	sinceAt := r.URL.Query().Get("since_at")
	if sinceAt == "" {
		log.Println(err)
		return nil, nil, helper.NewBadRequestError("since_at: cannot be blank.")
	}

	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 例えば 2020-01-01T00:00:00+09:00 でリクエストされても 2020-01-01T00:00:00 09:00 と変換されてしまうため、それをreplaceしている。
	sinceAtFormatted, err := time.ParseInLocation("2006-01-02T15:04:05+09:00", strings.Replace(sinceAt, " ", "+", 1), jst)
	if err != nil {
		log.Println(err)
		return nil, nil, helper.NewBadRequestError(err.Error())
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		log.Println(err)
		return nil, nil, helper.NewBadRequestError("limit: cannot be blank.")
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Println(err)
		return nil, nil, helper.NewBadRequestError(err.Error())
	}

	// here user means selected user to see user profile
	userName := r.URL.Query().Get("user_name")
	if err := validation.Validate(userName, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
		log.Println(err)
		return nil, nil, helper.NewBadRequestError(err.Error())
	}

	// db connect
	db, err := mysql.NewDB()
	if err != nil {
		log.Println(err)
		return nil, nil, helper.NewInternalServerError(err.Error())
	}
	defer db.Close()
	tx := mysql.NewDBTransaction(db)

	// repository
	userRepo := repository.NewUserRepository(db)
	postingRepo := repository.NewPostingRepository(db)
	likeRepo := repository.NewLikeRepository(db)

	// UseCase
	u := usecase.NewGetPostings(r.Context(), tx, tokenUserName, sinceAtFormatted, int8(limitInt), userName, userRepo, postingRepo, likeRepo)
	if postings, likes, err = u.GetPostingsUseCase(); err != nil {
		log.Println(err)
		if err == usecase.ErrDecodeImage {
			return nil, nil, helper.NewBadRequestError(err.Error())
		}
		if err == repository.ErrNotExistsData {
			return nil, nil, helper.NewBadRequestError(err.Error())
		}
		return nil, nil, helper.NewInternalServerError(err.Error())
	}
	return
}

func deletePosting(r *http.Request) error {
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
	if tokenUserName == lib.GuestUserName {
		log.Println(errMsgGuestUserForbidden)
		return helper.NewForbiddenError(errMsgGuestUserForbidden)
	}

	// get request parameter
	vars := mux.Vars(r)
	postingID, ok := vars["posting_id"]
	if !ok || postingID == "0" {
		log.Println(err)
		return helper.NewBadRequestError("posting_id: cannot be blank")
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
	userRepo := repository.NewUserRepository(db)
	postingRepo := repository.NewPostingRepository(db)

	// UseCase
	u := usecase.NewDeletePosting(r.Context(), tx, int64(id), tokenUserName, userRepo, postingRepo)
	if err = u.DeletePostingUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
