package repository

import (
    "SpaceBookProject/internal/domain"
    "database/sql"
)

type notificationRepoSQL struct {
    db *sql.DB
}

func NewNotificationRepository(db *sql.DB) NotificationRepository {
    return &notificationRepoSQL{db: db}
}

func (r *notificationRepoSQL) Create(n *domain.Notification) error {
    q := `INSERT INTO notifications (user_id, type, message) VALUES ($1, $2, $3) RETURNING id, created_at`
    return r.db.QueryRow(q, n.UserID, n.Type, n.Message).Scan(&n.ID, &n.CreatedAt)
}

func (r *notificationRepoSQL) ListByUser(userID int) ([]domain.Notification, error) {
    q := `SELECT id, user_id, type, message, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`
    rows, err := r.db.Query(q, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var res []domain.Notification
    for rows.Next() {
        var n domain.Notification
        if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Message, &n.CreatedAt); err != nil {
            return nil, err
        }
        res = append(res, n)
    }
    return res, rows.Err()
}
