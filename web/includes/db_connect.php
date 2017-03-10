<?php
include_once 'psl-config.php';   // As functions.php is not included
$mysqli = new mysqli(HOST, USER, PASSWORD, DATABASE);

$pdo = new PDO('mysql:host=localhost;dbname=WebServer', 'root', 'abc123');

$monitorDB = new PDO('mysql:host='.HOST.";dbname=".DATABASE_MONITOR, USER, PASSWORD);

?>