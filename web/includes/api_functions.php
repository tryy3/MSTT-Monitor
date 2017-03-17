<?php
    /* Download API */
    function download() {

    }

    /* Client API */
    function createClient($db, $ip) {
        $stmt = $db->prepare("INSERT INTO clients(ip) VALUES (?)");
        $stmt->execute(array($ip));
        if ($stmt->rowCount() <= 0) {
            return array("error" => true, "message" => "Something went wrong when creating a new client.");
        }
        $id = intval($db->lastInsertId());

        $out = sendServerRequest($db, "/update/client", array(
            "type" => "update",
            "id" => $id
        ));

        if ($out["error"]) {
            return $out;
        }
        return array("error" => false, "message" => $id);
    }

    function deleteClient($db, $id) {
        $getClientStmt = $db->prepare("DELETE FROM clients WHERE id=?");
        $getClientStmt->execute(array($id));
        if ($getClientStmt->rowCount() <= 0) {
            return array("error" => true, "message" => "Client does not exists.");
        }

        $out = sendServerRequest($db, "/update/client", array(
            "type" => "delete",
            "id" => intval($id)
        ));

        if ($out["error"]) {
            return $out;
        }
        return array("error" => false, "message" => "Client successfully deleted.");
    }

    function editClient($db, $id, $key, $value) {
        if ($key != "ip" && $key != "namn") {
            return array("error" => true, "message"=>"Invalid key.");
        }
        $stmt = $db->prepare("UPDATE clients SET ".$key."=? WHERE ID =?");
        $stmt->execute(array($value, $id));
        if ($stmt->rowCount() <= 0) {
            return array("error" => false, "message" => "Nothing changed.");
        }

        if ($key == "ip") {
            $out = sendServerRequest($db, "/update/client", array(
                "type" => "update",
                "id" => intval($id)
            ));

            if ($out["error"]) {
                return $out;
            }
        }
        
        return array("error" => false, "message" => "Successfully edited the client.");
    }

    function editClientGroup($db, $id, $group, $type) {
        if ($type != "add" && $type != "remove") {
            return array("error" => true, "message" => "Invalid type.");
        }

        $groupStmt = $db->prepare("SELECT * FROM groups WHERE group_name=?");
        $groupStmt->execute(array($group));
        if ($groupStmt->rowCount()<=0) {
            return array("error" => true, "message" => "Non existant group name.");
        }

        $clientStmt = $db->prepare("SELECT group_names FROM clients WHERE id=?");
        $clientStmt->execute(array($id));
        if ($clientStmt->rowCount()<=0) {
            return array("error" => true, "message" => "Can't find a client with this id.");
        }

        $groups = $clientStmt->fetch(PDO::FETCH_ASSOC)["group_names"];
        $groups = explode(',', $groups);

        switch(strtolower($type)) {
            case "remove":
                $key = array_search($group, $groups);
                if ($key === false) {
                    return array("error" => true, "message" => "Client does not belong to this group.");
                }
                unset($groups[$key]);
                break;
            case "add":
                if (in_array($group, $groups)) {
                    return array("error" => true, "message" => "Client is already in this group.");
                }
                array_push($groups, $group);
                break;
        }

        $type = ($type=="add") ? "insert" : "delete";
        $out = sendServerRequest($db, "/update/group", array(
            "type" => $type,
            "clientid" => intval($id),
            "name" => $group
        ));

        if ($out["error"]) {
            return $out;
        }

        $clientUpdateStmt = $db->prepare("UPDATE clients SET group_names=? WHERE id=?");
        $clientUpdateStmt->execute(array(trim(implode(",",$groups), ","), $id));
        return array("error" => false, "message" => "Updated the clients group successfully.");
    }

    /* Group API */
    function groupExists($db, $group) {
        $stmt = $db->prepare("SELECT * FROM groups WHERE group_name=?");
        $stmt->execute(array($group));
        if ($stmt->rowCount() <= 0) {
            return array("error" => false, "message" => false);
        }
        return array("error" => false, "message" => true);
    }

    function deleteGroup($db, $group) {
        $deleteGroupStmt = $db->prepare("DELETE FROM groups WHERE group_name=?");
        $deleteGroupStmt->execute(array($group));
        if ($deleteGroupStmt->rowCount() <= 0) {
            return array("error" => true, "message" => "Group does not exists.");
        }

        $updateClientStmt = $db->prepare("UPDATE clients SET group_names=? WHERE id=?");
        $getClientStmt = $db->query("SELECT id, group_names FROM clients");

        while($row = $getClientStmt->fetch(PDO::FETCH_ASSOC)) {
            $groups = explode(",", $row["group_names"]);
            $key = array_search($group, $groups);
            if ($key === false) {
                continue;
            }
            unset($groups[$key]);
            $updateClientStmt->execute(array(trim(implode(",",$groups), ","), $row["id"]));
        }

        $out = sendServerRequest($db, "/update/client", array(
            "type" => "delete",
            "name" => $group
        ));

        if ($out["error"]) {
            return $out;
        }

        return array("error" => false, "message" => "Group deleted.");
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

        $out = sendServerRequest($db, "/update/group", array(
            "type" => "update",
            "name" => $id
        ));

        if ($out["error"]) {
            return $out;
        }
        return array("error"=>false, "message" => "Successfully edited the command.");
    }

    function removeCommandGroup($db, $id, $group) {
        $getCommandIDStmt = $db->prepare("SELECT command_id FROM groups WHERE id = ?");
        $getCommandIDStmt->execute(array($id));
        $command = $getCommandIDStmt->Fetch(PDO::FETCH_ASSOC);
        
        $checkIDStmt = $db->prepare("DELETE FROM groups WHERE id = ?");
        $checkIDStmt->execute(array($id));

        if ($checkIDStmt->rowCount() <= 0) {
            return array("error"=>true, "message"=>"ID does not exists.");
        }

        $out = sendServerRequest($db, "/update/command", array(
            "type" => "delete",
            "id" => intval($command["command_id"]),
            "group" => $group
        ));

        if ($out["error"]) {
            return $out;
        }

        return array("error"=>false, "message"=>"Successfully removed the command from the group.");
    }

    function addCommandGroup($db, $groupName, $commandID) {
        $checkGroupStmt = $db->prepare("SELECT command_id FROM groups WHERE group_name = ?");
        $checkCommandStmt = $db->prepare("SELECT namn FROM commands WHERE id = ? LIMIT 1");
        $insertCommandStmt = $db->prepare("INSERT INTO groups(command_id, group_name, next_check, stop_error) VALUES (?,?,?,?)");

        $checkGroupStmt->execute(array($groupName));
        if ($checkGroupStmt->rowCount() > 0) {
            while($row = $checkGroupStmt->fetch(PDO::FETCH_ASSOC)) {
                if ($row["command_id"] == $commandID) {
                    return array("error"=>true, "message"=>"Group already have access to this command.");
                }
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

        $out = sendServerRequest($db, "/update/command", array(
            "type" => "insert",
            "id" => intval($commandID),
            "group" => $groupName
        ));

        if ($out["error"]) {
            return $out;
        }

        return array("error"=>false, "message"=>$group);
    }

    /* Command API */
    function createCommand($db) {
        $stmt = $db->exec("INSERT INTO commands VALUES ()");
        if ($stmt <= 0) {
            return array("error"=>true, "message"=>"Something went wrong when creating a new command.");
        }
        return array("error"=>false, "message"=>$db->lastInsertId());
    }

    function deleteCommand($db, $id) {
        $cmdStmt = $db->prepare("DELETE FROM commands WHERE id=?");
        $groupStmt = $db->prepare("DELETE FROM groups WHERE command_id=?");

        $cmdStmt->execute(array($id));
        if ($cmdStmt->rowCount() <= 0) {
            return array("error"=>true, "message"=>"Something went wrong when deleting the command.");
        }

        $groupStmt->execute(array($id));

        $out = sendServerRequest($db, "/update/command", array(
            "type" => "delete",
            "id" => intval($id)
        ));

        if ($out["error"]) {
            return $out;
        }

        return array("error"=>false, "message"=>"Removed the command, effected ".$groupStmt->rowCount()." groups");
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

        if ($key == "command") {
            $out = sendServerRequest($db, "/update/command", array(
                "type" => "update",
                "id" => intval($id)
            ));

            if ($out["error"]) {
                return $out;
            }
        }
        return array("error"=>false, "message" => "Successfully edited the command.");
    }

    /* Server API */
    function addServer($db, $ip) {
        $stmt = $db->prepare("INSERT INTO servers(ip) VALUES (?)");
        $stmt->execute(array($ip));
        if ($stmt->rowCount() <= 0) {
            return array("error" => true, "message" => "Something went wrong when adding a server.");
        }
        return array("error" => false, "message" => $db->lastInsertId());
    }

    function delServer($db, $id) {
        $stmt = $db->prepare("DELETE FROM servers WHERE id=?");
        $stmt->execute(array($id));
        if($stmt->rowCount() <= 0) {
            return array("error" => true, "message" => "Server ID does not exists.");
        }
        return array("error" => false, "message" => "Server removed.");
    }

    function editServer($db, $id, $key, $value) {
        if ($key != "ip" && $key != "namn") {
            return array("error" => true, "message"=>"Invalid key.");
        }
        $stmt = $db->prepare("UPDATE servers SET ".$key."=? WHERE id=?");
        $stmt->execute(array($value, $id));
        if ($stmt->rowCount() <= 0) {
            return array("error" => false, "message" => "Nothing changed.");
        }
        return array("error"=>false, "message" => "Successfully edited the server.");
    }

    /* General functions */
    function sendServerRequest($db, $base, $form) {
        $servers = getServers($db);
        $data = json_encode($form);

        $out = array("error" => false, "message" => "Sent a server request to all servers.");

        foreach($servers as $server) {
            $url = "http://".$server["ip"].$base;
            $ch = curl_init($url);
            curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "POST");
            curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
            curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
            curl_setopt($ch, CURLOPT_HTTPHEADER, array(
                "Content-Type: application/json",
                "Content-Length: " . strlen($data)
            ));

            $resp = curl_exec($ch);
            $status = curl_getinfo($ch, CURLINFO_HTTP_CODE);
            if ( $status != 200 ) {
                error_log("Error: call to URL $url failed with status $status, response $resp, curl_error " . curl_error($ch) . ", curl_errno " . curl_errno($ch));
                return array("error" => true, "message" => "One or more servers failed, check logs.");
            }

            curl_close($ch);
            $response = json_decode($resp, true);
            error_log("Debug: call to URL $url, response has a error ".$response["message"]);
            if ($response["error"]) {
                return array("error" => true, "message" => "One or more servers failed, check logs.");
            }
        }
        return $out;
    }
?>