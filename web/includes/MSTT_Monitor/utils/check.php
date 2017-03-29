<?php
    namespace MSTT_MONITOR\Utils;
    include_once(__DIR__.'/../../function.php');

    class Check {
        private $id = -1; // Mysql ID
        private $clientID = -1; // Client ID
        private $commandID = -1; // Command ID
        private $timestamp = -1; // Timestamp (Date)
        private $checked = false; // Checked (bool)
        private $error = false; // Error (bool)
        private $finished = false; // Finished (bool)
        private $response = array(); // Response (JSON)

        public function __construct($check) {
            if (isset($check["id"])) {
                $this->id = $check["id"];
            }
            
            if (isset($check["client_id"])) {
                $this->clientID = $check["client_id"];
            }
            
            if (isset($check["command_id"])) {
                $this->commandID = $check["command_id"];
            }
            
            if (isset($check["timestamp"])) {
                $this->timestamp = strtotime($check["timestamp"]);
            }
            
            if (isset($check["checked"])) {
                $this->checked = \toBool($check["checked"]);
            }
            
            if (isset($check["error"])) {
                $this->error = \toBool($check["error"]);
            }
            
            if (isset($check["finished"])) {
                $this->finished = \toBool($check["finished"]);
            }
            
            if (isset($check["response"])) {
                $this->response = json_decode(utf8_encode($check["response"]), true);
            }
        }

        public function getID() {
            return $this->id;
        }

        public function getClientID() {
            return $this->clientID;
        }

        public function getCommandID() {
            return $this->commandID;
        }

        public function getTimestamp() {
            return $this->timestamp;
        }

        public function getChecked() {
            return $this->checked;
        }

        public function getError() {
            return $this->error;
        }

        public function getFinished() {
            return $this->finished;
        }

        public function getResponse() {
            return $this->response;
        }
    }

    function getChecks($db, $clientID, $commandID = -1, $limit = 1) {
        if ($commandID == -1) {
            $stmt = $db->prepare("SELECT * FROM `checks` WHERE `client_id`=? ORDER BY timestamp DESC");
            $stmt->execute(array($clientID));
        } else {
            $stmt = $db->prepare("SELECT * FROM `checks` WHERE `client_id`=? AND `command_id`=? ORDER BY timestamp DESC");
            $stmt->execute(array($clientID, $commandID));
        }

        $checks = array();
        if ($stmt->rowCount() <= 0) {
            return $checks;
        }

        for ($i = 0; $i < $limit; $i++) {
            $result = $stmt->fetch(\PDO::FETCH_ASSOC);
            if (!$result) {
                return $checks;
            }
            array_push($checks, new Check($result));
        }
        return $checks;
    }

    function getAllChecks($db, $clientID, $limit = 1) {
        $stmt = $db->prepare("SELECT `command_id` FROM `checks` WHERE `client_id`=? GROUP BY `command_id`");
        $stmt->execute(array($clientID));
        if ($stmt->rowCount() <= 0) {
            return array();
        }
        $checks = array();
        while($row = $stmt->fetch(\PDO::FETCH_ASSOC)) {
            $checks = array_merge($checks, getChecks($db, $clientID, $row["command_id"], $limit));
        }
        return $checks;
    }
?>