CREATE TABLE IF NOT EXISTS booking_status_history (
  id          SERIAL PRIMARY KEY,
  booking_id  INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
  old_status  VARCHAR(50),
  new_status  VARCHAR(50) NOT NULL,
  changed_by  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  changed_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_bsh_booking_id  ON booking_status_history(booking_id);
CREATE INDEX IF NOT EXISTS idx_bsh_changed_at  ON booking_status_history(changed_at);
