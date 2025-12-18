-- Migration: Add booking status history table
CREATE TABLE IF NOT EXISTS booking_status_history (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    old_status VARCHAR(20) CHECK (old_status IN ('pending', 'approved', 'rejected', 'cancelled')),
    new_status VARCHAR(20) NOT NULL CHECK (new_status IN ('pending', 'approved', 'rejected', 'cancelled')),
    changed_by INTEGER NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    changed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_booking_status_history_booking_id ON booking_status_history(booking_id);
CREATE INDEX IF NOT EXISTS idx_booking_status_history_changed_at ON booking_status_history(changed_at DESC);
CREATE INDEX IF NOT EXISTS idx_booking_status_history_changed_by ON booking_status_history(changed_by);