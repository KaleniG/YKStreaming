<?php
session_start();

header("Access-Control-Allow-Origin: http://localhost:5173");
header("Access-Control-Allow-Credentials: true");
header("Access-Control-Allow-Headers: Content-Type");
header("Access-Control-Allow-Methods: GET, POST, OPTIONS");
header("Content-Type: application/json; charset=UTF-8");

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

function generateRememberToken($userId)
{
  $token = bin2hex(random_bytes(32));
  saveTokenToDatabase($userId, $token);
  return $token;
}

function saveTokenToDatabase($userId, $token)
{
  $pdo = getConn();
  $stmt = $pdo->prepare("UPDATE users SET remember_token = :token WHERE id = :id");
  $stmt->execute([
    ':token' => $token,
    ':id' => $userId
  ]);
}

function getUserByEmail($email)
{
  $pdo = getConn();
  $stmt = $pdo->prepare("SELECT id, email, password_hash FROM users WHERE email = :email");
  $stmt->execute([':email' => $email]);
  return $stmt->fetch(PDO::FETCH_ASSOC);
}

// ❗ Now POST-only check AFTER OPTIONS handling
if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
  json_error("Invalid request method", 405);
}

$input = json_decode(file_get_contents("php://input"), true);
$email = $input['email'] ?? null;
$password = $input['password'] ?? null;
$remember_me = $input['remember_me'] ?? false;

$user = getUserByEmail($email);

if ($user && password_verify($password, $user['password_hash'])) {

  $_SESSION['user_id'] = $user['id'];

  if ($remember_me) {
    // 🔥 FIX: generate token
    $token = generateRememberToken($user['id']);

    setcookie(
      "remember_token",
      $token,
      [
        'expires' => time() + 60 * 60 * 24 * 30,
        'path' => '/',
        'secure' => false,   // true when HTTPS
        'httponly' => true,
        'samesite' => 'Lax'
      ]
    );
  }

  echo json_encode(['success' => true, 'user' => $user]);
} else {
  json_error('Invalid credentials');
}
