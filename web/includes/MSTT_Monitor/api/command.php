<?php
    namespace MSTT_Monitor\API;

    class Command extends \MSTT_Monitor\Common {
        /**
         * create()
         *
         * Creates a new command in the database
         *
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function create() {
            $error = new ErrorAPI();
            $stmt = $this->db->exec("INSERT INTO commands VALUES ()");
            if ($stmt <= 0) {
                $error->setMessage("Something went wrong when creating a new command.");
                return $error;
            }
            $error->setError(false);
            $error->setMessage($this->db->lastInsertId());
            return $error;
        }

        /**
         * delete()
         *
         * Deletes a command from the database.
         *
         * @param int $id The id of the command.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function delete($id) {
            $error = new ErrorAPI();
            $cmdStmt = $this->db->prepare("DELETE FROM commands WHERE id=?");
            $groupStmt = $this->db->prepare("DELETE FROM groups WHERE command_id=?");

            $cmdStmt->execute(array($id));
            if ($cmdStmt->rowCount() <= 0) {
                $error->setMessage("Something went wrong when deleting the command.");
                return $error;
            }

            $groupStmt->execute(array($id));

            $error->setBaseURL("/update/command");
            $error->setForm(array( "type" => "update", "command_id" => intval($id) ));

            $error->setError(false);
            $error->setMessage("Removed the command, effected ".$groupStmt->rowCount()." groups");
            return $error;
        }

        /**
         * editCommand()
         * 
         * Edit the command value of a command.
         * 
         * @param int $id The row id of the command in the database.
         * @param string $value The new command value.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function editCommand($id, $value) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE commands SET command=? WHERE id=?");
            if (!$this->updateCommand($stmt, array($value, $id))) {
                $error->setMessage("Nothing changed.");
                return $error;
            }
            
            $error->setBaseURL("/update/command");
            $error->setForm(array( "type" => "update", "command_id" => intval($id) ));

            $error->setError(false);
            $error->setMessage("Successfully edited the command.");
            return $error;
        }

        /**
         * editName()
         * 
         * Edit a command name.
         * 
         * @param int $id The row id of the command in the database.
         * @param string $value The new command name.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function editName($id, $value) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE commands SET namn=? WHERE id=?");
            if (!$this->updateCommand($stmt, array($value, $id))) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            $error->setError(false);
            $error->setMessage("Successfully edited the command.");
            return $error;
        }

        /**
         * editDescription()
         * 
         * Edit a commands description.
         * 
         * @param int $id The row id of the command in the database.
         * @param string $value The new command description.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function editDescription($id, $value) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE commands SET description=? WHERE id=?");
            if (!$this->updateCommand($stmt, array($value, $id))) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            $error->setError(false);
            $error->setMessage("Successfully edited the command.");
            return $error;
        }

        /**
         * editName()
         * 
         * Edit a commands format.
         * 
         * @param int $id The row id of the command in the database.
         * @param string $value The new command format.
         * @return ErrorAPI An ErrorAPI instance containing information if the function got any errors or not.
         */
        public function editFormat($id, $value) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE commands SET format=? WHERE id=?");
            if (!$this->updateCommand($stmt, array($value, $id))) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            $error->setError(false);
            $error->setMessage("Successfully edited the command.");
            return $error;
        }

        /**
         * updateCommand()
         * 
         * Update a command statement.
         * 
         * @param PDO::Statement $stmt The command statement to run when updating a command.
         * @param array $value The array to execute.
         * @return bool A boolean if the command was ran successfully or not.
         */
        private function updateCommand($stmt, $value) {
            $stmt->execute($value);
            if ($stmt->rowCount() <= 0) {
                return false;
            }
            return true;
        }
    }
?>