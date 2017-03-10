<?php
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

        $groupStmt = $db->prepare('SELECT command_id FROM groups WHERE group_name=?');
        $commandStmt = $db->prepare('SELECT namn, description, format FROM commands WHERE id=?');
        foreach (explode(',', $client['group_names']) as $group) {
            $groupStmt->execute(array($group));
            while ($command = $groupStmt->fetch(PDO::FETCH_ASSOC)) {
                $commandStmt->execute(array($command['command_id']));
                
                $cmd = $commandStmt->fetch(PDO::FETCH_ASSOC);
                array_push($client['commands'], array(
                    'command_id'=>$command['command_id'],
                    'namn'=>$cmd['namn'],
                    'description'=>$cmd['description'],
                    'format'=>$cmd['format']
                ));
            }
        }

        $checkStmt = $db->prepare('SELECT timestamp, response, checked, error, finished FROM checks WHERE client_id=? AND command_id=? ORDER BY timestamp DESC');
        foreach($client['commands'] as $key => $cmd) {
            $checkStmt->execute(array($client['id'], $cmd['command_id']));
            $check = $checkStmt->fetch(PDO::FETCH_ASSOC);

            if (!isset($client['commands'][$key]['latest_check'])) {
                $c = array(
                    "timestamp" => $check['timestamp'],
                    "checked" => $check['checked'],
                    "error" => $check['error'],
                    "finished" => $check['finished']
                );

                if (strpos($check['response'], ",") !== false) {
                    $resp = explode(",", $check['response']);
                    
                    foreach($resp as $key => $value) {
                        if (strpos($value, "=")) {
                            $resp[$key] = explode("=", $value);
                        }
                    }
                    $c['response'] = $resp;
                } else {
                    $c['response'] = $check['response'];
                }
                $client['commands'][$key]['latest_check'] = $c;
            }
        }

        return $client;
    }

    function fa_icon($bool) {
        return ($bool) ? "fa-check fa-check-green" : "fa-close fa-close-red";
    }

    function generate_size_buttons($id) {
        $bytes = array("B", "KiB", "MiB", "GiB", "TiB");
        $bits = array("b", "Kib", "Mib", "Gib", "Tib");

        echo '<div class="btn-toolbar" role="toolbar">';
            echo '<div class="btn-group role="group">';
                foreach ($bytes as $format) {
                    echo '<button type="button" class="button-convert-size btn btn-default" data-target="'.$id.'" role="group">'.$format."</button>";
                }
            echo '</div><div class="btn-group" role="group">';
                foreach ($bits as $format) {
                    echo '<button type="button" class="button-convert-size btn btn-default" data-target="'.$id.'" role="group">'.$format."</button>";
                }
            echo "</div>";
        echo "</div>";
    }

    try {
        $client = getClient($monitorDB, $_GET['id']);
    } catch(PDOException $ex) {
        $client = array();
    }
?>

<div class="col-md-2"></div>

<div class="col-md-8">
    
	<h1 class="page-header">MSTT Monitor - Client page</h1>

    <?php if (empty($client)) : ?>

    <h3> Can't find this client</h3>

    <?php else: ?>

    <div class="row">
        <div class="col-md-3">
            <table class="table table-hover table-bordered">
                <tr>
                    <th colspan="2">User information</th>
                </tr>
                <tr>
                    <th>ID</th>
                    <td><?php echo $client['id'] ?></td>
                </tr>
                <tr>
                    <th>Namn</th>
                    <td><?php echo $client['namn'] ?></td>
                </tr>
                <tr>
                    <th>IP</th>
                    <td><?php echo $client['ip'] ?></td>
                </tr>
                <tr>
                    <th>Groups</th>
                    <td><?php echo $client['group_names'] ?></td>
                </tr>
                <tr>
                    <th>latest</th>
                    <td><?php echo $client['latest'] ?></td>
                </tr>
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
                            <?php echo $cmd['desc'] ?>
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
                            <b>Timestamp: </b><?php echo $cmd['latest_check']['timestamp']?><br>
                            <b>Checked: </b><i class="fa <?php echo fa_icon($cmd['latest_check']['checked']) ?>"></i><br>
                            <b>Finished: </b><i class="fa <?php echo fa_icon($cmd['latest_check']['finished']) ?>"></i><br>
                            <b>Error: </b><i class="fa <?php echo fa_icon($cmd['latest_check']['error']) ?>"></i>
                            <?php if (is_string($cmd['latest_check']['response'])) { ?>
                                <br>
                                <b>Response:</b> <span class="convert-size" <?php echo 'data-identifier="'.$cmd['command_id'].'" data-value="'.$cmd['latest_check']['response'].'"data-format="'.$cmd['format'].'"'?>>
                                    <?php
                                        echo $cmd['latest_check']['response'];
                                        generate_size_buttons($cmd['command_id']);
                                    ?>
                                </span>
                            <?php } ?>
                        </p>
                    </div>

                    <?php if (is_array($cmd['latest_check']['response'])) { ?>
                        <table class="table table-hover table-bordered table-condensed">
                            <?php foreach ($cmd['latest_check']['response'] as $value) { ?>
                                <tr>
                                    <th width="15%"><?php echo $value[0]; ?>
                                    <td width="85%"><?php echo $value[1]; ?>
                                </tr>
                            <?php } ?>
                        </table>
                    <?php } ?>
                </div>
            <?php } ?>
        </div>
    </div>

    <?php endif; ?>
</div>

<div class="col-md-2"></div>