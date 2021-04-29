package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	applicationLog "github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
	"github.com/gold-kou/ToeBeans/backend/app/lib"
)

func PasswordController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

	switch r.URL.Path {
	case "/password":
		switch r.Method {
		case http.MethodPut:
			err = changePassword(r)
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
			methods := []string{http.MethodPut}
			helper.ResponseNotAllowedMethod(w, errMsgNotAllowedMethod, methods)
		}
	case "/password-reset":
		switch r.Method {
		case http.MethodPost:
			err := passwordReset(r)
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
	case "/password-reset-email":
		switch r.Method {
		case http.MethodPost:
			err := passwordResetEmail(r)
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

func changePassword(r *http.Request) (err error) {
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
	var reqChangePassword *modelHTTP.RequestChangePassword
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	defer r.Body.Close()
	if err = json.Unmarshal(b, &reqChangePassword); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}

	// validation check
	err = reqChangePassword.ValidateParam()
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
	u := usecase.NewChangePassword(r.Context(), tx, tokenUserName, reqChangePassword, userRepo)
	if err = u.ChangePasswordUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData || err == usecase.ErrNotCorrectPassword {
			return helper.NewBadRequestError(err.Error())
		}
		if err == lib.ErrTokenInvalidNotExistingUserName {
			return helper.NewAuthorizationError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return
}

func passwordResetEmail(r *http.Request) error {
	// get request parameter
	var reqPasswordResetEmail *modelHTTP.Email
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	defer r.Body.Close()
	if err = json.Unmarshal(b, &reqPasswordResetEmail); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}

	// validation check
	err = reqPasswordResetEmail.ValidateParam()
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
	re := usecase.NewPasswordResetEmail(r.Context(), tx, reqPasswordResetEmail, userRepo)
	if err = re.PasswordResetEmailUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		if err == usecase.ErrOverPasswordResetCount {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}

func passwordReset(r *http.Request) error {
	// get request parameter
	var reqResetPassword *modelHTTP.RequestResetPassword
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}
	defer r.Body.Close()
	if err = json.Unmarshal(b, &reqResetPassword); err != nil {
		log.Println(err)
		return helper.NewBadRequestError(err.Error())
	}

	// validation check
	err = reqResetPassword.ValidateParam()
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
	re := usecase.NewPasswordReset(r.Context(), tx, reqResetPassword, userRepo)
	if err = re.PasswordResetUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return helper.NewBadRequestError(errMsgUserNameResetKeyNotExistsResetKeyExpired)
		}
		return helper.NewInternalServerError(err.Error())
	}
	return err
}
