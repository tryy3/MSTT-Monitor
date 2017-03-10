<?php
    function download() {

    }

    function createClient() {

    }

    function createGroup() {

    }

    function createCommand($db) {
        $stmt = $db->exec("INSERT INTO commands VALUES ()");
        if ($stmt <= 0) {
            return array("error"=>true, "message"=>"Something went wrong when creating a new command.");
        }
        return array("error"=>false, "message"=>$db->lastInsertId());
    }

    function deleteClient() {
        
    }

    function deleteCommand($db, $id) {
        $cmdStmt = $db->prepare("DELETE FROM commands WHERE id=?");
        $groupStmt = $db->prepare("DELETE FROM groups WHERE command_id=?");

        $cmdStmt->execute(array($id));
        if ($cmdStmt->rowCount() <= 0) {
            return array("error"=>true, "message"=>"Something went wrong when deleting the command.");
        }

        $groupStmt->execute(array($id));

        return array("error"=>false, "message"=>"Removed the command, effected ".$groupStmt->rowCount()." groups");
    }

    function deleteGroup() {
        
    }

    function editClient() {
        
    }

    function editCommand($db, $id, $key, $value) {
        if ($key != "namn" && $key != "command" && $key != "description" && $key != "format") {
            return array("error" => true, "message"=>"Invalid key.");
        }
        $stmt = $db->prepare("UPDATE commands SET ".$key."=? WHERE id=?");
        $stmt->execute(array($value, $id));
        if ($stmt->rowCount() <= 0) {
            return array("error" => false, "message" => "Nothing changed.");
        }
        return array("error"=>false, "message" => "Successfully edited the command.");
    }

    function editGroup($db, $id, $key, $value) {
        if ($key != "next_check" && $key != "stop_error") {
            return array("error" => true, "message"=>"Invalid key.");
        }
        if ($key == "next_check" && !is_numeric($value)) {
            return array("error" => true, "messsage" => "Invalid value, value require a number.");
        }
        if ($key == "stop_error") {
            $value = toBool($value);
            if (!is_bool($value))
                return array("error" => true, "message" => "Invalid value, value require a boolean.");
        }
        $stmt = $db->prepare("UPDATE groups SET ".$key."=? WHERE id=?");
        $stmt->execute(array($value, $id));
        if ($stmt->rowCount() <= 0) {
            return array("error" => false, "message" => "Nothing changed.");
        }
        return array("error"=>false, "message" => "Successfully edited the command.");
    }

    function removeCommandGroup($db, $ID) {
        $checkIDStmt = $db->prepare("DELETE FROM groups WHERE id = ?");
        $checkIDStmt->execute(array($ID));

        if ($checkIDStmt->rowCount() <= 0) {
            return array("error"=>true, "message"=>"ID does not exists.");
        }
        return array("error"=>false, "message"=>"Successfully removed the command from the group.");
    }

    function addCommandGroup($db, $groupName, $commandID) {
        $checkGroupStmt = $db->prepare("SELECT command_id FROM groups WHERE group_name = ?");
        $checkCommandStmt = $db->prepare("SELECT namn FROM commands WHERE id = ? LIMIT 1");
        $insertCommandStmt = $db->prepare("INSERT INTO groups(command_id, group_name, next_check, stop_error) VALUES (?,?,?,?)");

        $checkGroupStmt->execute(array($groupName));
        if ($checkGroupStmt->rowCount() <= 0) {
            return array("error"=>true, "message"=>"Group does not exists.");
        }

        while($row = $checkGroupStmt->fetch(PDO::FETCH_ASSOC)) {
            if ($row["command_id"] == $commandID) {
                return array("error"=>true, "message"=>"Group already have access to this command.");
            }
        }

        $checkCommandStmt->execute(array($commandID));
        if ($checkCommandStmt->rowCount() <= 0) {
            return array("error"=>true, "message"=>"Command does not exists.");
        }
        $cmd = $checkCommandStmt->fetch(PDO::FETCH_ASSOC)['namn'];

        $group = array($commandID, $groupName, -1, 0);
        $insertCommandStmt->execute($group);
        $group[0] = $db->lastInsertId();
        $group[1] = $cmd;
        return array("error"=>false, "message"=>$group);
    }
?>