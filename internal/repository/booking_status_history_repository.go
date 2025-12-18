package repository

import (
	"database/sql"
	"time"

	"SpaceBookProject/internal/domain"
)

type BookingStatusHistoryRepository struct {
	db sqlDB
}

// sqlDB объединяет *sql.DB и *sql.Tx
type sqlDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

func NewBookingStatusHistoryRepository(db sqlDB) *BookingStatusHistoryRepository {
	return &BookingStatusHistoryRepository{db: db}
}

// Create записывает изменение статуса в историю
func (r *BookingStatusHistoryRepository) Create(history *domain.BookingStatusHistory) error {
	const query = `
		INSERT INTO booking_status_history 
		(booking_id, old_status, new_status, changed_by, reason, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at;
	`

	return r.db.QueryRow(
		query,
		history.BookingID,
		history.OldStatus,
		history.NewStatus,
		history.ChangedBy,
		history.Reason,
		history.ChangedAt,
	).Scan(&history.ID, &history.CreatedAt)
}

// GetByBookingID возвращает всю историю статусов для бронирования
func (r *BookingStatusHistoryRepository) GetByBookingID(bookingID int) ([]domain.BookingStatusHistory, error) {
	const query = `
		SELECT 
			id, booking_id, old_status, new_status, 
			changed_by, changed_at, reason, created_at
		FROM booking_status_history
		WHERE booking_id = $1
		ORDER BY changed_at DESC, id DESC
	`

	rows, err := r.db.Query(query, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.BookingStatusHistory
	for rows.Next() {
		var h domain.BookingStatusHistory
		var oldStatus sql.NullString
		var reason sql.NullString

		err := rows.Scan(
			&h.ID,
			&h.BookingID,
			&oldStatus,
			&h.NewStatus,
			&h.ChangedBy,
			&h.ChangedAt,
			&reason,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Обрабатываем nullable поля
		if oldStatus.Valid {
			status := domain.BookingStatus(oldStatus.String)
			h.OldStatus = &status
		}

		if reason.Valid {
			h.Reason = &reason.String
		}

		history = append(history, h)
	}

	return history, rows.Err()
}

// RecordInitialStatus записывает начальный статус при создании бронирования
func (r *BookingStatusHistoryRepository) RecordInitialStatus(
	bookingID int,
	status domain.BookingStatus,
	userID int,
) error {
	history := &domain.BookingStatusHistory{
		BookingID: bookingID,
		OldStatus: nil, // Нет старого статуса при создании
		NewStatus: status,
		ChangedBy: userID,
		ChangedAt: time.Now(),
		Reason:    nil,
	}
	return r.Create(history)
}
