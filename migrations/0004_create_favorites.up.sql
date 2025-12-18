CREATE TABLE IF NOT EXISTS favorites (
  user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  space_id    INTEGER NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, space_id)
);

CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);