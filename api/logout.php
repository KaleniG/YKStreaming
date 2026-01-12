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

include "conn.php";
$pdo = getConn();

function json_error($msg, $code = 400)
{
  http_response_code($code);
  echo json_encode(['success' => false, 'error' => $msg]);
  exit();
}

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
  json_error("Invalid request method", 405);
}

if (isset($_SESSION['user_id'])) {
  $stmt = $pdo->prepare("UPDATE users SET remember_token = NULL WHERE id = :id");
  $stmt->execute([':id' => $_SESSION['user_id']]);
} elseif (isset($_COOKIE['remember_token'])) {
  $stmt = $pdo->prepare("UPDATE users SET remember_token = NULL WHERE remember_token = :token");
  $stmt->execute([':token' => $_COOKIE['remember_token']]);
}

setcookie(
  'remember_token',
  '',
  [
    'expires'  => time() - 3600,
    'path'     => '/',
    'secure'   => true,
    'httponly' => true,
    'samesite' => 'Lax',
  ]
);
unset($_COOKIE['remember_token']);
unset($_SESSION['user_id']);

header('Content-Type: application/json');
echo json_encode(['success' => true]);
