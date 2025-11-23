package repository

import (
	"database/sql"
	"time"

	"SpaceBookProject/internal/domain"
)

type SpaceRepository struct {
	db *sql.DB
}

func NewSpaceRepository(db *sql.DB) *SpaceRepository {
	return &SpaceRepository{db: db}
}

func (r *SpaceRepository) Create(ownerID int, req *domain.CreateSpaceRequest) (*domain.Space, error) {
	query := `
        INSERT INTO spaces (owner_id, title, description, area_m2, price, phone, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
        RETURNING id, created_at, updated_at;
    `
	now := time.Now()

	space := &domain.Space{
		OwnerID:     ownerID,
		Title:       req.Title,
		Description: req.Description,
		AreaM2:      req.AreaM2,
		Price:       req.Price,
		Phone:       req.Phone,
		IsActive:    true,
	}

	err := r.db.QueryRow(
		query,
		ownerID,
		req.Title,
		req.Description,
		req.AreaM2,
		req.Price,
		req.Phone,
		now,
	).Scan(&space.ID, &space.CreatedAt, &space.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return space, nil
}

func (r *SpaceRepository) ListActive() ([]domain.Space, error) {
	query := `
        SELECT id, owner_id, title, description, area_m2, price, phone, is_active, created_at, updated_at
        FROM spaces
        WHERE is_active = TRUE
        ORDER BY created_at DESC;
    `
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spaces []domain.Space
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
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		spaces = append(spaces, s)
	}
	return spaces, nil
}

func (r *SpaceRepository) ListByOwner(ownerID int) ([]domain.Space, error) {
	query := `
        SELECT id, owner_id, title, description, area_m2, price, phone, is_active, created_at, updated_at
        FROM spaces
        WHERE owner_id = $1
        ORDER BY created_at DESC;
    `
	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spaces []domain.Space
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
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		spaces = append(spaces, s)
	}
	return spaces, nil
}
