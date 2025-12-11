package worker

import (
    "context"
    "log"
    "strconv"
    "time"

    "SpaceBookProject/internal/domain"
    "SpaceBookProject/internal/services"
)

type BookingEventWorker struct {
    events        <-chan domain.BookingEvent
    notifications services.NotificationService
}

func NewBookingEventWorker(
    events <-chan domain.BookingEvent,
    notifications services.NotificationService,
) *BookingEventWorker {
    return &BookingEventWorker{
        events:        events,
        notifications: notifications,
    }
}

func (w *BookingEventWorker) Run(ctx context.Context) {
    log.Println("[worker] booking event worker started")
    defer log.Println("[worker] booking event worker stopped")

    for {
        select {
        case <-ctx.Done():
            return

        case evt := <-w.events:
            msg := buildNotificationMessage(evt)

            n := &domain.Notification{
                UserID:  evt.TenantID,
                Type:    string(evt.Type),
                Message: msg,
            }

            // через сервис
            if err := w.notifications.CreateNotification(n); err != nil {
                log.Printf("[worker] failed create notification: %v", err)
            } else {
                log.Printf("[worker] notification created id=%d for user=%d", n.ID, n.UserID)
            }

            log.Printf(
                "[worker] event=%s booking_id=%d space_id=%d tenant_id=%d at=%s",
                evt.Type,
                evt.BookingID,
                evt.SpaceID,
                evt.TenantID,
                evt.At.Format(time.RFC3339),
            )
        }
    }
}

func buildNotificationMessage(evt domain.BookingEvent) string {
    switch evt.Type {
    case domain.BookingEventCreated:
        return "Your booking #" + strconv.Itoa(evt.BookingID) + " was created."
    case domain.BookingEventApproved:
        return "Your booking #" + strconv.Itoa(evt.BookingID) + " was approved."
    case domain.BookingEventRejected:
        return "Your booking #" + strconv.Itoa(evt.BookingID) + " was rejected."
    case domain.BookingEventCancelled:
        return "Your booking #" + strconv.Itoa(evt.BookingID) + " was cancelled."
    default:
        return "Booking event: " + string(evt.Type)
    }
}
