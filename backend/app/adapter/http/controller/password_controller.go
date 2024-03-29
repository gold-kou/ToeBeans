package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

var errMsgPasswordResetKeyWrong = "password reset key is wrong"
var errMsgPasswordResetKeyExpired = "password reset key is expired"

func PasswordController(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/password":
		switch r.Method {
		case http.MethodPut:
			err := changePassword(r)
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
	u := usecase.NewChangePassword(tx, tokenUserName, reqChangePassword, userRepo)
	if err = u.ChangePasswordUseCase(r.Context()); err != nil {
		log.Println(err)
		if err == usecase.ErrNotCorrectPassword {
			return helper.NewBadRequestError(err.Error())
		}
		if err == usecase.ErrTokenInvalidNotExistingUserName {
			return helper.NewAuthorizationError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return
}

func passwordResetEmail(r *http.Request) error {
	// get request parameter
	var reqPasswordResetEmail *modelHTTP.RequestSendPasswordResetEmail
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
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	// UseCase
	re := usecase.NewPasswordResetEmail(tx, reqPasswordResetEmail, userRepo, passwordResetRepo)
	if err = re.PasswordResetEmailUseCase(r.Context()); err != nil {
		log.Println(err)
		if err == usecase.ErrNotExistsData {
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
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	// UseCase
	re := usecase.NewPasswordReset(tx, reqResetPassword, userRepo, passwordResetRepo)
	if err = re.PasswordResetUseCase(r.Context()); err != nil {
		log.Println(err)
		switch err {
		case usecase.ErrNotExistsData:
			return helper.NewBadRequestError(errMsgNotExists)
		case usecase.ErrPasswordResetKeyWrong:
			return helper.NewBadRequestError(errMsgPasswordResetKeyWrong)
		case usecase.ErrPasswordResetKeyExpired:
			return helper.NewBadRequestError(errMsgPasswordResetKeyExpired)
		default:
			return helper.NewInternalServerError(err.Error())
		}
	}
	return err
}
