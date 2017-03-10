<?php
include_once 'register.inc.php';    
include_once 'db_connect.php';
include_once 'functions.php';
 
sec_session_start();
?>
<!DOCTYPE html>
<html lang="sv">
    <head>
        <meta charset="UTF-8">
        <title>Registrera</title>
        <script type="text/JavaScript" src="../js/sha512.js"></script> 
        <script type="text/JavaScript" src="../js/forms.js"></script>
    </head>
    <body>
       
        <div id="main">
        <div id="kunddiv">
                <h1>Inställningar</h1>
        </div>
          <!-- Registration form to be output if the POST variables are not
        set or if the registration script caused an error. -->
        <h1 style="text-align: center;">Lägg till en ny användare.</h1>
        <?php
        if (!empty($error_msg)) {
            echo $error_msg;
        }
        ?>
        <form action="<?php echo esc_url($_SERVER['PHP_SELF']); ?>"method="post" name="registration_form" class="basic-grey">
            <label>
            <span>Användarnamn :</span><input type='text'name='username' id='username' />
            </label>
            <label>
            <span>E-post adress :</span><input type="text" name="email" id="email" /><br>
            </label>
            <label>
            <span>Lösenord :</span><input type="password"name="password" id="password"/><br>
            </label>
            <label>
            <span>Lösenord igen :</span><input type="password" name="confirmpwd" id="confirmpwd" /><br>
            </label>
            <label>
            <span>&nbsp;</span>
            <input type="button" value="Registrera" class="button" onclick="return regformhash(this.form, this.form.username, this.form.email, this.form.password, this.form.confirmpwd);" /> 
            </lable>
        </form>
        </div>
</body>
</html>
        