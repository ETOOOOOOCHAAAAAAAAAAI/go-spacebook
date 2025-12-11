package repository

import "SpaceBookProject/internal/domain"

type NotificationRepository interface {
    Create(n *domain.Notification) error
    ListByUser(userID int) ([]domain.Notification, error)
}
