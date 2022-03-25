package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type GetNotificationsUseCaseInterface interface {
	GetNotificationsUseCase() (*model.Notification, error)
}

type GetNotifications struct {
	tx               mysql.DBTransaction
	tokenUserName    string
	visitedName      string
	userRepo         *repository.UserRepository
	notificationRepo *repository.NotificationRepository
}

func NewGetNotifications(tx mysql.DBTransaction, tokenUserName, visitedName string, userRepo *repository.UserRepository, notificationRepo *repository.NotificationRepository) *GetNotifications {
	return &GetNotifications{
		tx:               tx,
		tokenUserName:    tokenUserName,
		visitedName:      visitedName,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

func (n *GetNotifications) GetNotificationsUseCase(ctx context.Context) (notifications []model.Notification, err error) {
	notifications, err = n.notificationRepo.GetNotifications(ctx, n.visitedName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return nil, ErrNotExistsData
		}
		return
	}
	return
}
