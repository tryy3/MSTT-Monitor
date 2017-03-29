<?php
	$configString = file_get_contents("config.json");
	$config = json_decode($configString, true);
	use \MSTT_MONITOR\Utils;

	try {
		$limit = -1;
		foreach ($config["ClientList"] as $val) {
			if (!isset($val["Warnings"])) {
				continue;
			}
			$v = $val["Warnings"][count($val["Warnings"])-1]["Amount"];
			if ($v > $limit) {
				$limit = $v;
			}
		}
		$clients = Utils\getAllClients($monitorDB, $limit);
		echo count($clients);
	} catch (PDOException $ex) {
		$clients = array();
	}
?>

<div class="col-md-2"></div>

<div class="col-md-8">
	<h1 class="page-header">MSTT Monitor</h1>
	<div class="alert" style="display:none;">
        <b class="alert-header"></b><span class="alert-text"></span>
	</div>
	<form class="form-inline">
		<div class="form-group pull-left">
			<input type="text" class="create-client form-control" placeholder="Klient IP">
			<button type="button" class="btn btn-create-client btn-default" data-toggle="false">Skapa ny klient</button>
		</div>
	</form>
	<div class="form-inline">
		<div class="form-group pull-right">
			<input type="text" class="search form-control" placeholder="Sök efter en klient">
		</div>
	</div>
	<span class="counter pull-right"></span>
	<table class="table table-hover table-bordered results">
	    <thead>
			<tr>
				<?php foreach ($config["ClientList"] as $val) { ?>
					<th><?php echo $val["Namn"] ?></th>
				<?php } ?>
			</tr>
	        <tr class="warning no-result">
	            <!-- colspan måste vara samma antal som kolumer -->
	            <td colspan="5"></i class="fa fa-warning"></i> No result</td>
	        </tr>
	    </thead>
		<tbody>
			<?php foreach ($clients as $cl) { ?>
				<tr class="clickable-row" data-href="?page=client&id=<?php echo $cl->getID()?>">
					<?php foreach($config["ClientList"] as $val) {
						$elem = "td";
						if (isset($val["Bold"]) && $val["Bold"]) {
							$elem = "th";
						}
						if (isset($val["Key"]) && is_string($val["Key"])) {
							echo "<".$elem.">".$cl->get($val["Key"])."</".$elem.">";
							continue;
						}

						if (isset($val["Function"])) {
							switch ($val["Function"]) {
								case 'warnings':
									if (!isset($val["Check"]) || $val["Check"] < 0) {
										echo "<".$elem.">Invalid check</".$elem.">";
										continue;
									}
									$checks = $cl->getChecksByCommandID($val["Check"]);
									if (count($checks) <= 0) {
										echo "<".$elem.">No results</".$elem.">";
										continue;
									}

									$count = $val["Warnings"][count($val["Warnings"])-1]["Amount"];
									$count = ($count > count($checks)) ? count($checks) : $count;
									$fails = 0;
									for ($i = 0; $i < $count; $i++) {
										if ($checks[$i]->getError()) {
											$fails++;
											continue;
										}
										if ($checks[$i]->getResponse()["error"] != "") {
											$fails++;
											continue;
										}
									}
									$prev = array();
									foreach ($val["Warnings"] as $warning) {
										if ($warning["Amount"] > $fails) {
											continue;
										}
										$prev = $warning;
									}
									echo "<".$elem."><i class='fa fa-".$prev["Symbol"]."' style='color:".$prev["Color"]."'></i></".$elem.">";
									break;
								default:
									echo "<".$elem.">Invalid Function</".$elem.">";
									break;
							}
						}
					} ?>
				</tr>
			<?php } ?>
		</tbody>
	</table>	    
</div>
<div class="col-md-2"></div>