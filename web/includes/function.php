<?php
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

function ParseParams($params, $value) {
	foreach($params as $p) {
		$value = $value[$p];
	}
	return $value;
}

function isActive($key, $checks, $true = "selected", $false = "") {
    if (is_array($checks)) {
        foreach ($checks as $check) {
            if (is_a($check, 'MSTT_MONITOR\Utils\Command')) {
                if ($key == $check->getCommandID()) {
					if ($true == 'id') {
						return $check->getID();
					}
                    return $true;
                }
            }
        }
    } else if (strpos($checks, ',') !== false) {
        foreach (explode(",", $checks) as $check) {
            if ($key == $check) {
                return $true;
            }
        }
    } else {
        if ($key == $checks) {
            return $true;
        }
    }
    return $false;
}
?>