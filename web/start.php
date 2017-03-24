<?php
    include_once("includes/Graph.php");
	$configString = file_get_contents("config.json");
	$config = json_decode($configString, true);

    $graphs = array();

    foreach ($config["FrontPage"] as $opt) {
        $graph = new Graph();
        $graph->Parse($opt);
        $graph->FillDataPoints($monitorDB);
        array_push($graphs, $graph);
    }
?>

<div class="col-md-2"></div>

<div class="col-md-8">
<h1>Start</h1>
    <div class="row">
        <?php
            foreach($graphs as $k => $g) {
                echo "<div id='graphCheck' data-check='".$k."' style=\"width: 100%; height: 400px; display: inline-block;\"></div>";
            }
        ?>
    </div>
</div>
<div class="col-md-2"></div>