CREATE TABLE IF NOT EXISTS spaces (
                                      id SERIAL PRIMARY KEY,
                                      owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                      title VARCHAR(255) NOT NULL,
                                      description TEXT,
                                      area_m2 NUMERIC(10,2) NOT NULL,
                                      price INTEGER NOT NULL,
                                      phone VARCHAR(20),
                                      is_active BOOLEAN NOT NULL DEFAULT TRUE,
                                      created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_spaces_owner_id ON spaces(owner_id);
CREATE INDEX IF NOT EXISTS idx_spaces_active ON spaces(is_active);
