package repository

import (
	"database/sql"
	"errors"
	"time"

	"SpaceBookProject/internal/domain"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, role, first_name, last_name, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, password_hash, role, first_name, last_name, phone, created_at, updated_at
		FROM users
		WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(id int) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, password_hash, role, first_name, last_name, phone, created_at, updated_at
		FROM users
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, phone = $3, updated_at = $4
		WHERE id = $5`

	user.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		query,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) SaveRefreshToken(userID int, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) 
		DO UPDATE SET token = $2, expires_at = $3, created_at = $4`

	_, err := r.db.Exec(query, userID, token, expiresAt, time.Now())
	return err
}

func (r *UserRepository) GetRefreshToken(token string) (int, error) {
	var userID int
	var expiresAt time.Time

	query := `
		SELECT user_id, expires_at
		FROM refresh_tokens
		WHERE token = $1`

	err := r.db.QueryRow(query, token).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrUserNotFound
		}
		return 0, err
	}

	if time.Now().After(expiresAt) {
		return 0, errors.New("refresh token expired")
	}

	return userID, nil
}

func (r *UserRepository) DeleteRefreshToken(userID int) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}
