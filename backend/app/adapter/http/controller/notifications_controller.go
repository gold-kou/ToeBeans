package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/application/usecase"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

func NotificationsController(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		notifications, err := getNotifications(r)
		switch err := err.(type) {
		case nil:
			var httpNotifications []modelHTTP.ResponseGetNotification
			for _, n := range notifications {
				httpNotification := modelHTTP.ResponseGetNotification{
					VisitorName: n.VisitorName,
					ActionType:  n.Action,
					CreatedAt:   n.CreatedAt,
				}
				httpNotifications = append(httpNotifications, httpNotification)
			}
			resp := modelHTTP.ResponseGetNotifications{
				VisitedName: notifications[0].VisitedName,
				Actions:     httpNotifications,
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
}

func getNotifications(r *http.Request) (notifications []model.Notification, err error) {
	tokenUserName, err := context.GetTokenUserName(r.Context())
	if err != nil {
		log.Println(err)
		return nil, helper.NewInternalServerError(err.Error())
	}

	// get request parameter
	visitedName := r.URL.Query().Get("user_name")

	// validation check
	if err = validation.Validate(visitedName, validation.Required, validation.Length(modelHTTP.MinVarcharLength, modelHTTP.MaxVarcharLength), is.Alphanumeric); err != nil {
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
	notificationRepo := repository.NewNotificationRepository(db)

	// UseCase
	u := usecase.NewGetNotifications(r.Context(), tx, tokenUserName, visitedName, userRepo, notificationRepo)
	if notifications, err = u.GetNotificationsUseCase(); err != nil {
		log.Println(err)
		if err == usecase.ErrNotExistsData {
			return nil, helper.NewBadRequestError(err.Error())
		}
		return nil, helper.NewInternalServerError(err.Error())
	}
	return
}
