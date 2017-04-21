<?php
	use \MSTT_MONITOR\Utils;

    function getGroups($db) {
        $out = array();
        $stmt = $db->query("SELECT id, command_id, group_name, next_check, stop_error FROM groups");
        while($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
            if (!isset($out[$row['group_name']]))
                $out[$row['group_name']] = array('id'=>$row['id'],'group_name'=>$row['group_name']);
            
            if (!isset($out[$row['group_name']]['commands']))
                $out[$row['group_name']]['commands'] = array();
            
            $cmd = array(
                'id' => $row['id'],
                'command_id' => $row['command_id'],
                'next_check' => $row['next_check'],
                'stop_error' => $row['stop_error']
            );

            array_push($out[$row['group_name']]['commands'], $cmd);
        }
        return $out;
    }

    function getCommands($db) {
        $out = array();
        $stmt = $db->query("SELECT id, command, namn, description, format FROM commands");
        while($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
            $out[$row['id']] = $row;
        }
        return $out;
    }

    function isSelected($input, $val) {
        if ($input == $val) return 'selected="selected"';
    }

    function findCommand($id, $commands) {
        foreach ($commands as $cmd) {
            if ($cmd->getCommandID() == $id) {
                return $cmd->getName();
            }
        }
    }

    try {
        $groups = Utils\getAllGroups($monitorDB);
        $commands = Utils\getAllCommands($monitorDB);
        //$groups = getGroups($monitorDB);
        //$commands = getCommands($monitorDB);
        $servers = getServers($monitorDB);
    } catch(PDOException $ex) {
        echo $ex;
    }
?>
<div class="col-md-1"></div>

