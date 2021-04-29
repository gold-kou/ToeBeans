package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	applicationLog "github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func CommentController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

	switch {
	case r.URL.Path == "/comment":
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
			case *helper.ForbiddenError:
				helper.ResponseForbidden(w, err.Error())
			case *helper.InternalServerError:
				helper.ResponseInternalServerError(w, err.Error())
			default:
				helper.ResponseInternalServerError(w, err.Error())
			}
		default:
			methods := []string{http.MethodPost}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	case r.URL.Path == "/comments":
		switch r.Method {
		case http.MethodGet:
			comments, err := getComments(r)
			switch err := err.(type) {
			case nil:
				var httpComments []modelHTTP.ResponseGetComment
				var resp modelHTTP.ResponseGetComments
				if len(comments) >= 1 {
					for _, c := range comments {
						httpComment := modelHTTP.ResponseGetComment{
							CommentId:   c.ID,
							UserName:    c.UserName,
							CommentedAt: c.CreatedAt,
							Comment:     c.Comment,
						}
						httpComments = append(httpComments, httpComment)
					}
					resp = modelHTTP.ResponseGetComments{
						PostingId: comments[0].PostingID,
						Comments:  httpComments,
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
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	case strings.HasPrefix(r.URL.Path, "/comment/"):
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

func registerComment(r *http.Request) error {
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
	userRepo := repository.NewUserRepository(db)
	postingRepo := repository.NewPostingRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// UseCase
	u := usecase.NewRegisterComment(r.Context(), tx, tokenUserName, reqRegisterComment, userRepo, postingRepo, commentRepo, notificationRepo)
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

func getComments(r *http.Request) (comments []model.Comment, err error) {
	// authorization
	cookie, err := r.Cookie(helper.CookieIDToken)
	if err != nil {
		log.Println(err)
		return nil, helper.NewAuthorizationError(err.Error())
	}
	tokenUserName, err := lib.VerifyToken(cookie.Value)
	if err != nil {
		return nil, helper.NewAuthorizationError(err.Error())
	}

	// get request parameter
	postingID := r.URL.Query().Get("posting_id")
	if postingID == "" {
		log.Println(err)
		return nil, helper.NewBadRequestError("posting_id: cannot be blank.")
	}
	id, err := strconv.Atoi(postingID)
	if err != nil {
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
	commentRepo := repository.NewCommentRepository(db)

	// UseCase
	u := usecase.NewGetComments(r.Context(), tx, tokenUserName, int64(id), userRepo, postingRepo, commentRepo)
	if comments, err = u.GetCommentsUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return nil, helper.NewBadRequestError(err.Error())
		}
		return nil, helper.NewInternalServerError(err.Error())
	}
	return
}

func deleteComment(r *http.Request) error {
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
