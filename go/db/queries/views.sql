-- name: ViewStreamAsUser :one
SELECT view_stream_as_user(sqlc.arg(user_id), sqlc.arg(stream_key)) 
AS stream_found;

-- name: ViewStreamAsGuest :one
SELECT view_stream_as_guest(sqlc.arg(guest_token), sqlc.arg(stream_key))
AS stream_found;

-- name: UnviewStreamAsUser :exec
SELECT unview_stream_as_user(sqlc.arg(user_id), sqlc.arg(stream_key))
AS stream_found;

-- name: UnviewStreamAsGuest :exec
SELECT unview_stream_as_guest(sqlc.arg(guest_token), sqlc.arg(stream_key))
AS stream_found;

-- name: RemoveStreamViewers :exec
DELETE FROM views
WHERE stream_id = sqlc.arg(stream_id);