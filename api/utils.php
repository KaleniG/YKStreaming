<?php

function call_api(string $url): bool
{
  $ch = curl_init();
  curl_setopt($ch, CURLOPT_URL, $url);
  curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
  curl_exec($ch);
  return (curl_getinfo($ch, CURLINFO_HTTP_CODE) == 200);
}
