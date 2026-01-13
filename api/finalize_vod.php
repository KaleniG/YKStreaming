<?php
$flvPath = $_POST['path'] ?? '';

if (!$flvPath) {
  http_response_code(403);
  echo "No vod path prvided";
  exit;
}

$mp4Path = pathinfo($flvPath, PATHINFO_DIRNAME) . '/' . pathinfo($flvPath, PATHINFO_FILENAME) . ".mp4";

$input = escapeshellarg($flvPath);
$output = escapeshellarg($mp4Path);

if (file_exists($flvPath)) {
  // try fast copy first
  $cmd = "timeout 5 ffmpeg -y -i $input -c copy $output 2>&1";
  exec($cmd);

  if ($code !== 0) {
    // fallback: re-encode
    $cmd = "timeout 5 ffmpeg -y -i $input -c:v libx264 -c:a aac $output 2>&1";
    exec($cmd);
  }
}

http_response_code(200);
echo "OK";
exit;
