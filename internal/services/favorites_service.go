package services

import (
	"context"

	"SpaceBookProject/internal/domain"
	"SpaceBookProject/internal/repository"
)

type FavoritesService struct {
	repo *repository.FavoritesRepository
}

func NewFavoritesService(repo *repository.FavoritesRepository) *FavoritesService {
	return &FavoritesService{repo: repo}
}

func (s *FavoritesService) Add(ctx context.Context, userID, spaceID int) error {
	return s.repo.Add(ctx, userID, spaceID)
}

func (s *FavoritesService) Remove(ctx context.Context, userID, spaceID int) error {
	return s.repo.Remove(ctx, userID, spaceID)
}

func (s *FavoritesService) List(ctx context.Context, userID int) ([]domain.Space, error) {
	return s.repo.List(ctx, userID)
}
