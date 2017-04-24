<?php
    namespace MSTT_Monitor\API;

    class Client extends \MSTT_Monitor\Common {
        /**
         * create()
         *
         * Creates a new client in the database
         *
         * @param int $ip The client IP.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function create($ip) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("INSERT INTO clients(ip) VALUES (?)");
            $stmt->execute(array($ip));
            if ($stmt->rowCount() <= 0) {
                $error->setMessage("Something went wrong when creating a new client.");
                return $error;
            }
            $id = intval($this->db->lastInsertId());

            $error->setBaseURL("/update/client");
            $error->setForm(array( "type" => "update", "id" => $id ));

            $error->setError(false);
            $error->setMessage($id);
            return $error;
        }

        /**
         * delete()
         *
         * Deletes a client from the database.
         *
         * @param int $id The id of the command.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function delete($id) {
            $error = new ErrorAPI();
            $getClientStmt = $this->db->prepare("DELETE FROM clients WHERE id=?");
            $getClientStmt->execute(array($id));
            if ($getClientStmt->rowCount() <= 0) {
                $error->setMessage("Client does not exists.");
                return $error;
            }

            $error->setBaseURL("/update/client");
            $error->setForm(array( "type" => "delete", "id" => intval($id) ));

            $error->setError(false);
            $error->setMessage("Client successfully deleted.");
            return $error;
        }

        /**
         * editIP()
         * 
         * Changes the IP of an existing client.
         * 
         * @param int $id The row id of the client in the database.
         * @param string $value The new ip.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function editIP($id, $value) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE clients SET ip=? WHERE ID=?");
            $stmt->execute(array($value, $id));
            if ($stmt->rowCount() <= 0) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            $error->setBaseURL("/update/client");
            $error->setForm(array( "type" => "update", "id" => intval($id) ));

            $error->setError(false);
            $error->setMessage("Successfully edited the client.");
            return $error;
        }

        /**
         * editName()
         * 
         * Changes the Name of an existing client.
         * 
         * @param int $id The row id of the client in the database.
         * @param string $value The new name.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function editName($id, $value) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE clients SET namn=? WHERE ID=?");
            $stmt->execute(array($value, $id));
            if ($stmt->rowCount() <= 0) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            $error->setError(false);
            $error->setMessage("Successfully edited the client.");
            return $error;
        }

        /**
         * addGroup()
         * 
         * Add a group to the client.
         * 
         * @param string $id The client ID.
         * @param string $groups The group to add to the client.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function addGroup($id, $group) {
            $error = new ErrorAPI();
            if (!$this->groupExists($group)) {
                $error->setMessage("Non existant group name.");
                return $error;
            }

            $groups = $this->getGroups($id);
            $groups = explode(',', $groups);

            if (in_array($group, $groups)) {
                $error->setMessage("Client is already in this group.");
                return $error;
            }

            array_push($groups, $group);
            if (!$this->setGroup($id, trim(implode(',', $groups), ','))) {
                $error->setMessage("Something went wrong in the database.");
                return $error;
            }

            $error->setBaseURL("/update/group");
            $error->setForm(array( "type" => "insert", "id" => intval($id), "group_name" => $group ));

            $error->setError(false);
            $error->setMessage("Updated the clients group successfully.");
            return $error;
        }

        /**
         * delGroup()
         * 
         * Delete a group from a client.
         * 
         * @param string $id The client ID.
         * @param string $groups The group to delete from the client.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function delGroup($id, $group) {
            $error = new ErrorAPI();
            if (!$this->groupExists($group)) {
                $error->setMessage("Non existant group name.");
                return $error;
            }

            $groups = $this->getGroups($id);
            if ($groups === "") {
                $error->setMessage("Can't find a client with this id.");
                return $error;
            }
            $groups = explode(',', $groups);

            $key = array_search($group, $groups);
            if ($key === false) {
                $error->setMessage("Client does not belong to this group.");
                return $error;
            }
            unset($groups[$key]);

            if (!$this->setGroup($id, trim(implode(',', $groups), ','))) {
                $error->setMessage("Something went wrong in the database.");
                return $error;
            }

            $error->setBaseURL("/update/group");
            $error->setForm(array( "type" => "delete", "id" => intval($id), "group_name" => $group ));

            $error->setError(false);
            $error->setMessage("Updated the clients group successfully.");
            return $error;
        }
    }
?>