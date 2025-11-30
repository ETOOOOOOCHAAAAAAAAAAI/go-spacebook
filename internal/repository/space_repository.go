package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"SpaceBookProject/internal/domain"
)

type SpaceFilter struct {
	Query    *string
	MinPrice *int
	MaxPrice *int
	MinArea  *float64
	MaxArea  *float64
}

var ErrSpaceNotFound = errors.New("space not found")

type SpaceRepository struct {
	db *sql.DB
}

func NewSpaceRepository(db *sql.DB) *SpaceRepository {
	return &SpaceRepository{db: db}
}

func (r *SpaceRepository) GetByID(id int) (*domain.Space, error) {
	const query = `
        SELECT id, owner_id, title, description, area_m2, price, phone, created_at, updated_at
        FROM spaces
        WHERE id = $1
    `

	s := &domain.Space{}
	err := r.db.QueryRow(query, id).Scan(
		&s.ID,
		&s.OwnerID,
		&s.Title,
		&s.Description,
		&s.AreaM2,
		&s.Price,
		&s.Phone,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrSpaceNotFound
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SpaceRepository) ListFiltered(f SpaceFilter) ([]domain.Space, error) {
	query := `
		SELECT id, owner_id, title, description, area_m2, price, phone, created_at, updated_at
		FROM spaces
	`
	var (
		conds []string
		args  []any
		i     = 1
	)

	if f.Query != nil && *f.Query != "" {
		conds = append(conds, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", i, i+1))
		pattern := "%" + *f.Query + "%"
		args = append(args, pattern, pattern)
		i += 2
	}
	if f.MinPrice != nil {
		conds = append(conds, fmt.Sprintf("price >= $%d", i))
		args = append(args, *f.MinPrice)
		i++
	}
	if f.MaxPrice != nil {
		conds = append(conds, fmt.Sprintf("price <= $%d", i))
		args = append(args, *f.MaxPrice)
		i++
	}
	if f.MinArea != nil {
		conds = append(conds, fmt.Sprintf("area_m2 >= $%d", i))
		args = append(args, *f.MinArea)
		i++
	}
	if f.MaxArea != nil {
		conds = append(conds, fmt.Sprintf("area_m2 <= $%d", i))
		args = append(args, *f.MaxArea)
		i++
	}

	if len(conds) > 0 {
		query += " WHERE " + strings.Join(conds, " AND ")
	}

	query += " ORDER BY created_at DESC, id DESC"

	rows, err := r.db.Query(query, args...)
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

func (r *SpaceRepository) Create(space *domain.Space) error {
	now := time.Now()

	query := `
		INSERT INTO spaces (owner_id, title, description, area_m2, price, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		space.OwnerID,
		space.Title,
		space.Description,
		space.AreaM2,
		space.Price,
		space.Phone,
		now,
		now,
	).Scan(&space.ID, &space.CreatedAt, &space.UpdatedAt)

	return err
}
