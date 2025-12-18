package repository

import (
	"context"
	"database/sql"

	"SpaceBookProject/internal/domain"
)

type FavoritesRepository struct {
	db *sql.DB
}

func NewFavoritesRepository(db *sql.DB) *FavoritesRepository {
	return &FavoritesRepository{db: db}
}

func (r *FavoritesRepository) Add(ctx context.Context, userID, spaceID int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO favorites (user_id, space_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, space_id) DO NOTHING
	`, userID, spaceID)
	return err
}

func (r *FavoritesRepository) Remove(ctx context.Context, userID, spaceID int) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM favorites
		WHERE user_id = $1 AND space_id = $2
	`, userID, spaceID)
	return err
}

func (r *FavoritesRepository) List(ctx context.Context, userID int) ([]domain.Space, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.owner_id, s.title, s.description, s.area_m2, s.price, s.phone, s.created_at, s.updated_at
		FROM spaces s
		JOIN favorites f ON f.space_id = s.id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC, s.id DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Space
	for rows.Next() {
		var s domain.Space
		if err := rows.Scan(
			&s.ID,
			&s.OwnerID,
			&s.Title,
			&s.Description,
			&s.AreaM2,
			&s.Price,
			&s.Phone,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}
