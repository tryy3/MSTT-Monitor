<?php
    namespace MSTT_Monitor\API;

    class Group extends \MSTT_Monitor\Common {
        /**
         * exists()
         *
         * Checks if a group actually exists or not.
         *
         * @param string $group The group to check exists or not.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function exists($group) {
            $error = new ErrorAPI(false);
            $error->setMessage($this->groupExists($group));
            return $error;
        }

        /**
         * delete()
         *
         * Deletes a group from the database.
         * This will also remove the group from any clients that is in the group.
         *
         * @param string $group The group name to remove.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function delete($group) {
            $error = new ErrorAPI();
            $deleteGroupStmt = $this->db->prepare("DELETE FROM groups WHERE group_name=?");
            $deleteGroupStmt->execute(array($group));
            if ($deleteGroupStmt->rowCount() <= 0) {
                $error->setMessage("Group does not exists.");
                return $error;
            }

            $updateClientStmt = $this->db->prepare("UPDATE clients SET group_names=? WHERE id=?");

            $getClientStmt = $db->query("SELECT id, group_names FROM clients");
            while($row = $getClientStmt->fetch(PDO::FETCH_ASSOC)) {
                $groups = explode(",", $row["group_names"]);
                $key = array_search($group, $groups);
                if ($key === false) {
                    continue;
                }
                unset($groups[$key]);
                $updateClientStmt->execute(array(trim(implode(",",$groups),","),$row["id"]));
            }

            $error->setBaseURL("/update/client");
            $error->setForm(array("type" => "delete", "name" => $group));

            $error->setError(false);
            $error->setMessage("Group deleted.");
            return $error;
        }
        
        /**
         * editNextCheck()
         *
         * Change when the next check for a group command is gonna happen.
         *
         * @param int $id The row id of the group in the database.
         * @param int $value The value to set the next check to.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function editNextCheck($id, $value) {
            $error = new ErrorAPI();
            if (!is_numeric($value)) {
                $error->setMessage("Invalid value, value require a number.");
                return $error;
            }
            $stmt = $this->db->prepare("UPDATE groups SET next_check=? WHERE id=?");
            $stmt->execute(array($value, $id));
            if ($stmt->rowCount() <= 0) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            $error->setBaseURL("/update/group");
            $error->setForm(array("type"=>"update","name"=>$id));

            $error->setError(false);
            $error->setMessage("Successfully edited the command.");
            return $error;
        }
        
        /**
         * editStopError()
         *
         * Change when the next check for a group command is gonna happen.
         *
         * @param int $id The row id of the group in the database.
         * @param bool $value The value to set the stop error to.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function editStopError($id, $value) {
            $error = new ErrorAPI();
            $val = toBool($value);
            if (!is_bool($val)) {
                $error->setMessage("Invalid value, value require a boolean.");
                return $error;
            }

            $stmt = $this->db->prepare("UPDATE groups SET stop_error=? WHERE id=?");
            $stmt->execute(array($val, $id));
            if ($stmt->rowCount() <= 0) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            $error->setBaseURL("/update/group");
            $error->setForm(array("type"=>"update","name"=>$id));

            $error->setError(false);
            $error->setMessage("Successfully edited the command.");
            return $error;
        }

        public function removeCommand($id, $group) {
            $errors = new ErrorAPI();

            $getCommandStmt = $this->db->prepare("SELECT command_id FROM Groups WHERE id=?");
            $getCommandStmt->execute(array($id));
            $command = $getCommandStmt->fetch(PDO::FETCH_ASSOC);

            $checkIDStmt = $this->db->prepare("DELETE FROM groups WHERE id=?");
            $checkIDStmt->execute(array($id));

            if($checkIDStmt->rowCount() <= 0) {
                $errors->setMessage("ID does not exists.");
                return $errors;
            }

            $errors->setBaseURL("/update/command");
            $errors->setForm(array("type"=>"delete", "id"=>intval($command["command_id"]), "group"=>$group));

            $errors->setError(false);
            $errors->setMessage("Successfully removed the command from the group.");
            return $errors;
        }

        public function addCommand($groupName, $commandID) {
            $errors = new ErrorAPI();

            $checkGroupStmt = $this->db->prepare("SELECT command_id FROM groups WHERE group_name = ?");
            $checkCommandStmt = $this->db->prepare("SELECT namn FROM commands WHERE id = ? LIMIT 1");
            $insertCommandStmt = $this->db->prepare("INSERT INTO groups(command_id, group_name, next_check, stop_error) VALUES (?,?,?,?)");

            $checkGroupStmt->execute(array($group));
            if ($checkGroupStmt->rowCount() > 0) {
                while($row = $checkGroupStmt->fetch(PDO::FETCH_ASSOC)) {
                    if ($row["command_id"] == $commandID) {
                        $errors->setMessage("Group already have access to this command.");
                        return $errors;
                    }
                }
            }

            $checkCommandStmt->execute(array($commandID));
            if ($checkCommandStmt->rowCount() <= 0) {
                $errors->setMessage("Command does not exists.");
                return $errors;
            }
            $cmd = $checkCommandStmt->fetch(PDO::FETCH_ASSOC)["namn"];

            $insertCommandStmt->execute(array($commandID, $groupName, -1, 0));

            $errors->setBaseURL("/update/command");
            $errors->setForm(array("type"=>"insert", "id"=>intval($commandID), "group"=>$groupName));

            $errors->setError(false);
            $errors->setMessage(array($this->db->lastInsertId(), $cmd, -1, 0));
            return $errors;
        }
    }
?>