<?php
session_start();

$origin = "http://localhost:5173";
header("Access-Control-Allow-Origin: $origin");
header("Access-Control-Allow-Credentials: true");
header("Content-Type: application/json; charset=UTF-8");
header("Access-Control-Allow-Methods: GET, POST, OPTIONS");
header("Access-Control-Allow-Headers: Content-Type, Authorization");

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
  http_response_code(200);
  exit();
}

include "conn.php";

$pdo = getConn();

$stmt = $pdo->prepare(
  "SELECT
  u.name AS streamer_name,
  s.key,
  s.name,
  s.has_custom_thumbnail AS uses_thumbnail,
  s.thumbnail_format,
  (s.active = TRUE AND s.ended_at IS NULL) AS is_live,
  s.is_vod,
  COUNT(v.id) FILTER ( WHERE v.watching = TRUE) AS live_viewers
  FROM streams s
  JOIN users u ON s.user_id = u.id
  LEFT JOIN views v ON v.stream_id = s.id
  WHERE (s.active = TRUE AND s.ended_at IS NULL) OR (s.is_vod = TRUE AND s.ended_at IS NOT NULL)
  GROUP BY s.id, u.name"
); // LAWLESSNESS GOING ON HERE

$stmt->execute();
$streams = $stmt->fetchAll();

echo json_encode(["success" => true, "streams" => $streams]);
