<?php
    include_once("includes/db_connect.php");
    header('Content-Type: application/json');

    // API version utifall man måste ändra något kritiskt
    // så finns det bakåtkompatibelitet så att klienter
    // som ska uppdatera inte behöver misslyckas.
    $version = "1.0";
    if (isset($_POST["api"])) {
        $version = $_POST["api"];
    }
    if (isset($_GET["api"])) {
        $version = $_GET["api"];
    }

    if (isset($_POST["software"])) {
        $software = $_POST["software"];
    }
    if (isset($_GET["software"])) {
        $software = $_GET["software"];
    }

    switch($version) {
        case "1.0":
            $versions = getVersions($monitorDB);
            if (isset($software)) {
                if (isset($versions[$software])) {
                    echo json_encode($versions[$software], JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES | JSON_NUMERIC_CHECK);
                } else {
                    echo json_encode(array("error"=>"Can't find a software with the name: ".$software));
                }
            } else {
                echo json_encode($versions, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES | JSON_NUMERIC_CHECK);
            }
            break;
        default:
            echo json_encode(array("error"=>"Unsupported api version."));
    }

    function getVersions($db) {
        $out = array();

        $stmt = $db->query('SELECT name, checksum, version, patch, patch_checksum FROM uploads');
        while ($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
            if (!isset($out[$row['name']]))
                $out[$row['name']] = array("versions" => array());

            if (!isset($out[$row['name']]['versions'][$row['version']]))
                $out[$row['name']]['versions'][$row['version']] = array();

            $file = "/uploads/".$row['name'].$row['version'].".exe";
            
            $out[$row['name']]['versions'][$row['version']]['download'] = $file;
            $out[$row['name']]['versions'][$row['version']]['checksum'] = $row['checksum'];
            $out[$row['name']]['versions'][$row['version']]['patch'] = (bool)$row['patch'];
            $out[$row['name']]['versions'][$row['version']]['patch_checksum'] = $row['patch_checksum'];
            if ($row['patch'])
                $out[$row['name']]['versions'][$row['version']]['patch_download'] = $file.".patch";
        }
        return $out;
    }
?>