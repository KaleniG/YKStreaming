CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  remember_token VARCHAR(64)
);

CREATE TABLE streams (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  key TEXT NOT NULL UNIQUE,
  has_custom_thumbnail BOOLEAN,
  is_vod BOOLEAN,
  is_active BOOLEAN DEFAULT FALSE,
  total_views INTEGER DEFAULT 0,
  live_viewers INTEGER DEFAULT 0,
  started_at TIMESTAMPTZ,
  ended_at TIMESTAMPTZ,
  user_id INT NOT NULL,
  CONSTRAINT fk_stream_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE views (
  id SERIAL PRIMARY KEY,
  is_watching BOOLEAN DEFAULT TRUE,
  guest_token TEXT,
  user_id INT,
  stream_id INT NOT NULL,
  CONSTRAINT fk_view_user_id FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_view_stream_id FOREIGN KEY (stream_id) REFERENCES streams(id)
);

CREATE UNIQUE INDEX uniq_user_stream ON views(user_id, stream_id)
WHERE user_id IS NOT NULL;

CREATE UNIQUE INDEX uniq_guest_stream ON views(guest_token, stream_id)
WHERE guest_token IS NOT NULL;

CREATE OR REPLACE FUNCTION start_stream(
  p_stream_key TEXT
)
RETURNS BOOLEAN AS $$
DECLARE
    rows_changed INT;
BEGIN
  UPDATE streams
  SET is_active = TRUE, started_at = NOW()
  WHERE key = p_stream_key
    AND ended_at IS NULL
    AND is_active IS DISTINCT FROM TRUE;

  GET DIAGNOSTICS rows_changed = ROW_COUNT;

  RETURN rows_changed > 0;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION view_stream_as_user(
    p_user_id INT,
    p_stream_key TEXT
)
RETURNS BOOLEAN AS $$
DECLARE
    v_stream_id INT;
    v_xmax xid;
BEGIN
    -- 1. Find the stream
    SELECT s.id
    INTO v_stream_id
    FROM streams s
    WHERE s.key = p_stream_key
      AND s.is_active = TRUE
      AND s.ended_at IS NULL
    LIMIT 1;

    IF v_stream_id IS NULL THEN
      RETURN FALSE;
    END IF;

    -- 2. Insert or update the view
    INSERT INTO views (user_id, stream_id, is_watching)
    VALUES (p_user_id, v_stream_id, TRUE)
    ON CONFLICT (user_id, stream_id) WHERE user_id IS NOT NULL
    DO UPDATE
    SET is_watching = TRUE
    WHERE views.is_watching IS DISTINCT FROM TRUE
    RETURNING xmax INTO v_xmax;

    -- 3. If xmax = 0, it was an insert
    IF v_xmax = 0 THEN
      UPDATE streams
      SET total_views = total_views + 1
      WHERE id = v_stream_id;
    END IF;

    -- 4. If less than 300 total_viewers just manually add the live_viewers
    IF (SELECT total_views FROM streams WHERE id = v_stream_id) <= 300 THEN
      UPDATE streams
      SET live_viewers = live_viewers + 1
      WHERE id = v_stream_id;
    END IF;

    -- 5. Always return true because insert/update happened
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION view_stream_as_guest(
    p_guest_token TEXT,
    p_stream_key TEXT
)
RETURNS BOOLEAN AS $$
DECLARE
    v_stream_id INT;
    v_xmax xid;
BEGIN
    -- 1. Find the stream
    SELECT s.id
    INTO v_stream_id
    FROM streams s
    WHERE s.key = p_stream_key
      AND s.is_active = TRUE
      AND s.ended_at IS NULL
    LIMIT 1;

    IF v_stream_id IS NULL THEN
      RETURN FALSE;
    END IF;

    -- 2. Insert or update the view
    INSERT INTO views (guest_token, stream_id, is_watching)
    VALUES (p_guest_token, v_stream_id, TRUE)
    ON CONFLICT (guest_token, stream_id) WHERE guest_token IS NOT NULL
    DO UPDATE
    SET is_watching = TRUE
    WHERE views.is_watching IS DISTINCT FROM TRUE
    RETURNING xmax INTO v_xmax;

    -- 3. If xmax = 0, it was an insert
    IF v_xmax = 0 THEN
      UPDATE streams
      SET total_views = total_views + 1
      WHERE id = v_stream_id;
    END IF;

    -- 4. If less than 300 total_viewers just manually add the live_viewers
    IF (SELECT total_views FROM streams WHERE id = v_stream_id) <= 300 THEN
      UPDATE streams
      SET live_viewers = live_viewers + 1
      WHERE id = v_stream_id;
    END IF;

    -- 5. Always return true because insert/update happened
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION unview_stream_as_user(
  p_user_id INT,
  p_stream_key TEXT
)
RETURNS BOOLEAN AS $$
DECLARE
  v_stream_id INT;
  v_updated BOOLEAN := FALSE;
BEGIN
  -- 1. Find the stream
  SELECT s.id
  INTO v_stream_id
  FROM streams s
  WHERE s.key = p_stream_key
    AND s.is_active = TRUE
    AND s.ended_at IS NULL
  LIMIT 1;

  -- If stream not found, return false
  IF v_stream_id IS NULL THEN
    RETURN FALSE;
  END IF;

  -- 2. Update the view
  UPDATE views AS v
  SET is_watching = FALSE
  WHERE v.stream_id = v_stream_id
    AND v.user_id = p_user_id
    AND v.is_watching IS DISTINCT FROM FALSE
  RETURNING TRUE INTO v_updated;

  -- 3. If less than 300 total_viewers just manually remove the live_viewers
  UPDATE streams
  SET live_viewers = GREATEST(live_viewers - 1, 0)
  WHERE id = v_stream_id
    AND total_views <= 300;

  RETURN COALESCE(v_updated, FALSE);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION unview_stream_as_guest(
  p_guest_token TEXT,
  p_stream_key TEXT
)
RETURNS BOOLEAN AS $$
DECLARE
  v_stream_id INT;
  v_updated BOOLEAN := FALSE;
BEGIN
  -- 1. Find the stream
  SELECT s.id
  INTO v_stream_id
  FROM streams s
  WHERE s.key = p_stream_key
    AND s.is_active = TRUE
    AND s.ended_at IS NULL
  LIMIT 1;

  -- If stream not found, return false
  IF v_stream_id IS NULL THEN
    RETURN FALSE;
  END IF;

  -- 2. Update the view
  UPDATE views AS v
  SET is_watching = FALSE
  WHERE v.stream_id = v_stream_id
    AND v.guest_token = p_guest_token
    AND v.is_watching IS DISTINCT FROM FALSE
  RETURNING TRUE INTO v_updated;

  -- 3. If less than 300 total_viewers just manually remove the live_viewers
  UPDATE streams
  SET live_viewers = GREATEST(live_viewers - 1, 0)
  WHERE id = v_stream_id
    AND total_views <= 300;

  RETURN COALESCE(v_updated, FALSE);
END;
$$ LANGUAGE plpgsql;