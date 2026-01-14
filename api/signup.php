<?php
header("Access-Control-Allow-Origin: http://localhost:5173");
header("Access-Control-Allow-Credentials: true");
header("Access-Control-Allow-Headers: Content-Type");
header("Access-Control-Allow-Methods: GET, POST, OPTIONS");
header("Content-Type: application/json; charset=UTF-8");

session_start();

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

function addUser($name, $email, $password)
{
  $password_hash = password_hash($password, PASSWORD_DEFAULT);
  $pdo = getConn();
  $stmt = $pdo->prepare(
    "INSERT INTO users (name, email, password_hash) VALUES (:name, :email, :password_hash) RETURNING id"
  );
  $stmt->execute([
    ':name' => $name,
    ':email' => $email,
    ':password_hash' => $password_hash
  ]);
  return $stmt->fetch(PDO::FETCH_ASSOC);
}

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
  json_error("Invalid request method", 405);
}

$input = json_decode(file_get_contents("php://input"), true);
$email = $input['email'] ?? null;
$username = $input['name'] ?? null;
$password = $input['password'] ?? null;

if (!$email || !$password || !$username) {
  json_error("Missing fields");
}

$stmt = $pdo->prepare("SELECT * FROM users WHERE email = :email OR name = :name");
$stmt->execute([':email' => $email, ':name' => $username]);
$user = $stmt->fetch(PDO::FETCH_ASSOC);
if ($user) {
  echo json_encode(["success" => false, "existing" => true]);
  exit();
} else {
  $user = addUser($username, $email, $password);
}

if ($user) {
  $_SESSION['user_id'] = $user['id'];
  echo json_encode(["success" => true, "user" => $user]);
  exit();
} else {
  json_error("Failed to register");
}
