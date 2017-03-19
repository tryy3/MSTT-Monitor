<?php
function active($page, $pagename){
	if($page == $pagename){
		$result = 'class="active"';
	}
	else {
		$result = "";
	}
	return $result;
}

function aktiv($page, $pagename){
	
	if($page == $pagename){
		
		$result = 'active';
	}
	else {
		$result = "";
	}
	return $result;
}

function getLatestVersion($db, $name) {
	$stmt = $db->prepare('SELECT version FROM uploads WHERE name=? ORDER BY version DESC');
	$stmt->execute(array($name));
	return $stmt->fetch(PDO::FETCH_ASSOC)['version'];
}

function checksum($name) {
	$md5_checksum = hash_file('sha256', $name);
	if (!$md5_checksum) return false;
	return $md5_checksum;
}

function toBool($var) {
	if (!is_string($var)) return $var;
	switch (strtolower($var)) {
		case '1':
		case 'true':
		case 'on':
		case 'yes':
		case 'y':
			return true;
		case '0':
		case 'false':
		case 'off':
		case 'no':
		case 'n':
			return false;
		default:
			return $var;
	}
}

function getServers($db) {
	$out = array();
	$stmt = $db->query("SELECT id, ip, namn FROM servers");
	while($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
		array_push($out, $row);
	}
	return $out;
}
?>