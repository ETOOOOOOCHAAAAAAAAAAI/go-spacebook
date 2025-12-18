package repository

import (
	"database/sql"
	"time"

	"SpaceBookProject/internal/domain"
)

type BookingHistoryRepository struct {
	db *sql.DB
}

func NewBookingHistoryRepository(db *sql.DB) *BookingHistoryRepository {
	return &BookingHistoryRepository{db: db}
}

// Add saves booking status change to history
func (r *BookingHistoryRepository) Add(
	bookingID int,
	oldStatus *string,
	newStatus string,
	changedBy int,
) error {

	query := `
		INSERT INTO booking_status_history
			(booking_id, old_status, new_status, changed_by, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		query,
		bookingID,
		oldStatus,
		newStatus,
		changedBy,
		time.Now(),
	)

	return err
}

// List returns status change history for booking
func (r *BookingHistoryRepository) List(
	bookingID int,
) ([]domain.BookingStatusHistory, error) {

	rows, err := r.db.Query(`
		SELECT
			id,
			booking_id,
			old_status,
			new_status,
			changed_by,
			created_at
		FROM booking_status_history
		WHERE booking_id = $1
		ORDER BY created_at ASC
	`, bookingID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.BookingStatusHistory

	for rows.Next() {
		var h domain.BookingStatusHistory
		if err := rows.Scan(
			&h.ID,
			&h.BookingID,
			&h.OldStatus,
			&h.NewStatus,
			&h.ChangedBy,
			&h.CreatedAt,
		); err != nil {
			return nil, err
		}

		history = append(history, h)
	}

	return history, nil
}
