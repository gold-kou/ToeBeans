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
	ctx              context.Context
	tx               mysql.DBTransaction
	tokenUserName    string
	visitedName      string
	userRepo         *repository.UserRepository
	notificationRepo *repository.NotificationRepository
}

func NewGetNotifications(ctx context.Context, tx mysql.DBTransaction, tokenUserName, visitedName string, userRepo *repository.UserRepository, notificationRepo *repository.NotificationRepository) *GetNotifications {
	return &GetNotifications{
		ctx:              ctx,
		tx:               tx,
		tokenUserName:    tokenUserName,
		visitedName:      visitedName,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

func (n *GetNotifications) GetNotificationsUseCase() (notifications []model.Notification, err error) {
	notifications, err = n.notificationRepo.GetNotifications(n.ctx, n.visitedName)
	if err != nil {
		return
	}
	return
}
