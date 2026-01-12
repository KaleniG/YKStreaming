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

$response = ['success' => true, 'logged_in' => false];

if (isset($_SESSION['user_id'])) {
  $stmt = $pdo->prepare("SELECT id, email, name FROM users WHERE id = :id");
  $stmt->execute([':id' => $_SESSION['user_id']]);
  $user = $stmt->fetch(PDO::FETCH_ASSOC);
  if ($user) $response = ['success' => true, 'logged_in' => true, 'user' => $user];
} elseif (isset($_COOKIE['remember_token'])) {
  $user = getUserByRememberToken($_COOKIE['remember_token']);
  if ($user) {
    $_SESSION['user_id'] = $user['id'];
    $response = ['success' => true, 'logged_in' => true, 'user' => $user];
  }
}

header('Content-Type: application/json');
echo json_encode($response);
