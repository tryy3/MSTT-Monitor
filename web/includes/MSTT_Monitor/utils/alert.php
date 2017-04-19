<?php
    namespace MSTT_MONITOR\Utils;
    include_once(__DIR__.'/../../function.php');

    class Alert {
        private $id = -1; // Mysql ID
        private $clientID = -1; // Client ID
        private $alertID = -1; // Alert ID
        private $value = ""; // Alert value
        private $timestamp = -1; // Alert timestamp

        public function __construct($check) {
            if (isset($check["id"])) {
                $this->id = $check["id"];
            }
            
            if (isset($check["client_id"])) {
                $this->clientID = $check["client_id"];
            }
            
            if (isset($check["alert_id"])) {
                $this->alertID = $check["alert_id"];
            }
            
            if (isset($check["value"])) {
                $this->value = $check["value"];
            }
            
            if (isset($check["timestamp"])) {
                $this->timestamp = strtotime($check["timestamp"]);
            }
        }

        public function getID() {
            return $this->id;
        }

        public function getClientID() {
            return $this->clientID;
        }

        public function getAlertID() {
            return $this->alertID;
        }

        public function getValue() {
            return $this->value;
        }

        public function getTimestamp() {
            return $this->timestamp;
        }

        public function getDate() {
            return date("Y/m/d H:i:s", $this->timestamp);
        }
    }
    
    function getAlerts($db, $clientID, $alertID, $limit = 10) {
        $stmt = $db->prepare("SELECT * FROM `alert` WHERE `client_id`=? AND `alert_id`=?");
        $stmt->execute(array($clientID, $alertID));

        $alerts = array();
        if ($stmt->rowCount() <= 0) {
            return $alerts;
        }

        for ($i = 0; $i < $limit; $i++) {
            $result = $stmt->fetch(\PDO::FETCH_ASSOC);
            if (!$result) {
                return $checks;
            }

            $alerts = array_push(new Alert($result));
        }
        return $alerts;
    }
?>