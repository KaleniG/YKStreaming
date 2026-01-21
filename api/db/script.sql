DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  remember_token VARCHAR(64)
);

CREATE INDEX idx_email ON users(email);
CREATE INDEX idx_remember_token ON users(remember_token);

DROP TABLE IF EXISTS streams CASCADE;
CREATE TABLE streams (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  key TEXT NOT NULL UNIQUE,
  has_custom_thumbnail BOOLEAN,
  is_vod BOOLEAN,
  is_active BOOLEAN DEFAULT FALSE,
  total_views INTEGER DEFAULT 0,
  started_at TIMESTAMPTZ,
  ended_at TIMESTAMPTZ,
  user_id INT NOT NULL,
  CONSTRAINT fk_stream_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_key ON streams(key);

DROP TABLE IF EXISTS views CASCADE;
CREATE TABLE views (
  id SERIAL PRIMARY KEY,
  is_watching BOOLEAN DEFAULT TRUE,
  guest_token TEXT,
  user_id INT,
  stream_id INT NOT NULL,
  CONSTRAINT fk_view_user_id FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_view_stream_id FOREIGN KEY (stream_id) REFERENCES streams(id)
);

CREATE UNIQUE INDEX uniq_user_stream ON views (user_id, stream_id)
WHERE user_id IS NOT NULL;

CREATE UNIQUE INDEX uniq_guest_stream ON views (guest_token, stream_id)
WHERE guest_token IS NOT NULL;