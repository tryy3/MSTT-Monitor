<?php
	$configString = file_get_contents("config.json");
	$config = json_decode($configString, true);

    $a = array("1", "2", "3");
    echo array_rand($a);

    function parse_seconds($seconds)
    {
        $return = '';
        $m = 60;
        $h = $m * 60;
        $d = $h * 24;
        $days = (int)($seconds / $d);
        $seconds = $seconds % $d;
        $hours = (int)($seconds / $h);
        $seconds = $seconds % $h;
        $minutes = (int)($seconds / $m);
        $seconds = $seconds % $m;

        if ($days > 0) {
            $return .= $days.'d';
        }
        if($hours > 0){
            if($days > 0) {
                $return .= ' ';
            }
            $return .= $hours.'h';
        } 
        if($minutes > 0) {
            if($hours > 0) {
                $return .= ' ';
            }
            $return .= $minutes.'m';
        }
        if(empty($return))
            return $seconds.'s';
        else
            return $return;
    }
	
    function getClient($db, $id) {
        $client = array();

        $stmt = $db->prepare('SELECT * FROM clients WHERE id=?');
        $stmt->execute(array($id));
        $client = $stmt->fetch(PDO::FETCH_ASSOC);
        if (empty($client)) {
            return $client;
        }
        
        $timestampStmt = $db->prepare('SELECT timestamp FROM checks WHERE client_id=? ORDER BY timestamp DESC');
        $timestampStmt->execute(array($client['id']));
        $client['latest'] = $timestampStmt->fetch(PDO::FETCH_ASSOC)['timestamp'];
        $client['commands'] = array();
        if ($client['group_names'] != "") {
            $client['groups'] = explode(',', $client['group_names']);
        } else {
            $client['groups'] = array();
        }

        $groupStmt = $db->prepare('SELECT command_id FROM groups WHERE group_name=?');
        $commandStmt = $db->prepare('SELECT namn, description, format FROM commands WHERE id=?');
        foreach ($client['groups'] as $group) {
            $groupStmt->execute(array($group));
            while ($command = $groupStmt->fetch(PDO::FETCH_ASSOC)) {
                $commandStmt->execute(array($command['command_id']));
                
                $cmd = $commandStmt->fetch(PDO::FETCH_ASSOC);
                $client['commands'][$command['command_id']] = array(
                    'command_id'=>$command['command_id'],
                    'namn'=>$cmd['namn'],
                    'description'=>$cmd['description'],
                    'format'=>$cmd['format']
                );
            }
        }
        
        $client['checks'] = array();
        $allChecksStmt = $db->prepare('SELECT timestamp, response, checked, error, finished, command_id FROM checks WHERE client_id=? ORDER BY timestamp DESC');
        $allChecksStmt->execute(array($client['id']));
        while($row = $allChecksStmt->fetch(PDO::FETCH_ASSOC)) {
            if (!isset($client['checks'][$row["command_id"]])) {
                $client['checks'][$row["command_id"]] = array(
                    "timestamp" => $row["timestamp"],
                    "checked" => $row["checked"],
                    "error" => $row["error"],
                    "finished" => $row["finished"],
                    "response" => json_decode($row['response'], true),
                    "command_id" => $row["command_id"]
                );
            }
        }

        return $client;
    }

    function navigateResp($array, $params) {
        foreach($params as $param) {
            if (!isset($array[$param])) {
                $array = "";
                break;
            }
            $array = $array[$param];
        }
        return $array;
    }

    function getCheck($stmt, $id, $options) {
        $from = strtotime($options["from"]);
        $to = time();
        if (isset($options["to"]) && $options["to"] != "now") {
            $options["to"] = strtotime($to);
        }

        $stmt->execute(array($id, $options["check"], $from, $to));
        $check = array();
        while($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
            $x = strtotime($row["timestamp"]);
            $resp = json_decode($row["response"], true);
            foreach($options["DataPointsOptions"] as $k => $opt) {
                if (!isset($check[$k])) {
                    $check[$k] = array_merge($options["DataOptions"], array("dataPoints" => array()));
                }
                $data = array();
                $v = navigateResp($resp, $opt["Params"]);
                if ($v == "") {
                    continue;
                }
                $data["x"] = $x;
                $data["y"] = $v;
                $data["label"] = date("Y/m/d H:i:s", $x);
                if (isset($opt["YFormat"])) {
                    switch($opt["YFormat"]) {
                        case "procent":
                            $data["toolTipContent"] = "<span style='\"'color: {color};'\"'>{label}:</span> ".number_format($v,2)."%";
                            break;
                        case "GB":
                            $data["toolTipContent"] = "<span style='\"'color: {color};'\"'>{label}:</span> ".number_format($v/1000/1000/1000,2)." GB";
                            break;
                        case "MB":
                            $data["toolTipContent"] = "<span style='\"'color: {color};'\"'>{label}:</span> ".number_format($v/1000/1000,2)." MB";
                            break;
                        case "KB":
                            $data["toolTipContent"] = "<span style='\"'color: {color};'\"'>{label}:</span> ".number_format($v/1000,2)." KB";
                            break;
                        case "B":
                            $data["toolTipContent"] = "<span style='\"'color: {color};'\"'>{label}:</span> ".number_format($v,2)." B";
                            break;
                    }
                    unset($opt["YFormat"]);
                }
                unset($opt["Params"]);
                array_push($check[$k]["dataPoints"], array_merge($data, $opt));
            }
        }
        return $check;
    }

    function getGroups($db) {
        $groups = array(); 

        $stmt = $db->query("SELECT group_name FROM groups");
        if (!$stmt) return $groups;
		while($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
            array_push($groups, $row["group_name"]);
        }
        return array_unique($groups);
    }

    function getCommands($db) {
        $cmds = array();
        $stmt = $db->query("SELECT command, namn, id FROM commands");
        if (!$stmt) return $cmds;
        return $stmt->fetchAll(PDO::FETCH_ASSOC);
    }

    function fa_icon($bool) {
        return ($bool) ? "fa-check fa-check-green" : "fa-close fa-close-red";
    }

    function convertMemory($type) {
        return ($type) ? "Swap" : "RAM";
    }

    function generate_size_buttons_array($id, $array) {
        echo '<span class="btn-group btn-group-xs" role="group" style="display:inline-block;">';
            foreach($array as $format) {
                echo '<button type="button" class="button-convert-size btn btn-default" data-target="'.$id.'" role="group">'.$format.'</button>';
            }
        echo "</span>";
    }

    function generate_size_buttons($id, $format) {
        $bytes = array("B", "KiB", "MiB", "GiB", "TiB");
        $bits = array("b", "Kib", "Mib", "Gib", "Tib");

        echo '<span class="btn-toolbar" role="toolbar" style="display:inline-block;">';
            if ($format == "memory" || $format == "disc") {
                generate_size_buttons_array($id,$bytes);
            } else {
                generate_size_buttons_array($id,$bytes);
                generate_size_buttons_array($id,$bits);
            }
        echo "</span>";
    }

    $checks = array();
    try {
        $client = getClient($monitorDB, $_GET['id']);
        $groups = getGroups($monitorDB);
        $commands = getCommands($monitorDB);

        $stmt = $monitorDB->prepare("SELECT timestamp, response, checked, error, finished FROM checks WHERE client_id=? AND command_id=? AND timestamp >= FROM_UNIXTIME(?) AND timestamp <= FROM_UNIXTIME(?)");
        foreach($config["ClientGraphs"] as $check) {
            $arr = array();
            if (!isset($check["Check"]) || !isset($check["From"])) {
                break;
            }
            if(isset($check["To"])) {
                $arr["to"] = $check["To"];
            }
            $arr["check"] = $check["Check"];
            $arr["from"] = $check["From"];
            $arr["DataPointsOptions"] = $check["DataPointsOptions"];
            $arr["DataOptions"] = $check["DataOptions"];
            $ch = getCheck($stmt, $client["id"], $arr);

            $a = array();
            if (isset($check["ChartOptions"])) {
                $a["ChartOptions"] = $check["ChartOptions"];
            } else {
                $a["ChartOptions"] = array();
            }
            $a["dataPoints"] = $ch;
            array_push($checks, $a);
        }
    } catch(PDOException $ex) {
        $client = array();
        $groups = array();
        $commands = array();
    }
