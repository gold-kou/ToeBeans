package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	applicationLog "github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
	"github.com/gold-kou/ToeBeans/backend/app/lib"
)

func UserController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

	switch {
	case r.URL.Path == "/user":
		switch r.Method {
		case http.MethodPost:
			err = registerUser(r)
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
		case http.MethodGet:
			user, err := getUser(r)
			switch err := err.(type) {
			case nil:
				resp := modelHTTP.ResponseGetUser{
					UserName:         user.Name,
					Icon:             user.Icon,
					SelfIntroduction: user.SelfIntroduction,
					PostingCount:     user.PostingCount,
					LikeCount:        user.LikeCount,
					LikedCount:       user.LikedCount,
					FollowCount:      user.FollowCount,
					FollowedCount:    user.FollowedCount,
					CreatedAt:        user.CreatedAt,
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
			case *helper.NotFoundError:
				helper.ResponseNotFound(w, err.Error())
			case *helper.InternalServerError:
				helper.ResponseInternalServerError(w, err.Error())
			default:
				helper.ResponseInternalServerError(w, err.Error())
			}
		case http.MethodPut:
			err = updateUser(r)
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
			case *helper.InternalServerError:
				helper.ResponseInternalServerError(w, err.Error())
			default:
				helper.ResponseInternalServerError(w, err.Error())
			}
		default:
			methods := []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete}
			helper.ResponseNotAllowedMethod(w, "not allowed method", methods)
		}
	case strings.HasPrefix(r.URL.Path, "/user-activation/"):
		switch r.Method {
		case http.MethodGet:
			err = activateUser(r)
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
			methods := []string{http.MethodGet}
			helper.ResponseNotAllowedMethod(w, "not allowed method", methods)
		}
	default:
		helper.ResponseInternalServerError(w, errMsgControllerPath)
	}
}

func registerUser(r *http.Request) error {
	// get request parameter
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
	u := usecase.NewRegisterUser(r.Context(), tx, reqRegisterUser, userRepo)
	if err = u.RegisterUserUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrDuplicateData {
			return helper.NewBadRequestError(err.Error() + ", the user name or email have been already used.")
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}

func getUser(r *http.Request) (user model.User, err error) {
	// authorization
	cookie, err := r.Cookie(helper.CookieIDToken)
	if err != nil {
		log.Println(err)
		return user, helper.NewAuthorizationError(err.Error())
	}
	tokenUserName, err := lib.VerifyToken(cookie.Value)
	if err != nil {
		return user, helper.NewAuthorizationError(err.Error())
	}

	// get request parameter
	userName := r.URL.Query().Get("user_name")

	// validation check
	if userName != "" {
		if err = validation.Validate(userName, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
			log.Println(err)
			return model.User{}, helper.NewBadRequestError(err.Error())
		}
	} else {
		// パラメータuser_nameが指定されていなければIDトークンのユーザ名を使う
		userName = tokenUserName
	}

	// db connect
	db, err := mysql.NewDB()
	if err != nil {
		log.Println(err)
		return model.User{}, helper.NewInternalServerError(err.Error())
	}
	defer db.Close()
	tx := mysql.NewDBTransaction(db)

	// repository
	userRepo := repository.NewUserRepository(db)

	// UseCase
	u := usecase.NewGetUser(r.Context(), tx, tokenUserName, userName, userRepo)
	if user, err = u.GetUserUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return model.User{}, helper.NewNotFoundError(err.Error())
		}
		if err == lib.ErrTokenInvalidNotExistingUserName {
			return model.User{}, helper.NewAuthorizationError(err.Error())
		}
		return model.User{}, helper.NewInternalServerError(err.Error())
	}
	return
}

func updateUser(r *http.Request) (err error) {
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
	u := usecase.NewUpdateUser(r.Context(), tx, tokenUserName, reqUpdateUser, userRepo)
	if err = u.UpdateUserUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData || err == usecase.ErrDecodeImage {
			return helper.NewBadRequestError(err.Error())
		}
		if err == lib.ErrTokenInvalidNotExistingUserName {
			return helper.NewAuthorizationError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return
}

func deleteUser(r *http.Request) (err error) {
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
	likeRepo := repository.NewLikeRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	followRepo := repository.NewFollowRepository(db)

	// UseCase
	u := usecase.NewDeleteUser(r.Context(), tx, tokenUserName, userRepo, postingRepo, likeRepo, commentRepo, followRepo)
	if err = u.DeleteUserUseCase(); err != nil {
		log.Println(err)
		if err == lib.ErrTokenInvalidNotExistingUserName {
			return helper.NewAuthorizationError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return
}

func activateUser(r *http.Request) (err error) {
	// get path parameter
	vars := mux.Vars(r)
	userName, ok := vars["user_name"]
	if !ok || userName == "" {
		log.Println(err)
		return helper.NewBadRequestError("user_name cannot be blank")
	}
	activationKey, ok := vars["activation_key"]
	if !ok || activationKey == "" {
		log.Println(err)
		return helper.NewBadRequestError("activation_key cannot be blank")
	}

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
	u := usecase.NewUserActivation(r.Context(), tx, userName, activationKey, userRepo)
	if err = u.UserActivationUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrUserActivationNotFound {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return
}
