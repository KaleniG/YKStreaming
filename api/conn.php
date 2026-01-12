<?php
function getConn()
{
  $hostname = "localhost";
  $username = "babylon";
  $database = "ykstreaming";
  $password = "SofSvi_37";

  $dsn = "pgsql:host=$hostname;dbname=$database";
  $options = [
    PDO::ATTR_EMULATE_PREPARES   => false,
    PDO::ATTR_ERRMODE            => PDO::ERRMODE_EXCEPTION,
    PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
  ];

  try {
    $conn = new PDO($dsn, $username, $password, $options);
    return $conn;
  } catch (PDOException $e) {
    throw new Exception("Errore di connessione al database." . $e->getMessage());
  }
}
