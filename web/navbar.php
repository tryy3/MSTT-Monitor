<nav class="navbar navbar-inverse navbar-default" style="border-radius:0px;">
  <div class="container-fluid">
    <!-- Brand and toggle get grouped for better mobile display -->
    <div class="navbar-header">
      <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1" aria-expanded="false">
        <span class="sr-only">Toggle navigation</span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
      </button>
      <a class="navbar-brand" href="index.php">MSTT Monitor</a>
    </div>

    <!-- Collect the nav links, forms, and other content for toggling -->
    <div class="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
      <ul class="nav navbar-nav">
      	<li <?php echo $active = active($page, 'start'); ?>><a href="?page=start">Start<span class="sr-only"></span></a></li>
        <li <?php echo $active = active($page, 'clients');?>><a href="?page=clients">Clients<span class="sr-only"></span></a></li>
        <li <?php echo $active = active($page, 'upload');?>><a href="?page=upload">Upload<span class="sr-only"></span></a></li>
        <li <?php echo $active = active($page, 'settings');?>><a href="?page=settings">Settings<span class="sr-only"></span></a></li>
        <li <?php echo $active = active($page, 'documentation');?>><a href="?page=documentation">Documentation<span class="sr-only"></span></a></li>
        <li <?php echo $active = active($page, 'success'); ?>><a href="?page=success">Success<span class="sr-only"></span></a></li>
        <li <?php echo $active = active($page, 'test'); ?>><a href="?page=test">Test<span class="sr-only"></span></a></li>
      </ul>
     
      <ul class="nav navbar-nav navbar-right">
        <li <?php echo $active = active($page, 'adminindex'); ?>><a href="?page=adminindex"><?php echo htmlentities($_SESSION['username']); ?></a></li>
        <li class="dropdown">
          <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">Konto<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><a href="../../">Intranätet</a></li>
            <li><a href="#">Länk 2</a></li>
            <li><a href="#">Länk 3</a></li>
            <li role="separator" class="divider"></li>
            <li><a href="includes/logout.php">Logga ut</a></li>
          </ul>
        </li>
      </ul>
    </div><!-- /.navbar-collapse -->
  </div><!-- /.container-fluid -->
</nav>