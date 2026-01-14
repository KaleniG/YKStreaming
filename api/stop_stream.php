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

include __DIR__ . "/conn.php";
include __DIR__ . "/utils.php";

function json_error($message, $code = 400)
{
  http_response_code($code);
  echo json_encode(["success" => false, "error" => $message]);
  exit;
}

function getUserByRememberToken(string $token)
{
  $pdo = getConn();
  $stmt = $pdo->prepare("SELECT id, email FROM users WHERE remember_token = :token");
  $stmt->execute(['token' => $token]);
  return $stmt->fetch(PDO::FETCH_ASSOC);
}

// Auto-login via remember_token
if (!isset($_SESSION['user_id']) && !empty($_COOKIE['remember_token'])) {
  $user = getUserByRememberToken($_COOKIE['remember_token']);
  if ($user) {
    $_SESSION['user_id'] = $user['id'];
  }
}

if (!isset($_SESSION['user_id'])) {
  json_error("Unauthorized", 401);
}

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
  json_error("Invalid request method", 405);
}

// Parse JSON input
$input = json_decode(file_get_contents("php://input"), true);
$key = $input['key'] ?? null;

if (!$key) {
  json_error("Missing 'key' field");
}

$pdo = getConn();
$stmt = $pdo->prepare("SELECT id, active, is_vod FROM streams WHERE key = :key_to_stop AND user_id = :user_id");
$stmt->execute([
  ":key_to_stop" => $key,
  ":user_id" => $_SESSION['user_id']
]);

$streamToStop = $stmt->fetch(PDO::FETCH_ASSOC);

if ($streamToStop) {
  if ($streamToStop["active"]) {
    if ($streamToStop["is_vod"] && file_exists("/var/www/recordings/{$key}.flv")) {
      call_api("http://localhost:8080/control/record/stop?app=live&name={$key}&rec=vod");
    }
    call_api("http://localhost:8080/control/drop/publisher?app=live&name={$key}");
  }

  // Stop stream
  $pdo = getConn();
  $stmt = $pdo->prepare("UPDATE streams SET active = FALSE WHERE key = :stream_id AND user_id = :user_id");
  if ($stmt->execute([
    ":stream_id" => $streamToStop["id"],
    ":user_id" => $_SESSION['user_id']
  ])) {
    echo json_encode(["success" => true]);
  } else {
    json_error("Error");
  }
} else {
  json_error("Error");
}
