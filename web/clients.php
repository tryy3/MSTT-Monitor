<?php
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
	            <th>#</th>
	            <th class="col-md-3 col-xs-3">Namn</th>
	            <th class="col-md-3 col-xs-3">IP</th>
	            <th class="col-md-3 col-xs-3">Grupper</th>
	            <th class="col-md-3 col-xs-3">Senaste check</th>
	        </tr>
	        <tr class="warning no-result">
	            <!-- colspan måste vara samma antal som kolumer -->
	            <td colspan="5"></i class="fa fa-warning"></i> No result</td>
	        </tr>
	    </thead>
		<tbody>
			<?php foreach ($clients as $cl) { ?>
				<tr class="clickable-row" data-href="?page=client&id=<?php echo $cl["id"]?>">
					<th scope="row"><?php echo $cl["id"]?></th>
					<td data-edit="true"><?php echo $cl["namn"]?></th>
					<td data-edit="true"><?php echo $cl["ip"]?></th>
					<td data-edit="true"><?php echo $cl["group_names"]?></th>
					<td><?php echo $cl["timestamp"]?></th>
				</tr>
			<?php } ?>
		</tbody>
	</table>	    
</div>
<div class="col-md-2"></div>