<div class="col-md-10">
    <div class="alert" style="display:none;">
        <b class="alert-header"></b><span class="alert-text"></span>
    </div>
    <div class="row">
        <div class="col-md-2">
            <div class="list-group">
                <a href="#" class="list-group-item active">
                    <h4 class="list-group-item-heading">Groups</h4>
                    <div class="list-group-item-text">
                        <div class="input-group">
                            <input type="text" class="group-name form-control" placeholder="Group name">
                            <span class="input-group-btn">
                                <button class="btn btn-default add-group" type="button">Add group</button>
                            </span>
                        </div>
                    </div>
                </a>
                <?php foreach ($groups as $group) { ?>
                    <button type="button" class="list-group-item checks" data-target="<?php echo $group->getName()?>">
                        <h4 class="list-group-item-heading">
                            <?php echo $group->getName() ?>
                            <i class="delete-group fa fa-close fa-close-red fa-lg pull-right"></i>
                        </h4>
                    </button>
                <?php } ?>
            </div>
        </div>
        <div class="col-md-4 group-list">
            <?php foreach($groups as $group) { ?>
                <div class="panel panel-primary checks-item drop-group" data-check="<?php echo $group->getName()?>" style="display:none;">
                    <div class="panel-heading"> 
                        <h3 class="panel-title"><?php echo $group->getName() ?></h3>
                        <select
                            class="selectpicker"
                            data-for="toggle_group"
                            data-group="<?php echo $group->getName() ?>"
                            multiple>
                            <?php foreach ($commands as $cmd) { ?>
                                <option 
                                    <?php echo isActive($cmd->getCommandID(), $group->getCommands()) ?>
                                    data-id="<?php echo $cmd->getCommandID()?>"
                                    data-cmd="<?php echo isActive($cmd->getCommandID(), $group->getCommands(), 'id')?>">
                                    <?php echo $cmd->getName() ?>
                                </option>
                            <?php } ?>
                        </select>
                    </div>
                    
                    <table class="table table-groups table-hover table-bordered table-condensed">
                        <tr>
                            <th width=5%>ID</th>
                            <th width=45%>Namn</th>
                            <th width=25%>Nästa Check</th>
                            <th width=25%>Stop Error</th>
                        </tr>
                        <?php foreach($group->getCommands() as $cmd) { ?>
                            <tr>
                                <td><?php echo $cmd->getID() ?></td>
                                <td><?php echo findCommand($cmd->getCommandID(), $commands) ?></td>
                                <td 
                                    contenteditable
                                    data-previous="<?php echo $cmd->getNextCheck()?>"
                                    data-for="group"
                                    data-target="next_check"
                                    data-id="<?php echo $cmd->getID()?>">
                                    <?php echo $cmd->getNextCheck()?>
                                </td>
                                <td>
                                    <select
                                        data-for="group"
                                        data-id="<?php echo $cmd->getID()?>"
                                        data-target="stop_error"
                                        class="form-control"
                                        style="width:80%; display:inline">
                                        <option <?php echo isSelected(toBool($cmd->getStopError()), true) ?>>True</option>
                                        <option <?php echo isSelected(toBool($cmd->getStopError()), false) ?>>False</option>
                                    </select>
                                    <i data-id="<?php echo $cmd->getID()?>" class="remove-command-group fa fa-close fa-close-red fa-lg pull-right"></i>
                                </td>
                            </tr>
                        <?php } ?>
                    </table>
                </div>
            <?php } ?>
        </div>
        <div class="col-md-6">
            <div class="panel panel-primary drop-command">
                <div class="panel-heading">
                    <h3 class="panel-title">Commands</h3>
                    <div class="panel-body">
                        <div class="input-group">
                            <button class="btn btn-default add-command" type="button">Add Command</button>
                        </div>
                    </div>
                </div>
                <table class="table table-commands table-hover table-bordered table-condensed">
                    <tr>
                        <th width=5%>ID</th>
                        <th width=15%>Namn</th>
                        <th width=25%>Command</th>
                        <th width=35%>Description</th>
                        <th width=20%>Format</th>
                    </tr>
                    <?php foreach($commands as $cmd) { ?>
                        <tr>
                            <td><?php echo $cmd->getCommandID()?></td>
                            <td
                                contenteditable
                                data-for="cmd"
                                data-id="<?php echo $cmd->getCommandID()?>"
                                data-target="namn"
                                data-previous="<?php echo $cmd->getName()?>">
                                <?php echo $cmd->getName()?>
                            </td>
                            <td
                                contenteditable
                                data-for="cmd"
                                data-id="<?php echo $cmd->getCommandID()?>"
                                data-target="command"
                                data-previous="<?php echo $cmd->getCommand()?>">
                                <?php echo $cmd->getCommand()?>
                            </td>
                            <td
                                contenteditable
                                data-for="cmd"
                                data-id="<?php echo $cmd->getCommandID()?>"
                                data-target="description"
                                data-previous="<?php echo $cmd->getDescription()?>">
                                <?php echo $cmd->getDescription()?>
                            </td>
                            <td>
                                <select
                                    data-for="cmd"
                                    data-id="<?php echo $cmd->getCommandID()?>"
                                    data-target="format" class="form-control"
                                    style="width:80%; display:inline">
                                    <option <?php echo isSelected($cmd->getFormat(), '') ?>>Nothing</option>
                                    <option <?php echo isSelected($cmd->getFormat(), 'memory') ?>>Memory</option>
                                    <option <?php echo isSelected($cmd->getFormat(), 'disc') ?>>Disc</option>
                                    <option <?php echo isSelected($cmd->getFormat(), 'procent') ?>>Procent</option>
                                    <option <?php echo isSelected($cmd->getFormat(), 'date') ?>>Date</option>
                                    <option <?php echo isSelected($cmd->getFormat(), 'seconds') ?>>Seconds</option>
                                    <option <?php echo isSelected($cmd->getFormat(), 'network') ?>>Network</option>
                                </select>
                                <i data-id="<?php echo $cmd->getCommandID()?>" class="delete-command fa fa-close fa-close-red fa-lg pull-right"></i>
                            </td>
                        </tr>
                    <?php } ?>
                </table>
            </div>
        </div>
    </div>
    <div class="row">
        <div class="col-md-3">
            <div class="panel panel-primary">
                <div class="panel-heading">
                    <h4 class="panel-title">Servers</h4>
                    <div class="panel-body">
                        <div class="input-group">
                            <input type="text" class="server-ip form-control" placeholder="Server IP">
                            <span class="input-group-btn">
                                <button class="btn btn-default add-server" type="button">Add Server</button>
                            </span>
                        </div>
                    </div>
                </div>
                <table class="table table-servers table-hover table-bordered table-condensed">
                    <tr>
                        <th width=10%>ID</th>
                        <th width=45%>Namn</th>
                        <th width=45%>IP</th>
                    </tr>
                    <?php foreach($servers as $server) { ?>
                        <tr>
                            <td><?php echo $server['id']?></td>
                            <td contenteditable data-for="server" data-id="<?php echo $server['id']?>" data-target="namn" data-previous="<?php echo $server['namn']?>"><?php echo $server['namn']?></td>
                            <td>
                                <span contenteditable data-for="server" data-id="<?php echo $server['id']?>" data-target="ip" data-previous="<?php echo $server['ip'] ?>"><?php echo $server['ip']?></span>
                                <i data-id="<?php echo $server['id']?>" class="delete-server fa fa-close fa-close-red fa-lg pull-right"></i>
                            </td>
                        </tr>
                    <?php } ?>
                </table>
            </div>
        </div>
    </div>
</div>

<div class="col-md-1"></div>