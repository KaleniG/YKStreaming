<?php

// --------------------------------------------------
// 1) Get RTMP/HLS data
// --------------------------------------------------
$streamKey = $_POST['name'] ?? '';
$app       = $_POST['app'] ?? 'live';

// Ignore invalid / stop updates
if (!$streamKey || (isset($_POST['time']) && $_POST['time'] == 0)) {
  http_response_code(200);
  exit;
}

// --------------------------------------------------
// 2) Connect to PostgreSQL
// --------------------------------------------------
include "conn.php";
$pdo = getConn();

$stmt = $pdo->prepare("SELECT id, active FROM streams WHERE key = :key LIMIT 1");
$stmt->execute([':key' => $streamKey]);
$stream = $stmt->fetch(PDO::FETCH_ASSOC);

// --------------------------------------------------
// 3) Validate stream
// --------------------------------------------------
if (!$stream) {
  http_response_code(403);
  echo "Invalid stream key";
  exit;
}

if ($stream["active"] == false) {
  http_response_code(403);
  echo "Stream ended";
  exit;
}

// --------------------------------------------------
// 4) Screenshot + lock directories
// --------------------------------------------------
$saveDir = "/var/www/stream_screenshots";
$lockDir = $saveDir . "/.locks";

if (!is_dir($saveDir)) {
  mkdir($saveDir, 0777, true);
}
if (!is_dir($lockDir)) {
  mkdir($lockDir, 0777, true);
}

// --------------------------------------------------
// 5) Lock: once per ~minute per stream
// --------------------------------------------------
$lockFile = $lockDir . "/{$streamKey}.lock";

if (file_exists($lockFile) && time() - filemtime($lockFile) < 55) {
  http_response_code(200);
  exit;
}
touch($lockFile);

// --------------------------------------------------
// 6) Take screenshot from HLS (async + safe)
// --------------------------------------------------
// HLS path: adjust this to your actual HLS output folder
$hlsPath = "/var/www/hls/{$streamKey}.m3u8";

$screenshotPath = "$saveDir/{$streamKey}.jpg";

// Only attempt if HLS file exists
if (file_exists($hlsPath)) {
  $cmd = sprintf(
    'timeout 5 ffmpeg -y -loglevel error -i %s -frames:v 1 -q:v 3 %s > /dev/null 2>&1 &',
    escapeshellarg($hlsPath),
    escapeshellarg($screenshotPath)
  );

  exec($cmd);
}

// --------------------------------------------------
http_response_code(200);
echo "OK";
