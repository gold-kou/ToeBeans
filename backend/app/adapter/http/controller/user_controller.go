package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gold-kou/ToeBeans/backend/app"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func UserController(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/users":
		switch r.Method {
		case http.MethodGet:
			user, postingCount, likeCount, likedCount, followCount, followedCount, err := getUser(r)
			switch err := err.(type) {
			case nil:
				resp := modelHTTP.ResponseGetUser{
					UserName:         user.Name,
					Icon:             user.Icon,
					SelfIntroduction: user.SelfIntroduction,
					PostingCount:     postingCount,
					LikeCount:        likeCount,
					LikedCount:       likedCount,
					FollowCount:      followCount,
					FollowedCount:    followedCount,
					CreatedAt:        user.CreatedAt,
				}
				w.Header().Set(helper.HeaderKeyContentType, helper.HeaderValueApplicationJSON)
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					log.Println(fmt.Errorf("error: " + err.Error()))
				}
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
			methods := []string{http.MethodGet}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	case strings.HasPrefix(r.URL.Path, "/users/"):
		switch r.Method {
		case http.MethodPost:
			err := registerUser(r)
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
		case http.MethodPut:
			err := updateUser(r)
			switch err := err.(type) {
			case nil:
				helper.ResponseSimpleSuccess(w)
			case *helper.BadRequestError:
				helper.ResponseBadRequest(w, err.Error())
			case *helper.AuthorizationError:
				helper.ResponseUnauthorized(w, err.Error())
			case *helper.ForbiddenError:
				helper.ResponseForbidden(w, err.Error())
			case *helper.ConflictError:
				helper.ResponseConflictError(w, err.Error())
			case *helper.InternalServerError:
				helper.ResponseInternalServerError(w, err.Error())
			default:
				helper.ResponseInternalServerError(w, err.Error())
			}
		case http.MethodDelete:
			err := deleteUser(r)
			switch err := err.(type) {
			case nil:
				helper.ResponseSimpleSuccess(w)
			case *helper.BadRequestError:
				helper.ResponseBadRequest(w, err.Error())
			case *helper.AuthorizationError:
				helper.ResponseUnauthorized(w, err.Error())
			case *helper.ForbiddenError:
				helper.ResponseForbidden(w, err.Error())
			case *helper.ConflictError:
				helper.ResponseConflictError(w, err.Error())
			case *helper.InternalServerError:
				helper.ResponseInternalServerError(w, err.Error())
			default:
				helper.ResponseInternalServerError(w, err.Error())
			}
		default:
			methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	case strings.HasPrefix(r.URL.Path, "/user-activation/"):
		switch r.Method {
		case http.MethodGet:
			err := activateUser(r)
			switch err := err.(type) {
			case nil:
				var link string
				if app.IsProduction() {
					link = "https://toebeans.ml/login"
				} else {
					link = "http://localhost:3000/login"
				}
				params := map[string]string{
					"loginLink": link,
				}
				w.Header().Set(helper.HeaderKeyContentType, helper.HeaderValueHTML)
				t, err := template.ParseFiles("../view/template/user-activation.html")
				if err != nil {
					log.Println(err.Error())
				}
				if err := t.Execute(w, params); err != nil {
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
			methods := []string{http.MethodGet}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	default:
		helper.ResponseInternalServerError(w, errMsgControllerPath)
	}
}

func registerUser(r *http.Request) error {
	// get request parameter
	vars := mux.Vars(r)
	userName, _ := vars["user_name"]
	var reqRegisterUser *modelHTTP.RequestRegisterUser
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	defer r.Body.Close()
	if err = json.Unmarshal(b, &reqRegisterUser); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}

	// validation check
	if err = validation.Validate(userName, validation.Required, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	err = reqRegisterUser.ValidateParam()
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
	u := usecase.NewRegisterUser(tx, userName, reqRegisterUser, userRepo)
	if err = u.RegisterUserUseCase(r.Context()); err != nil {
		log.Println(err)
		if err == usecase.ErrDuplicateData {
			return helper.NewBadRequestError(err.Error() + ", the user name or email has been already used.")
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}

func getUser(r *http.Request) (user model.User, postingCount, likeCount, likedCount, followCount, followedCount int64, err error) {
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		log.Println(err)
		err = helper.NewInternalServerError(err.Error())
		return
	}

	// get request parameter
	targetUserName := r.URL.Query().Get("user_name")

	// validation check
	if targetUserName != "" {
		if err = validation.Validate(targetUserName, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
			log.Println(err)
			err = helper.NewBadRequestError(err.Error())
			return
		}
	} else {
		// パラメータuser_nameが指定されていなければIDトークンのユーザ名を使う
		targetUserName = tokenUserName
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
	positngRepo := repository.NewPostingRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	followRepo := repository.NewFollowRepository(db)

	// UseCase
	u := usecase.NewGetUser(tx, tokenUserName, targetUserName, userRepo, positngRepo, likeRepo, followRepo)
	if user, postingCount, likeCount, likedCount, followCount, followedCount, err = u.GetUserUseCase(r.Context()); err != nil {
		log.Println(err)
		if err == usecase.ErrNotExistsData {
			err = helper.NewNotFoundError(err.Error())
			return
		}
		if err == usecase.ErrTokenInvalidNotExistingUserName {
			err = helper.NewAuthorizationError(err.Error())
			return
		}
		err = helper.NewInternalServerError(err.Error())
		return
	}
	return
}

func updateUser(r *http.Request) (err error) {
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
	userName, _ := vars["user_name"]
	var reqUpdateUser *modelHTTP.RequestUpdateUser
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	defer r.Body.Close()
	if err = json.Unmarshal(b, &reqUpdateUser); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}

	// validation check
	if err = validation.Validate(userName, validation.Required, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	if reqUpdateUser.Password != "" {
		if err = validation.Validate(reqUpdateUser.Password, validation.By(modelHTTP.PasswordValidation)); err != nil {
			log.Println(err)
			return helper.NewBadRequestError(err.Error())
		}
	}
	if reqUpdateUser.SelfIntroduction != "" {
		if err = validation.Validate(reqUpdateUser.SelfIntroduction, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength)); err != nil {
			log.Println(err)
			return helper.NewBadRequestError(err.Error())
		}
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
	u := usecase.NewUpdateUser(tx, userName, reqUpdateUser, userRepo)
	if err = u.UpdateUserUseCase(r.Context()); err != nil {
		log.Println(err)
		if err == usecase.ErrNotExistsData || err == usecase.ErrDecodeImage {
			return helper.NewBadRequestError(err.Error())
		}
		if err == usecase.ErrNotExitsUser {
			return helper.NewConflictError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return
}

func deleteUser(r *http.Request) (err error) {
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

	// get path parameter
	vars := mux.Vars(r)
	userName, _ := vars["user_name"]

	// validation check
	if err = validation.Validate(userName, validation.Required, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
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
	passwordResetRepo := repository.NewPasswordResetRepository(db)
	postingRepo := repository.NewPostingRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	followRepo := repository.NewFollowRepository(db)

	// UseCase
	u := usecase.NewDeleteUser(tx, userName, userRepo, passwordResetRepo, postingRepo, likeRepo, commentRepo, followRepo)
	if err = u.DeleteUserUseCase(r.Context()); err != nil {
		log.Println(err)
		if err == usecase.ErrNotExitsUser {
			return helper.NewConflictError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return
}

func activateUser(r *http.Request) (err error) {
	// get path parameter
	vars := mux.Vars(r)
	userName, _ := vars["user_name"]
	activationKey, _ := vars["activation_key"]

	// validation check
	if err = validation.Validate(userName, validation.Required, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	if err = validation.Validate(activationKey, validation.Required, is.UUID); err != nil {
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
	u := usecase.NewUserActivation(tx, userName, activationKey, userRepo)
	if err = u.UserActivationUseCase(r.Context()); err != nil {
		log.Println(err)
		if err == repository.ErrUserActivationNotFound {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return
}
