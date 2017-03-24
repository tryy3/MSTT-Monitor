<?php
	$configString = file_get_contents("config.json");
	$config = json_decode($configString, true);

	function getClients($db) {
		$clients = array();
	
		$stmt = $db->query('SELECT * FROM clients');
		if (!$stmt) return $clients;
		while($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
			$timestampStmt = $db->prepare('SELECT timestamp FROM checks WHERE client_id=? ORDER BY timestamp DESC');
			$timestampStmt->execute(array($row['id']));
			$row['timestamp'] = $timestampStmt->fetch(PDO::FETCH_ASSOC)['timestamp'];
			array_push($clients, $row);
		}
		return $clients;
	}

	function getCheck($db, $client, $check) {
		$stmt = $db->prepare("SELECT error, response FROM clients WHERE client_id=? AND command_id=? ORDER BY timestamp DESC");
		$stmt->execute(array($client, $check));
		if ($stmt->rowCount() <= 0) {
			return array();
		}
		return $stmt->fetchAll(PDO::FETCH_ASSOC);
	}

	try {
		$clients = getClients($monitorDB);
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
				<tr class="clickable-row" data-href="?page=client&id=<?php echo $cl['id']?>">
					<?php foreach($config["ClientList"] as $val) {
						$elem = "td";
						if (isset($val["Bold"]) && $val["Bold"]) {
							$elem = "th";
						}
						if (isset($val["Key"]) && is_string($val["Key"])) {
							echo "<".$elem.">".$cl[$val["Key"]]."</".$elem.">";
							return
						}

						if (isset($val["Function"])) {
							switch ($val["Function"]) {
								case 'warnings':
									if (!isset($val["Check"]) || $val["Check"] < 0) {
										echo "<".$elem.">Invalid check</".$elem.">";
										return;
									}
									$check = $val["Check"];

									$checks = getCheck($monitorDB, $cl["id"], $check);
									if (sizeof($checks) <= 0) {
										echo "<".$elem.">No results</".$elem.">";
									}

									$count = $val["Warnings"][sizeof($val["Warnings"])-1]["Amount"];

									$fails = 0;
									for ($i = 0; $i < $count; $i++) {
										if ($checks[$i]["error"]) {
											$fails++;
											continue;
										}
										$resp = json_decode($checks[$i]["response"], true);
										if ($resp["error"] != "") {
											$fails++;
											continue;
										}
									}

									$prev = array();
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