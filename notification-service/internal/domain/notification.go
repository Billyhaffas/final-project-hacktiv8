package domain

import (
	"context"
)

type NotificationUsecase interface {
	CheckAndSendNotification(ctx context.Context, userID int) (bool, string, error)
}
