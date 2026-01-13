<?php
$streamKey = $_POST['name'] ?? '';

if (!$streamKey) {
  http_response_code(403);
  echo "No stream key provided";
  exit;
}

include __DIR__ . "/conn.php";
$pdo = getConn();

$stmt = $pdo->prepare("SELECT id, is_vod FROM streams WHERE key = :key LIMIT 1");
$stmt->execute([':key' => $streamKey]);
$stream = $stmt->fetch(PDO::FETCH_ASSOC);

if (!$stream) {
  http_response_code(403);
  echo "Stream not found";
  exit;
}

$path = "/var/www/recordings";
$files = glob("{$path}/{$streamKey}*.flv");
if (!$files) {
  return;
}
usort($files, fn($a, $b) => filemtime($b) - filemtime($a));
$latestFLV = "/var/www/recordings/" . basename($files[0]);

if (!$stream["is_vod"]) {
  unlink($latestFLV);
} else {
  $input = escapeshellarg($latestFLV);
  $output = escapeshellarg($path . "/" . $streamKey . ".mp4");

  if (file_exists($latestFLV)) {
    // try fast copy first
    $cmd = "timeout 5 ffmpeg -y -i $input -c copy $output 2>&1";
    exec($cmd, $out, $code);

    if ($code !== 0) {
      // fallback: re-encode
      $cmd = "timeout 5 ffmpeg -y -i $input -c:v libx264 -c:a aac $output 2>&1";
      exec($cmd, $out, $code);
    }
  }
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
