package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"
	applicationLog "github.com/gold-kou/ToeBeans/app/adapter/http/log"
	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/application/usecase"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
	"github.com/gold-kou/ToeBeans/app/lib"
)

func UserController(w http.ResponseWriter, r *http.Request) {
	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	l.LogHTTPAccess(r)

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
				UserName:     user.Name,
				Icon:         user.Icon,
				PostingCount: user.PostingCount,
				LikeCount:    user.LikeCount,
				LikedCount:   user.LikedCount,
				CreatedAt:    user.CreatedAt,
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
	case http.MethodPut:
		err = updateUser(r)
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
	case http.MethodDelete:
		err := deleteUser(r)
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
		methods := []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete}
		helper.ResponseNotAllowedMethod(w, "not allowed method", methods)
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
	_, err = lib.VerifyHeaderToken(r)
	if err != nil {
		log.Println(err)
		return model.User{}, helper.NewAuthorizationError(err.Error())
	}

	// get request parameter
	userName := r.URL.Query().Get("user_name")

	// validation check
	if err = validation.Validate(userName, validation.Required, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
		log.Println(err)
		return model.User{}, helper.NewBadRequestError(err.Error())
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
	u := usecase.NewGetUser(r.Context(), tx, userName, userRepo)
	if user, err = u.GetUserUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return model.User{}, helper.NewBadRequestError(err.Error())
		}
		return model.User{}, helper.NewInternalServerError(err.Error())
	}
	return
}

func updateUser(r *http.Request) (err error) {
	// authorization
	userName, err := lib.VerifyHeaderToken(r)
	if err != nil {
		log.Println(err)
		return helper.NewAuthorizationError(err.Error())
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
	err = reqUpdateUser.ValidateParam()
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
	u := usecase.NewUpdateUser(r.Context(), tx, userName, reqUpdateUser, userRepo)
	if err = u.UpdateUserUseCase(); err != nil {
		log.Println(err)
		if err == repository.ErrNotExistsData {
			return helper.NewBadRequestError(err.Error())
		}
		return helper.NewInternalServerError(err.Error())
	}
	return
}

func deleteUser(r *http.Request) (err error) {
	// authorization
	userName, err := lib.VerifyHeaderToken(r)
	if err != nil {
		log.Println(err)
		return helper.NewAuthorizationError(err.Error())
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
	u := usecase.NewDeleteUser(r.Context(), tx, userName, userRepo, postingRepo, likeRepo, commentRepo, followRepo)
	if err = u.DeleteUserUseCase(); err != nil {
		log.Println(err)
		return helper.NewInternalServerError(err.Error())
	}
	return
}
