<?php
include_once 'includes/db_connect.php';
include_once 'includes/functions.php';

sec_session_start();
 
if (!login_check($mysqli)) {
    header('Location: index.php');
}
?>

<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <title>MSTT Monitor Home</title>
        <link rel="stylesheet" href="css/bootstrap.min.css" />
        <link rel="stylesheet" href="css/main.css" />
    </head>
    <body>
        <nav class="navbar navbar-inverse navbar-fixed-top">
            <div class="container">
                <div class="navbar-header">
                    <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
                        <span class="sr-only">Toggle navigation</span>
                        <span class="icon-bar"></span>
                        <span class="icon-bar"></span>
                        <span class="icon-bar"></span>
                    </button>
                    <a href="#" class="navbar-brand">MSTT Monitor</a>
                </div>
                <div id="navbar" class="collapse navbar-collapse">
                    <ul class="nav navbar-nav">
                    <li class="active"><a href="#">Home</a></li>
                    <li><a href="clients.php">Clients</a></li>
                    <li><a href="#">Settings</a></li>
                    </ul>
                </div>
            </div>
        </nav>

        <div class="container">
            <div class="jumbotron">
                <h1>MSTT Monitor</h1>
                <p>Skriv något nice här, yeeee</p>
            </div>
        </div>
        <script src="/js/jquery-3.1.1.js"></script>
        <script src="/js/bootstrap.min.js"></script>
    </body>
</html>