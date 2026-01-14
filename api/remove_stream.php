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
$stmt = $pdo->prepare("SELECT id, has_custom_thumbnail, active, is_vod FROM streams WHERE key = :key_to_remove AND user_id = :user_id");
$stmt->execute([
  ":key_to_remove" => $key,
  ":user_id" => $_SESSION['user_id']
]);

$streamToRemove = $stmt->fetch(PDO::FETCH_ASSOC);


// Delete stream
if ($streamToRemove) {
  if ($streamToRemove["active"]) {
    if ($streamToRemove["is_vod"] && file_exists("/var/www/recordings/{$key}.flv")) {
      call_api("http://localhost:8080/control/record/stop?app=live&name={$key}&rec=vod");
    }
    call_api("http://localhost:8080/control/drop/publisher?app=live&name={$key}");
  }

  $stmt = $pdo->prepare("DELETE FROM views WHERE stream_id = :stream_id");
  $stmt->execute([
    ":stream_id" => $streamToRemove["id"]
  ]);

  $stmt = $pdo->prepare("DELETE FROM streams WHERE id = :stream_id AND user_id = :user_id");
  $stmt->execute([
    ":stream_id" => $streamToRemove["id"],
    ":user_id" => $_SESSION['user_id']
  ]);

  if ($streamToRemove["has_custom_thumbnail"]) {
    $thumbnail = "/var/www/thumbnails/{$key}";

    foreach (glob($thumbnail . '.*') as $file) {
      if (is_file($file)) {
        unlink($file);
      }
    }
  }

  $recording = "/var/www/recordings/{$key}";
  foreach (glob($recording . '.*') as $file) {
    if (is_file($file)) {
      unlink($file);
    }
  }

  $screenshot = "/var/www/stream_screenshots/{$key}.jpg";
  $screenshotLock = "/var/www/stream_screenshots/.locks/{$key}.lock";

  if (file_exists($screenshot)) {
    unlink($screenshot);
  }

  // Remove lock file
  if (file_exists($screenshotLock)) {
    unlink($screenshotLock);
  }

  echo json_encode(["success" => true]);
} else {
  json_error("Stream key not found");
}
