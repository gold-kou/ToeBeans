package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gold-kou/ToeBeans/backend/app/domain/model"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	applicationLog "github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func CommentsController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

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
		helper.ResponseNotAllowedMethod(w, "not allowed method", methods)
	}
}

func getComments(r *http.Request) (comments []model.Comment, err error) {
	// authorization
	//tokenUserName, err := lib.VerifyHeaderToken(r)
	//if err != nil {
	//	log.Println(err)
	//	return nil, helper.NewAuthorizationError(err.Error())
	//}
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
