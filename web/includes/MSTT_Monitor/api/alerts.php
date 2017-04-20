<?php
    namespace MSTT_Monitor\API;

    class Alerts extends \MSTT_Monitor\Common {
        public function create($clientID) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("INSERT INTO alert_options(`client_id`) VALUES (?)");
            $stmt->execute(array($clientID));
            if ($stmt->rowCount() <= 0) {
                $error->setMessage("Something went wrong when creating a new alert option.");
                return $error;
            }
            $id = intval($this->db->lastInsertId());

            //$error->setBaseURL("/update/alert");
            //$error->setForm(array( "type" => "insert", "id" => $id, -1 ));

            $error->setError(false);
            $error->setMessage($id);
            return $error;
        }

        public function delete($id) {
            $error = new ErrorAPI();
            $getClientStmt = $this->db->prepare("DELETE FROM alert_options WHERE id=?");
            $getClientStmt->execute(array($id));
            if ($getClientStmt->rowCount() <= 0) {
                $error->setMessage("Alert option does not exists.");
                return $error;
            }

            //$error->setBaseURL("/update/alert");
            //$error->setForm(array( "type" => "delete", "id" => intval($id) ));

            $error->setError(false);
            $error->setMessage("Alert option successfully deleted.");
            return $error;
        }

        public function editCommand($id, $value) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE alert_options SET command_id=? WHERE ID=?");
            $stmt->execute(array($value, $id));
            if ($stmt->rowCount() <= 0) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            //$error->setBaseURL("/update/alert");
            //$error->setForm(array( "type" => "update", "id" => intval($id) ));

            $error->setError(false);
            $error->setMessage("Successfully edited the alert option.");
            return $error;
        }

        public function editAlert($id, $value) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE alert_options SET alert=? WHERE ID=?");
            $stmt->execute(array($value, $id));
            if ($stmt->rowCount() <= 0) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            //$error->setBaseURL("/update/alert");
            //$error->setForm(array( "type" => "update", "id" => intval($id) ));

            $error->setError(false);
            $error->setMessage("Successfully edited the alert option.");
            return $error;
        }

        public function editValue($id, $value) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE alert_options SET value=? WHERE ID=?");
            $stmt->execute(array($value, $id));
            if ($stmt->rowCount() <= 0) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            //$error->setBaseURL("/update/alert");
            //$error->setForm(array( "type" => "update", "id" => intval($id) ));

            $error->setError(false);
            $error->setMessage("Successfully edited the alert option.");
            return $error;
        }

        public function editCount($id, $value) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE alert_options SET count=? WHERE ID=?");
            $stmt->execute(array($value, $id));
            if ($stmt->rowCount() <= 0) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            //$error->setBaseURL("/update/alert");
            //$error->setForm(array( "type" => "update", "id" => intval($id) ));

            $error->setError(false);
            $error->setMessage("Successfully edited the alert option.");
            return $error;
        }

        public function editDelay($id, $value) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE alert_options SET delay=? WHERE ID=?");
            $stmt->execute(array($value, $id));
            if ($stmt->rowCount() <= 0) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            //$error->setBaseURL("/update/alert");
            //$error->setForm(array( "type" => "update", "id" => intval($id) ));

            $error->setError(false);
            $error->setMessage("Successfully edited the alert option.");
            return $error;
        }

        public function editService($id, $value) {
            $error = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE alert_options SET service=? WHERE ID=?");
            $stmt->execute(array($value, $id));
            if ($stmt->rowCount() <= 0) {
                $error->setMessage("Nothing changed.");
                return $error;
            }

            //$error->setBaseURL("/update/alert");
            //$error->setForm(array( "type" => "update", "id" => intval($id) ));

            $error->setError(false);
            $error->setMessage("Successfully edited the alert option.");
            return $error;
        }
    }
?>