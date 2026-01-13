BEGIN;

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  remember_token TEXT
);

DROP TABLE IF EXISTS streams CASCADE;
CREATE TABLE streams (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  key TEXT NOT NULL UNIQUE,
  has_custom_thumbnail BOOLEAN,
  thumbnail_format TEXT,
  is_vod BOOLEAN,
  active BOOLEAN DEFAULT FALSE,
  views INTEGER DEFAULT 0,
  started_at TIMESTAMPTZ,
  ended_at TIMESTAMPTZ,
  user_id INT NOT NULL,
  
  CONSTRAINT fk_stream_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

DROP TABLE IF EXISTS views CASCADE;
CREATE TABLE views (
  id SERIAL PRIMARY KEY,
  
  watching BOOLEAN DEFAULT TRUE,

  guest_token TEXT,
  user_id INT,
  stream_id INT NOT NULL,
  
  CONSTRAINT fk_view_user_id FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_view_stream_id FOREIGN KEY (stream_id) REFERENCES streams(id)
);

-- logged-in users
CREATE UNIQUE INDEX uniq_user_stream
ON views(user_id, stream_id)
WHERE user_id IS NOT NULL;

-- guests
CREATE UNIQUE INDEX uniq_guest_stream
ON views(guest_token, stream_id)
WHERE guest_token IS NOT NULL;

COMMIT;