<?php
    include_once("includes/Graph.php");
	$configString = file_get_contents("config.json");
	$config = json_decode($configString, true);
	use \MSTT_MONITOR\Utils;

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

    function formatResponse($response) {
        $json = json_encode($response, JSON_PRETTY_PRINT);
        $json = preg_replace('/&/', '&amp;', $json);
        $json = preg_replace('/</', '&lt;', $json);
        $json = preg_replace('/>/', '&gt;', $json);
        $json = preg_replace_callback('/("([a-zA-Z0-9]|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/', function($matches) {
            $cls = 'number';
            if (preg_match('/^"/', $matches[0])) {
                if (preg_match('/:$/', $matches[0])) {
                    $cls = 'key';
                } else {
                    $cls = 'string';
                }
            } else if (preg_match('/true|false/', $matches[0])) {
                $cls = 'boolean';
            } else if (preg_match('/null/', $matches[0])) {
                $cls = 'null';
            }
            return '<span class="'.$cls.'">'.$matches[0].'</span>';
        }, $json);
        return $json;
    }

    function getManualColor($time, $manuals) {
        $prev = array();
        foreach ($manuals as $manual) {
            $t = strtotime($manual["Since"]);
            if ($t > $time) {
                $prev = $manual;
            }
        }
        return $prev["Color"];
    }

    function getManualRefresh($id, $target, $commands, $check, $manuals) {
        $color = getManualColor($check, $manuals);
        return "<i class='fa fa-refresh refresh-check pull-right' style='color:".$color.
                    "' data-id='".$id.
                    "' data-target='".$target.
                    "' data-command='".getCommand($target, $commands)->getCommand().
                    "'></i>";
    }

    function fa_icon($bool) {
        return ($bool) ? "fa-check fa-check-green" : "fa-close fa-close-red";
    }

    function getCommand($id, $commands) {
        foreach ($commands as $command) {
            if ($command->getCommandID() == $id) {
                return $command;
            }
        }
    }

    $checks = array();
    try {
        $client = Utils\getClient($monitorDB, $_GET['id']);
        $groups = Utils\getAllGroups($monitorDB);
        $commands = Utils\getAllCommands($monitorDB);

        foreach($config["ClientGraphs"] as $check) {
            $graph = new Graph();
            $graph->Parse($check);
            $graph->FillDataPoints($monitorDB, $client->getID());
            array_push($checks, $graph);
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
        <div class="col-md-4">
            <table class="table table-hover table-bordered">
                <tr>
                    <th colspan="2">
                        User information
                        <button type="button" class="btn btn-danger delete-client inline -right" data-id="<?php echo $client->getID()?>" style="padding: 3px 6px;">Delete client</button>
                    </th>
                </tr>
                <?php foreach ($config["ClientTable"] as $v) { ?>
                    <tr>
                        <?php if (is_string($v["Key"])) {
                            echo "<th>".$v["Namn"]."</th>";
                            if (isset($v["Edit"]) && $v["Edit"]) {
                                echo "<td contenteditable data-previous='".$client->get($v["Key"])."' data-for='client' data-target='".$v["Key"]."' data-id='".$client->getID()."'>".$client->get($v["Key"])."</td>";
                            } else {
                                echo "<td>".$client->get($v["Key"])."</td>";
                            }
                        } else {
                            $check = $client->getChecksByCommandID($v["Key"]);
                            if (!$check || count($check) <= 0) {
                                if (isset($v["Manual"])) {
                                    echo "<th>".$v["Namn"].getManualRefresh($client->getID(), $v["Key"], $commands, -1, $v["Manual"])."</th><td>Invalid command ID.</td>";
                                } else {
                                    "<th>".$v["Namn"]."<th><td>Invalid command ID.</td>";
                                }
                                continue;
                            }
                            $check = $check[0];

                            if (isset($v["Manual"])) {
                                echo "<th>".$v["Namn"].getManualRefresh($client->getID(), $v["Key"], $commands, $check->getTimestamp(), $v["Manual"])."</th>";
                            } else {
                                echo "<th>".$v["Namn"]."</th>";
                            }

                            $value = ParseParams($v["Params"], $check->getResponse());
                            if ($value == "") {
                                echo "<td>Invalid Params</td>";
                                continue;
                            }
                            
                            $command = $client->getCommand($v["Key"]);
                            if ($command && $command->getFormat() == "seconds") {
                                echo "<td>".parse_seconds($value)."</td>";
                            } else {
                                echo "<td>".$value."</td>";
                            }
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
                <?php foreach ($client->getCommands() as $cmd) { ?>
                    <button type="button" class="list-group-item checks" data-target="<?php echo $cmd->getCommandID()?>">
                        <h4 class="list-group-item-heading"><?php echo $cmd->getName() ?></h4>
                        <p class="list-group-item-text">
                            <?php echo $cmd->getDescription() ?>
                        </p>
                    </button>
                <?php } ?>
            </div>
        </div>

        <div class="col-md-5">
            <?php foreach($client->getCommands() as $cmd) { ?>
                <div class="panel panel-primary checks-item" data-check="<?php echo $cmd->getCommandID()?>" style="display:none;">
                    <div class="panel-heading"> 
                        <h3 class="panel-title"><?php echo $cmd->getName() ?></h3>
                    </div>
                    <div class="panel-body">
                        <p>
                            <?php $c = $client->getChecksByCommandID($cmd->getCommandID());
                            if (!$c || count($c) <= 0) {
                                echo "<b>Can't find any check response with this command</b>";
                            } else {
                                $check = $c[0]; ?>

                                <b>Timestamp: </b><?php echo $check->getTimestamp() ?></br>
                                <b>Checked: </b><i class="fa <?php echo fa_icon($check->getChecked()) ?>"></i></br>
                                <b>Finished: </b><i class="fa <?php echo fa_icon($check->getFinished()) ?>"></i></br>
                                <b>Error: </b><i class="fa <?php echo fa_icon($check->getError()) ?>"></i></br>
                                <b>Response: </b>
                                <pre class="manual-output"><?php echo formatResponse($check->getResponse()) ?></pre>
                            <?php } ?>
                        </p>
                    </div>
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
                <?php foreach ($client->getGroups() as $group) { ?>
                    <div class="list-group-item drag group-dd" data-id="<?php echo $client->getID(); ?>" data-type="add" data-target="<?php echo $group->getName()?>">
                        <h4 class="list-group-item-heading"><?php echo $group->getName() ?></h4>
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
                    <div class="list-group-item drag" data-id="<?php echo $client->getID(); ?>" data-target="<?php echo $group->getName() ?>">
                        <h4 class="list-group-item-heading"><?php echo $group->getName() ?></h4>
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
                                    echo "<option data-id='".$cmd->getCommandID()."' data-cmd='".$cmd->getCommand()."'>".$cmd->getName()."</option>";
                                }
                            ?>
                        </select>
                        <span class="input-group-btn">
                            <button class="btn btn-default send-manual-command" data-id="<?php echo $client->getID()?>" type="submit">Send</button>
                        </span>
                    </div>
                </div>
                <div class="panel-body">
                    <pre class="manual-output">

                    </pre>
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