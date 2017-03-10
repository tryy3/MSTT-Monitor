<?php 
    var_dump(is_bool('0'));
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

    try {
        $groups = getGroups($monitorDB);
        $commands = getCommands($monitorDB);
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
                            <input type="text" class="form-control" placeholder="Group name">
                            <span class="input-group-btn">
                                <button class="btn btn-default" type="button">Add group</button>
                            </span>
                        </div>
                    </div>
                </a>
                <?php foreach ($groups as $group) { ?>
                    <button type="button" class="list-group-item checks" data-target="<?php echo $group['group_name']?>">
                        <h4 class="list-group-item-heading">
                            <?php echo $group['group_name'] ?>
                            <i class="delete-group fa fa-close fa-close-red fa-lg" style="padding-left:60%;"></i>
                        </h4>
                    </button>
                <?php } ?>
            </div>
        </div>
        <div class="col-md-4">
            <?php foreach($groups as $group) { ?>
                <div class="panel panel-primary checks-item drop-group" data-check="<?php echo $group['group_name']?>" style="display:none;">
                    <div class="panel-heading"> 
                        <h3 class="panel-title"><?php echo $group['group_name'] ?></h3>
                    </div>
                    
                    <table class="table table-groups table-hover table-bordered table-condensed">
                        <tr>
                            <th width=5%>ID</th>
                            <th width=45%>Namn</th>
                            <th width=25%>Nästa Check</th>
                            <th width=25%>Stop Error</th>
                        </tr>
                        <?php foreach($group['commands'] as $cmd) { ?>
                            <tr class="drag">
                                <td><?php echo $cmd['id'] ?></td>
                                <td><?php echo $commands[$cmd['command_id']]['namn']?></td>
                                <td contenteditable data-previous="<?php echo $cmd['next_check']?>" data-for="group" data-target="next_check" data-id="<?php echo $cmd['id']?>"><?php echo $cmd['next_check']?></td>
                                <td>
                                    <select data-for="group" data-id="<?php echo $cmd['id']?>" data-target="stop_error" class="form-control" style="width:80%; display:inline">
                                        <option <?php echo isSelected(toBool($cmd['stop_error']), true) ?>>True</option>
                                        <option <?php echo isSelected(toBool($cmd['stop_error']), false) ?>>False</option>
                                    </select>
                                    <i data-id="<?php echo $cmd['id']?>" class="remove-command-group fa fa-close fa-close-red fa-lg"></i>
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
                        <tr class="drag">
                            <td><?php echo $cmd['id']?></td>
                            <td contenteditable data-for="cmd" data-id="<?php echo $cmd['id']?>" data-target="namn" data-previous="<?php echo $cmd['namn']?>"><?php echo $cmd['namn']?></td>
                            <td contenteditable data-for="cmd" data-id="<?php echo $cmd['id']?>" data-target="command" data-previous="<?php echo $cmd['command']?>"><?php echo $cmd['command']?></td>
                            <td contenteditable data-for="cmd" data-id="<?php echo $cmd['id']?>" data-target="description" data-previous="<?php echo $cmd['description']?>"><?php echo $cmd['description']?></td>
                            <td>
                                <select data-for="cmd" data-id="<?php echo $cmd['id']?>" data-target="format" class="form-control" style="width:80%; display:inline">
                                    <option <?php echo isSelected($cmd['format'], '') ?>>Nothing</option>
                                    <option <?php echo isSelected($cmd['format'], 'bytes') ?>>Bytes</option>
                                    <option <?php echo isSelected($cmd['format'], 'bits') ?>>Bits</option>
                                </select>
                                <i data-id="<?php echo $cmd['id']?>" class="delete-command fa fa-close fa-close-red fa-lg"></i>
                            </td>
                        </tr>
                    <?php } ?>
                </table>
            </div>
        </div>
    </div>
</div>

<div class="col-md-1"></div>