-- name: GetUserStreams :many
SELECT
  s.name,
  s.key,
  s.is_active,
  s.ended_at,
  s.started_at,
  s.total_views,
  s.is_vod,
  COUNT(v.id) FILTER (WHERE v.is_watching = TRUE) AS live_viewers
FROM streams s
LEFT JOIN views v ON v.stream_id = s.id
WHERE s.user_id = sqlc.arg(user_id)
GROUP BY s.id
ORDER BY s.id DESC;

-- name: AddStream :one
INSERT INTO streams 
(
  key, 
  user_id, 
  name, 
  has_custom_thumbnail, 
  is_vod
)
VALUES 
(
  sqlc.arg(key), 
  sqlc.arg(user_id), 
  sqlc.arg(name), 
  sqlc.arg(has_custom_thumbnail), 
  sqlc.arg(is_vod)
)
ON CONFLICT (key) DO NOTHING 
RETURNING 1;

-- name: GetStreamRemovalDataByKey :one
SELECT id, has_custom_thumbnail, is_active, is_vod 
FROM streams 
WHERE key = sqlc.arg(key)
LIMIT 1;

-- name: RemoveStream :exec
DELETE FROM streams 
WHERE id = sqlc.arg(stream_id);

-- name: GetStreamStopDataByKey :one
SELECT id, is_active, is_vod 
FROM streams 
WHERE key = sqlc.arg(key)
LIMIT 1;

-- name: StopStream :exec
UPDATE streams 
SET is_active = FALSE 
WHERE id = sqlc.arg(stream_id);

-- name: GetPublicStreams :many
SELECT
  u.name AS streamer_name,
  s.key,
  s.name,
  s.has_custom_thumbnail,
  (s.is_active = TRUE AND s.ended_at IS NULL) AS is_live,
  s.is_vod,
  COUNT(v.id) FILTER ( WHERE v.is_watching = TRUE) AS live_viewers
FROM streams s
JOIN users u ON s.user_id = u.id
LEFT JOIN views v ON v.stream_id = s.id
WHERE (s.is_active = TRUE AND s.ended_at IS NULL) OR (s.is_vod = TRUE AND s.ended_at IS NOT NULL)
GROUP BY s.id, u.name;

-- name: GetStreamStatus :one
SELECT 
  is_vod,
  (is_active = TRUE AND ended_at IS NULL) AS is_live
FROM streams 
WHERE key = sqlc.arg(key) 
LIMIT 1;

-- name: StartStream :one
SELECT start_stream(sqlc.arg(stream_key))
AS started;

-- name: EndStream :one
UPDATE streams 
SET 
  ended_at = NOW(), 
  is_active = FALSE 
WHERE key = sqlc.arg(key)
RETURNING is_vod;