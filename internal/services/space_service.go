package services

import (
	"SpaceBookProject/internal/domain"
	"SpaceBookProject/internal/repository"
)

type SpaceService struct {
	repo *repository.SpaceRepository
}

func NewSpaceService(repo *repository.SpaceRepository) *SpaceService {
	return &SpaceService{repo: repo}
}

func (s *SpaceService) CreateSpace(ownerID int, req *domain.CreateSpaceRequest) (*domain.Space, error) {
	return s.repo.Create(ownerID, req)
}

func (s *SpaceService) ListActiveSpaces() ([]domain.Space, error) {
	return s.repo.ListActive()
}

func (s *SpaceService) ListOwnerSpaces(ownerID int) ([]domain.Space, error) {
	return s.repo.ListByOwner(ownerID)
}
