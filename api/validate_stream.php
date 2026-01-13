<?php
// get the stream key from Nginx
$streamKey = $_POST['name'] ?? '';

if (!$streamKey) {
  http_response_code(403);
  echo "No stream key provided";
  exit;
}

// Connect to PostgreSQL
include "conn.php";
$pdo = getConn();


// Check if the stream key exists in your database
$stmt = $pdo->prepare("SELECT id FROM streams WHERE key = :key AND active = FALSE AND ended_at IS NULL LIMIT 1");
$stmt->execute([':key' => $streamKey]);
$stream = $stmt->fetch(PDO::FETCH_ASSOC);

if ($stream) {
  // Key is valid → allow streaming
  $stmt = $pdo->prepare("UPDATE streams SET active = TRUE, started_at = NOW() WHERE id = :id");
  if ($stmt->execute([
    ':id' => $stream["id"]
  ])) {
    http_response_code(200);
    echo "OK";
    exit;
  }
} else {
  // Invalid key → reject stream
  http_response_code(403);
  echo "Invalid stream key";
  exit;
}
