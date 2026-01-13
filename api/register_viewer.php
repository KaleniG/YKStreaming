<?php
session_start();

$origin = "http://localhost:5173"; // React app
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
$pdo = getConn();

function generateToken($length = 32)
{
  $bytes = random_bytes($length);
  return bin2hex($bytes);
}

function json_error($msg, $code = 400)
{
  http_response_code($code);
  echo json_encode(['success' => false, 'error' => $msg]);
  exit();
}

function getUserByRememberToken(string $token)
{
  $pdo = getConn();
  $stmt = $pdo->prepare("SELECT id, email, name FROM users WHERE remember_token = :token");
  $stmt->execute([':token' => $token]);
  return $stmt->fetch(PDO::FETCH_ASSOC);
}

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
  json_error("Invalid request method", 405);
}

$input = json_decode(file_get_contents("php://input"), true);
$streamKey = $input['key'] ?? null;

if (!$streamKey) {
  json_error("Missing 'key' field");
}

if (isset($_SESSION['user_id'])) {
  $stmt = $pdo->prepare("SELECT id, email, name FROM users WHERE id = :id");
  $stmt->execute([':id' => $_SESSION['user_id']]);
  $user = $stmt->fetch(PDO::FETCH_ASSOC);
} elseif (isset($_COOKIE['remember_token'])) {
  $user = getUserByRememberToken($_COOKIE['remember_token']);
  if ($user) {
    $_SESSION['user_id'] = $user['id'];
  }
} elseif (!isset($_COOKIE["guest_token"])) {
  setcookie("guest_token", generateToken(), time() + 262746000, "/", "", false, true);
}

$stmt = $pdo->prepare("SELECT id, views FROM streams WHERE key = :stream_key");
$stmt->execute([':stream_key' => $streamKey]);
$stream = $stmt->fetch(PDO::FETCH_ASSOC);

if (!$stream) json_error("Missing stream");

if (isset($_SESSION['user_id'])) {
  $stmt = $pdo->prepare("SELECT watching FROM views WHERE user_id = :user_id AND stream_id = :stream_id LIMIT 1");
  if (!$stmt->execute([
    ':user_id' => $_SESSION['user_id'],
    ':stream_id' => $stream["id"]
  ])) {
    json_error("Failed to check views");
  }

  if ($stmt->rowCount() > 0) {
    if ($stmt->fetch(PDO::FETCH_ASSOC)["watching"] == false) {
      $stmt = $pdo->prepare("UPDATE views SET watching = TRUE WHERE user_id = :user_id AND stream_id = :stream_id");
      if (!$stmt->execute([
        ':user_id' => $_SESSION['user_id'],
        ':stream_id' => $stream["id"]
      ])) {
        json_error("Failed to save a view");
      }
    }
    header('Content-Type: application/json');
    echo json_encode(["success" => true]);
    exit;
  } else {
    $stmt = $pdo->prepare("INSERT INTO views (user_id, stream_id) VALUES (:user_id, :stream_id)");
    if (!$stmt->execute([
      ':user_id' => $_SESSION['user_id'],
      ':stream_id' => $stream["id"]
    ])) {
      json_error("Failed to save a view");
    }
  }
} else if (isset($_COOKIE["guest_token"])) {
  $stmt = $pdo->prepare("SELECT watching FROM views WHERE guest_token = :guest_token AND stream_id = :stream_id LIMIT 1");
  if (!$stmt->execute([
    ':guest_token' => $_COOKIE["guest_token"],
    ':stream_id' => $stream["id"]
  ])) {
    json_error("Failed to check views");
  }

  if ($stmt->rowCount() > 0) {
    if ($stmt->fetch(PDO::FETCH_ASSOC)["watching"] == true) {
      $stmt = $pdo->prepare("UPDATE views SET watching = TRUE WHERE guest_token = :guest_token AND stream_id = :stream_id");
      if (!$stmt->execute([
        ':guest_token' =>  $_COOKIE["guest_token"],
        ':stream_id' => $stream["id"]
      ])) {
        json_error("Failed to save a guest view");
      }
    }
    header('Content-Type: application/json');
    echo json_encode(["success" => true]);
    exit;
  } else {
    $stmt = $pdo->prepare("INSERT INTO views (guest_token, stream_id) VALUES (:guest_token, :stream_id)");
    if (!$stmt->execute([
      ':guest_token' => $_COOKIE["guest_token"],
      ':stream_id' => $stream["id"]
    ])) {
      json_error("Failed to save a guest view");
    }
  }
} else {
  json_error("How");
}

$stmt = $pdo->prepare("UPDATE streams SET views = :new_views WHERE id = :stream_id");
if (!$stmt->execute([
  ':new_views' => $stream["views"] + 1,
  ':stream_id' => $stream["id"]
])) {
  json_error("Failed to save a view");
}

header('Content-Type: application/json');
echo json_encode(["success" => true]);
