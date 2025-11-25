package repository

import (
	"database/sql"
	"errors"
	"time"

	"SpaceBookProject/internal/domain"
)

var (
	ErrBookingNotFound = errors.New("booking not found")
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(b *domain.Booking) error {
	const query = `
		INSERT INTO bookings (space_id, tenant_id, date_from, date_to, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, status, created_at, updated_at;
	`

	err := r.db.QueryRow(
		query,
		b.SpaceID,
		b.TenantID,
		b.DateFrom,
		b.DateTo,
		b.Status,
	).Scan(&b.ID, &b.Status, &b.CreatedAt, &b.UpdatedAt)

	return err
}

func (r *BookingRepository) GetByID(id int) (*domain.Booking, error) {
	const q = `
        SELECT id, space_id, tenant_id, date_from, date_to, status, created_at, updated_at
        FROM bookings
        WHERE id = $1`

	b := &domain.Booking{}
	err := r.db.QueryRow(q, id).Scan(
		&b.ID, &b.SpaceID, &b.TenantID,
		&b.DateFrom, &b.DateTo, &b.Status,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBookingNotFound
		}
		return nil, err
	}
	return b, nil
}

func (r *BookingRepository) ListByTenant(tenantID int) ([]domain.Booking, error) {
	const q = `
        SELECT id, space_id, tenant_id, date_from, date_to, status, created_at, updated_at
        FROM bookings
        WHERE tenant_id = $1
        ORDER BY date_from DESC, id DESC`

	rows, err := r.db.Query(q, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(
			&b.ID, &b.SpaceID, &b.TenantID,
			&b.DateFrom, &b.DateTo, &b.Status,
			&b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	return res, rows.Err()
}

func (r *BookingRepository) ListByOwner(ownerID int) ([]domain.Booking, error) {
	const q = `
        SELECT b.id, b.space_id, b.tenant_id, b.date_from, b.date_to,
               b.status, b.created_at, b.updated_at
        FROM bookings b
        JOIN spaces s ON s.id = b.space_id
        WHERE s.owner_id = $1
        ORDER BY b.date_from DESC, b.id DESC`

	rows, err := r.db.Query(q, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(
			&b.ID, &b.SpaceID, &b.TenantID,
			&b.DateFrom, &b.DateTo, &b.Status,
			&b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	return res, rows.Err()
}

func (r *BookingRepository) UpdateStatus(id int, status domain.BookingStatus) error {
	const q = `
        UPDATE bookings
        SET status = $1, updated_at = NOW()
        WHERE id = $2`

	res, err := r.db.Exec(q, status, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrBookingNotFound
	}
	return nil
}

func (r *BookingRepository) HasApprovedOverlap(
	spaceID int,
	from, to time.Time,
	excludeID *int,
) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM bookings
            WHERE space_id = $1
              AND status = 'approved'
              -- нет пересечения = (date_to <= from) OR (date_from >= to)
              -- нам нужны ИМЕННО пересекающиеся, поэтому NOT (...)
              AND NOT (date_to <= $2 OR date_from >= $3)
    `
	args := []any{spaceID, from, to}

	if excludeID != nil {
		query += " AND id <> $4"
		args = append(args, *excludeID)
	}

	query += ")"

	var exists bool
	if err := r.db.QueryRow(query, args...).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
