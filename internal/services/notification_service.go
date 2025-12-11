package services

import "SpaceBookProject/internal/domain"
import "SpaceBookProject/internal/repository"

type NotificationService interface {
    CreateNotification(n *domain.Notification) error
    ListUserNotifications(userID int) ([]domain.Notification, error)
}

type notificationService struct {
    repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
    return &notificationService{repo: repo}
}

func (s *notificationService) CreateNotification(n *domain.Notification) error {
    return s.repo.Create(n)
}

func (s *notificationService) ListUserNotifications(userID int) ([]domain.Notification, error) {
    return s.repo.ListByUser(userID)
}
