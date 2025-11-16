package services

import (
	"errors"
	"time"

	"SpaceBookProject/internal/auth"
	"SpaceBookProject/internal/domain"
	"SpaceBookProject/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTokenExpired       = errors.New("token has expired")
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthService(userRepo *repository.UserRepository, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, repository.ErrUserAlreadyExists
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.userRepo.SaveRefreshToken(user.ID, refreshToken, expiresAt); err != nil {
		return nil, err
	}
	user.PasswordHash = ""

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}
func (s *AuthService) Login(req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.userRepo.SaveRefreshToken(user.ID, refreshToken, expiresAt); err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}
func (s *AuthService) RefreshToken(refreshToken string) (*domain.AuthResponse, error) {
	userID, err := s.userRepo.GetRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.userRepo.SaveRefreshToken(user.ID, newRefreshToken, expiresAt); err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         *user,
	}, nil
}

func (s *AuthService) Logout(userID int) error {
	return s.userRepo.DeleteRefreshToken(userID)
}

func (s *AuthService) GetUserByID(userID int) (*domain.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *AuthService) ValidateToken(token string) (*auth.TokenClaims, error) {
	return s.jwtManager.ValidateToken(token)
}
