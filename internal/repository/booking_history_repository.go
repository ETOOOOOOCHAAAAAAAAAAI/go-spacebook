package repository

import (
	"database/sql"

	"SpaceBookProject/internal/domain"
)

type BookingHistoryRepository struct {
	db *sql.DB
}

func NewBookingHistoryRepository(db *sql.DB) *BookingHistoryRepository {
	return &BookingHistoryRepository{db: db}
}

func (r *BookingHistoryRepository) Add(bookingID int, oldStatus *string, newStatus string, changedBy int) error {
	const q = `
  INSERT INTO booking_status_history (booking_id, old_status, new_status, changed_by)
  VALUES ($1, $2, $3, $4)
 `
	_, err := r.db.Exec(q, bookingID, oldStatus, newStatus, changedBy)
	return err
}

func (r *BookingHistoryRepository) List(bookingID int) ([]domain.BookingStatusHistory, error) {
	const q = `
  SELECT id, booking_id, old_status, new_status, changed_by, changed_at
  FROM booking_status_history
  WHERE booking_id = $1
  ORDER BY changed_at ASC, id ASC
 `

	rows, err := r.db.Query(q, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.BookingStatusHistory
	for rows.Next() {
		var h domain.BookingStatusHistory
		if err := rows.Scan(&h.ID, &h.BookingID, &h.OldStatus, &h.NewStatus, &h.ChangedBy, &h.ChangedAt); err != nil {
			return nil, err
		}
		res = append(res, h)
	}
	return res, rows.Err()
}
