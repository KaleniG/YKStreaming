<?php
session_start();

// CORS headers
$origin = "http://localhost:5173";
header("Access-Control-Allow-Origin: $origin");
header("Access-Control-Allow-Credentials: true");
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

include "conn.php";

function generateToken($length = 32)
{
  $bytes = random_bytes($length);
  return bin2hex($bytes);
}

function getUserByRememberToken(string $token)
{
  $pdo = getConn();
  $stmt = $pdo->prepare("SELECT id, email FROM users WHERE remember_token = :token");
  $stmt->execute([':token' => $token]);
  return $stmt->fetch(PDO::FETCH_ASSOC);
}

// Session restore from remember_token cookie
if (isset($_COOKIE['remember_token']) && !isset($_SESSION['user_id'])) {
  $user = getUserByRememberToken($_COOKIE['remember_token']);
  if ($user) {
    $_SESSION['user_id'] = $user['id'];
  }
}

if (!isset($_SESSION['user_id'])) {
  json_error("Not authenticated", 401);
}

$pdo = getConn();

// Only POST
if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
  json_error("Invalid request method", 405);
}

// ✅ Generate unique stream key
$new_stream_key = null;
do {
  $new_stream_key = generateToken();
  $stmt = $pdo->prepare("SELECT 1 FROM streams WHERE key = :new_stream_key");
  $stmt->execute([":new_stream_key" => $new_stream_key]);
} while ($stmt->fetchColumn());

// ✅ Get POST fields
$name = $_POST['name'] ?? null;
$is_vod = isset($_POST['is_vod']) ? (bool)$_POST['is_vod'] : null;

// Validate POST fields
if (!$name || !isset($is_vod)) {
  json_error("Missing fields");
}

// FILE thumbnail
$has_custom_thumbnail = 0;
$thumbnailPath = null;
$ext = null;

if (isset($_FILES['thumbnail']) && $_FILES['thumbnail']['error'] === 0) {
  $file = $_FILES['thumbnail'];
  $ext = pathinfo($file['name'], PATHINFO_EXTENSION);

  $uploadDir = "/var/www/thumbnails/";
  if (!is_dir($uploadDir)) mkdir($uploadDir, 0755, true);

  $thumbnailFileName = $new_stream_key . "." . $ext;
  $thumbnailPath = $uploadDir . $thumbnailFileName;

  if (!move_uploaded_file($file['tmp_name'], $thumbnailPath)) {
    json_error("Failed to move uploaded file");
  }
  $has_custom_thumbnail = 1;
}


// Insert into DB
$stmt = $pdo->prepare("
    INSERT INTO streams (key, user_id, name, has_custom_thumbnail, thumbnail_format, is_vod)
    VALUES (:new_key, :user_id, :name, :has_custom_thumbnail, :thumbnail_format, :is_vod)
");

if (
  $stmt->execute([
    ":new_key" => $new_stream_key,
    ":user_id" => $_SESSION["user_id"],
    ":name" => $name,
    ":has_custom_thumbnail" => $has_custom_thumbnail,
    ":thumbnail_format" => empty($ext) ? null : $ext,
    ":is_vod" => $is_vod ? 1 : 0
  ])
) {
  echo json_encode([
    "success" => true,
    "new_key" => $new_stream_key
  ]);
} else {
  json_error("Failed to insert stream into DB");
}