?>

<div class="col-md-2"></div>

<div class="col-md-8">
    
	<h1 class="page-header">MSTT Monitor - Client page</h1>
	<div class="alert" style="display:none;">
        <b class="alert-header"></b><span class="alert-text"></span>
	</div>

    <?php if (empty($client)) : ?>

    <h3> Can't find this client</h3>

    <?php else: ?>

    <div class="row">
        <div class="col-md-3">
            <table class="table table-hover table-bordered">
                <tr>
                    <th colspan="2">
                        User information
                        <button type="button" class="btn btn-danger delete-client inline pull-right" data-id="<?php echo $client['id']?>" style="padding: 3px 6px;">Delete client</button>
                    </th>
                </tr>
                <?php foreach ($config["ClientTable"] as $v) { ?>
                    <tr>
                        <th><?php echo $v["Namn"] ?></th>
                        <?php if (is_string($v["Key"])) {
                            if (isset($v["Edit"]) && $v["Edit"]) {
                                echo "<td contenteditable data-previous='".$client[$v["Key"]]."' data-for='client' data-target='".$v["Key"]."' data-id='".$client['id']."'>".$client[$v["Key"]]."</td>";
                            } else {
                                echo "<td>";
                                echo $client[$v["Key"]];
                                echo "</td>";
                            }
                        } else {
                            echo "<td>";
                            if (isset($client["checks"][$v["Key"]]["response"])) {
                                $value = $client["checks"][$v["Key"]]["response"];
                                foreach($v["Params"] as $param) {
                                    if (!isset($value[$param])) {
                                        $value = "";
                                        break;
                                    }
                                    $value = $value[$param];
                                }
                                if ($value != "") {
                                    if (isset($client["commands"][$v["Key"]]) && $client["commands"][$v["Key"]]["format"] == "seconds") {
                                        echo parse_seconds($value);
                                    } else {
                                        echo $value;
                                    }
                                }
                            }
                            echo "</td>";
                        }?>
                    </tr>
                <?php } ?>
            </table>
        </div>

        <div class="col-md-3">
            <div class="list-group">
                <a href="#" class="list-group-item active">
                    <h4 class="list-group-item-heading">Checks</h4>
                </a>
                <?php foreach ($client['commands'] as $cmd) { ?>
                    <button type="button" class="list-group-item checks" data-target="<?php echo $cmd['command_id']?>">
                        <h4 class="list-group-item-heading"><?php echo $cmd['namn'] ?></h4>
                        <p class="list-group-item-text">
                            <?php echo $cmd['description'] ?>
                        </p>
                    </button>
                <?php } ?>
            </div>
        </div>

        <div class="col-md-6">
            <?php foreach($client['commands'] as $cmd) { ?>
                <div class="panel panel-primary checks-item" data-check="<?php echo $cmd['command_id']?>" style="display:none;">
                    <div class="panel-heading"> 
                        <h3 class="panel-title"><?php echo $cmd['namn'] ?></h3>
                    </div>
                    <div class="panel-body">
                        <p>
                            <b>Timestamp: </b><?php echo $client['checks'][$cmd['command_id']]['timestamp']?><br>
                            <b>Checked: </b><i class="fa <?php echo fa_icon($client['checks'][$cmd['command_id']]['checked']) ?>"></i><br>
                            <b>Finished: </b><i class="fa <?php echo fa_icon($client['checks'][$cmd['command_id']]['finished']) ?>"></i><br>
                            <b>Error: </b><i class="fa <?php echo fa_icon($client['checks'][$cmd['command_id']]['error']) ?>"></i>
                            <?php foreach($client['checks'][$cmd['command_id']]['response'] as $key => $value) {
                                if ($key == "error") continue; ?>
                                <br>
                                <b><?php echo ucwords($key); ?>: </b>
                                <?php if ($cmd['format'] == 'disc') {
                                    generate_size_buttons($cmd['command_id'], $cmd['format']);
                                }?>
                                <?php if (is_array($value)) { ?>
                                    <table class="table table-hover table-bordered table-condensed">
                                        <tr>
                                            <?php foreach($value[0] as $k => $vr) { ?>
                                                <th><?php echo ucwords($k) ?></th>
                                            <?php } ?>
                                        </tr>
                                        <?php foreach($value as $val) { ?>
                                            <tr>
                                                <?php foreach($val as $ky => $v) { ?>
                                                    <?php if ($ky == "size") { ?>
                                                        <td>
                                                            <span class="convert-size" <?php echo 'data-identifier="'.$cmd['command_id'].'" data-value="'.$v.'"data-format="'.$cmd['format'].'"'?>>
                                                                <?php echo $v; ?>
                                                            </span>
                                                            <span class="convert-size-format" data-identifier="<?php echo $cmd['command_id']?>">B</span>
                                                        </td>
                                                    <?php } else { ?>
                                                        <td><?php echo $v ?></td>
                                                    <?php } ?>
                                                <?php } ?>
                                            </tr>
                                        <?php } ?>
                                    </table>
                                <?php } else { ?>
                                    <?php if ($cmd['format'] == 'memory') {?>
                                        <?php if ($key == "type") {
                                            echo convertMemory($value);
                                        } elseif ($key == "size") {?>
                                            <span class="convert-size" <?php echo 'data-identifier="'.$cmd['command_id'].'" data-value="'.$value.'"data-format="'.$cmd['format'].'"'?>>
                                                <?php echo $value; ?>
                                            </span>
                                            <span class="convert-size-format" data-identifier="<?php echo $cmd['command_id']?>">B</span>
                                            <?php generate_size_buttons($cmd['command_id'], $cmd['format']);?>
                                        <?php }
                                    } else {
                                        echo $value;
                                    }?>
                                <?php } ?>
                            <?php } ?>
                        </p>
                    </div>

                    <!--<?php //if (is_array($client['checks'][$cmd['command_id']]['response'])) { ?>
                        <table class="table table-hover table-bordered table-condensed">
                            <?php //foreach ($client['checks'][$cmd['command_id']]['response'] as $value) { ?>
                                <pre>
                                <?php //echo $value ?>
                                </pre>
                                <tr>
                                    <th width="15%"><?php //echo $value[0]; ?>
                                    <td width="85%"><?php //echo $value[1]; ?>
                                </tr>
                            <?php //} ?>
                        </table>
                    <?php //} ?>-->
                </div>
            <?php } ?>
        </div>
    </div>
    
    <div class="row">
        <div class="col-md-3">
            <div class="list-group group-dd" data-type="add" >
                <div href="#" class="list-group-item active">
                    <h4 class="list-group-item-heading">Client Groups</h4>
                </div>
                <?php foreach ($client['groups'] as $group) { ?>
                    <div class="list-group-item drag group-dd" data-id="<?php echo $client['id']; ?>" data-type="add" data-target="<?php echo $group?>">
                        <h4 class="list-group-item-heading"><?php echo $group ?></h4>
                    </div>
                <?php } ?>
            </div>
        </div>
        <div class="col-md-3">
            <div class="list-group group-dd" data-type="remove">
                <div href="#" class="list-group-item active">
                    <h4 class="list-group-item-heading">Groups</h4>
                </div>
                <?php foreach($groups as $group) { ?>
                    <div class="list-group-item drag" data-id="<?php echo $client['id']; ?>" data-target="<?php echo $group ?>">
                        <h4 class="list-group-item-heading"><?php echo $group ?></h4>
                    </div>
                <?php }?>
            </div>
        </div>
        <div class="col-md-6">
            <div class="panel panel-primary">
                <div class="panel-heading">
                    <h3 class="panel-title">Manual check</h3>
                    
                    <span style="padding-right: 20px">
                        <input type="checkbox" name="Manual-Switch" checked>
                    </span>
                    
                    <span>
                        <input type="checkbox" name="Save-Mysql" checked>
                    </span>
                    
                    <div class="input-group" style="padding: 10px 0px">
                        <input type="text" class="check-command form-control manual-command" style="width:70%" placeholder="Command">
                        <select class="form-control manual-dropdown" style="width:30%" disabled>
                            <option selected disabled>Command List</option>
                            <?php
                                foreach ($commands as $cmd) {
                                    echo "<option data-id='".$cmd["id"]."' data-cmd='".$cmd["command"]."'>".$cmd["namn"]."</option>";
                                }
                            ?>
                        </select>
                        <span class="input-group-btn">
                            <button class="btn btn-default send-manual-command" type="submit">Send</button>
                        </span>
                    </div>
                </div>
                <div class="panel-body">
                    <div class="manual-output">

                    </div>
                </div>
            </div>
        </div>
    </div>
    <div class="row">
    <?php
        foreach($checks as $k => $c) {
            echo "<div id='graphCheck' data-check='".$k."' style=\"width: 100%; height: 400px; display: inline-block;\"></div>";
        }
    ?>
    </div>
    <?php endif; ?>
</div>

<div class="col-md-2"></div>