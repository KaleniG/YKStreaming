<?php
$streamKey = $_POST['name'] ?? '';

if (!$streamKey) {
  http_response_code(403);
  echo "No stream key provided";
  exit;
}

include __DIR__ . "/conn.php";
include __DIR__ . "/utils.php";
$pdo = getConn();

$stmt = $pdo->prepare("SELECT id, is_vod FROM streams WHERE key = :key LIMIT 1");
$stmt->execute([':key' => $streamKey]);
$stream = $stmt->fetch(PDO::FETCH_ASSOC);

if (!$stream) {
  http_response_code(403);
  echo "Stream not found";
  exit;
}

if ($stream["is_vod"] && file_exists("/var/www/recordings/{$streamKey}.flv")) {
  call_api("http://localhost:8080/control/record/stop?app=live&name={$streamKey}&rec=vod");
}

$stmt = $pdo->prepare("DELETE FROM views WHERE stream_id = :stream_id");
$stmt->execute([
  ":stream_id" => $stream["id"]
]);

$stmt = $pdo->prepare("UPDATE streams SET ended_at = NOW(), active = FALSE WHERE id = :id");
if ($stmt->execute([':id' => $stream["id"]])) {
  http_response_code(200);
  echo "OK";
  exit;
}

http_response_code(500);
echo "Something went wrong";
