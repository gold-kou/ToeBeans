package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type GetNotificationsUseCaseInterface interface {
	GetNotificationsUseCase() (*model.Notification, error)
}

type GetNotifications struct {
	ctx              context.Context
	tx               mysql.DBTransaction
	visitedName      string
	notificationRepo *repository.NotificationRepository
}

func NewGetNotifications(ctx context.Context, tx mysql.DBTransaction, visitedName string, notificationRepo *repository.NotificationRepository) *GetNotifications {
	return &GetNotifications{
		ctx:              ctx,
		tx:               tx,
		visitedName:      visitedName,
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
