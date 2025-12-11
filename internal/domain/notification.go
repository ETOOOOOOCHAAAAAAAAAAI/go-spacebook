package domain

import "time"

type Notification struct {
    ID        int       `json:"id" db:"id"`
    UserID    int       `json:"user_id" db:"user_id"`
    Type      string    `json:"type" db:"type"`
    Message   string    `json:"message" db:"message"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
