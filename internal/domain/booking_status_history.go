package domain

import "time"

type BookingStatusHistory struct {
	ID        int       `json:"id"`
	BookingID int       `json:"booking_id"`
	OldStatus *string   `json:"old_status"`
	NewStatus string    `json:"new_status"`
	ChangedBy int       `json:"changed_by"`
	ChangedAt time.Time `json:"changed_at"`
	CreatedAt time.Time `json:"created_at"`
}
