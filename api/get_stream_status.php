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

function json_error($msg, $code = 400)
{
  http_response_code($code);
  echo json_encode(['success' => false, 'error' => $msg]);
  exit();
}

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
  json_error("Invalid request method", 405);
}

$input = json_decode(file_get_contents("php://input"), true);
$streamKey = $input['stream_key'] ?? null;

if (!$streamKey) {
  json_error("Missing 'stream_key' field");
}

include "conn.php";
$pdo = getConn();

$stmt = $pdo->prepare("SELECT (active = TRUE AND ended_at IS NULL) AS is_live, is_vod FROM streams WHERE key = :stream_key");
$stmt->execute([":stream_key" => $streamKey]);

$result = $stmt->fetch(PDO::FETCH_ASSOC);

if ($result) {
  echo json_encode(["success" => true, "exists" => true, "is_live" => $result["is_live"], "is_vod" => $result["is_vod"]]);
} else {
  echo json_encode(["success" => true, "exists" => false]);
}
