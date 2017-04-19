<?php

include_once '../includes/db_connect.php';
include_once '../includes/function.php';
include_once 'includes/db_connect.php';
include_once 'includes/functions.php';
include_once 'includes/function.php';
include_once 'includes/MSTT_Monitor/utils/utils.php';
 
sec_session_start();

if(isset($_GET['page'])){
	
	$page = $_GET['page']; 
}
else {
	
	$page = "start";
}
?>

<?php if (login_check($mysqli) == true) : ?>
	
<!DOCTYPE html>
<html lang="en">
    <head>
       <?php include 'head.php'; ?>
    </head>
    <body>
		<?php include_once ('navbar.php'); ?>
		<div><!-- Start Allt-->	
		
		<div class="row">	
			<?php include ($page . '.php'); ?>
		</div>
		
		</div><!-- Stop Allt-->
		<?php include 'js/js.php';?>
	</body>	
</html>    
<?php else : header("location: ../index.php"); ?>   
	
<?php endif; ?> 

