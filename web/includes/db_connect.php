<?php
    include_once(__DIR__."/../../includes/db_connect.php");
    $monitorDB = new PDO('mysql:host='.HOST.";dbname=".DATABASE_MONITOR, USER, PASSWORD);
?>