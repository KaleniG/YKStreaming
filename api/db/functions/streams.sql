CREATE OR REPLACE FUNCTION start_stream(
  p_stream_key TEXT
)
RETURNS BOOLEAN AS $$
DECLARE
    rows_changed INT;
BEGIN
  UPDATE streams
  SET is_active = TRUE
  WHERE key = p_stream_key
    AND ended_at IS NULL
    AND is_active IS DISTINCT FROM TRUE;

  GET DIAGNOSTICS rows_changed = ROW_COUNT;

  RETURN rows_changed > 0;
END;
$$ LANGUAGE plpgsql;