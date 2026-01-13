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
function getUserByRememberToken(string $token)
{
  $pdo = getConn();
  $stmt = $pdo->prepare("SELECT id, email FROM users WHERE remember_token = :token");
  $stmt->execute([':token' => $token]);
  return $stmt->fetch(PDO::FETCH_ASSOC);
}

if (isset($_COOKIE['remember_token']) && !isset($_SESSION['user_id'])) {
  $user = getUserByRememberToken($_COOKIE['remember_token']);
  if ($user) {
    $_SESSION['user_id'] = $user['id'];
  }
}

if (!isset($_SESSION['user_id'])) {
  echo json_encode([]);
  exit();
}
$pdo = getConn();

$stmt = $pdo->prepare(
  "SELECT
  s.name,
  s.key,
  s.active,
  s.ended_at,
  s.views,
  s.is_vod,
  COUNT(v.id) FILTER (WHERE v.watching = TRUE) AS live_viewers
  FROM streams s
  LEFT JOIN views v ON v.stream_id = s.id
  WHERE s.user_id = :user_id
  GROUP BY s.id
  ORDER BY s.id DESC"
);

$stmt->execute([":user_id" => $_SESSION["user_id"]]);
$streams = $stmt->fetchAll(PDO::FETCH_ASSOC);

echo json_encode(["success" => true, "streams" => $streams]);
