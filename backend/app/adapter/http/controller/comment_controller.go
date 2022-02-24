package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"

	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func CommentController(w http.ResponseWriter, r *http.Request) {
	switch {
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
	case strings.HasPrefix(r.URL.Path, "/comments/"):
		switch r.Method {
		case http.MethodPost:
			err := registerComment(r)
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
		case http.MethodDelete:
			err := deleteComment(r)
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
			methods := []string{http.MethodPost, http.MethodDelete}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	default:
		helper.ResponseInternalServerError(w, errMsgControllerPath)
	}
}

func registerComment(r *http.Request) error {
	// not allowed to guest user
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		log.Println(err)
		return helper.NewInternalServerError(err.Error())
	}
	if tokenUserName == lib.GuestUserName {
		log.Println(errMsgGuestUserForbidden)
		return helper.NewForbiddenError(errMsgGuestUserForbidden)
	}

	// get request parameter
	vars := mux.Vars(r)
	paramPostingID, _ := vars["posting_id"]
	postingID, err := strconv.Atoi(paramPostingID)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
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
	if err = validation.Validate(postingID, validation.Required); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
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
	u := usecase.NewRegisterComment(r.Context(), tx, tokenUserName, postingID, reqRegisterComment, userRepo, postingRepo, commentRepo, notificationRepo)
	if err = u.RegisterCommentUseCase(); err != nil {
		log.Println(err)
		if err == usecase.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		if err == usecase.ErrDuplicateData {
			return helper.NewBadRequestError("Whoops, you already commented that")
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}

func getComments(r *http.Request) (comments []model.Comment, err error) {
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
	tokenUserName, e := context.GetTokenUserName(r.Context())
	if e != nil {
		log.Println(err)
		return nil, helper.NewInternalServerError(err.Error())
	}
	u := usecase.NewGetComments(r.Context(), tx, tokenUserName, int64(id), userRepo, postingRepo, commentRepo)
	if comments, err = u.GetCommentsUseCase(); err != nil {
		log.Println(err)
		if err == usecase.ErrNotExistsData {
			return nil, helper.NewBadRequestError(err.Error())
		}
		return nil, helper.NewInternalServerError(err.Error())
	}
	return
}

func deleteComment(r *http.Request) error {
	// not allowed to guest user
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		log.Println(err)
		return helper.NewInternalServerError(err.Error())
	}
	if tokenUserName == lib.GuestUserName {
		log.Println(errMsgGuestUserForbidden)
		return helper.NewForbiddenError(errMsgGuestUserForbidden)
	}

	// get request parameter
	vars := mux.Vars(r)
	paramCommentID, _ := vars["comment_id"]
	commentID, err := strconv.Atoi(paramCommentID)
	if err != nil {
		return helper.NewInternalServerError(err.Error())
	}

	// validation check
	if err = validation.Validate(commentID, validation.Required); err != nil {
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
	commentRepo := repository.NewCommentRepository(db)

	// UseCase
	u := usecase.NewDeleteComment(r.Context(), tx, tokenUserName, int64(commentID), userRepo, commentRepo)
	if err = u.DeleteCommentUseCase(); err != nil {
		log.Println(err)
		if err == usecase.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
