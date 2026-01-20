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

    -- 4. Always return true because insert/update happened
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

    -- 4. Always return true because insert/update happened
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

  RETURN COALESCE(v_updated, FALSE);
END;
$$ LANGUAGE plpgsql;