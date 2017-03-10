<?php
include_once 'includes/db_connect.php';
include_once 'includes/functions.php';

sec_session_start();
 
if (login_check($mysqli) == true) {
    header('Location: protected_page.php');
} else {
    $logged = 'utloggad';
}
?>

<!DOCTYPE html>
<html lang="sv">
    <head>
        <meta charset="utf-8">
        <title>MSTT Monitor logga in</title>
        <link rel="stylesheet" href="css/main.css" />
    </head>
    <body>
    	<?php
        if (isset($_GET['error'])) {
            echo '<p class="error">Error Logging In!</p>';
        }
        ?>
		<div class="loginmodal-container">
			<h1>MSTT Monitor</h1>
			<h1>Logga in</h1><br>
			<form action="includes/process_login.php" method="post" name="login_form">
				<input type="text" id="email" name="email" placeholder="Din@epost.se">
				<input type="password" id="password" name="password" placeholder="Lösenord">
				<input type="submit" class="login loginmodal-submit" value="logga in" onclick="formhash(this.form, this.form.password);">
			</form>
			<div class="login-help"><a href="includes/register.php">Registrera</a> - <a href="#">Glömt Lösenordet</a></div>
		</div>

        <script src="/js/forms.js"></script>
        <script src="/js/sha512.js"></script>
	</body>
</html>