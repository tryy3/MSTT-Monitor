<?php
    namespace MSTT_MONITOR\Utils;

    class AlertOptions {
        private $id = -1; // Mysql ID
        private $clientID = -1; // Client ID
        private $commandID = -1; // Command ID
        private $alert = ""; // Alert function name
        private $value = ""; // Alert value
        private $count = -1; // Alert count
        private $delay = -1; // Alert delay
        private $service = ""; // Alert services

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
            
            if (isset($check["alert"])) {
                $this->alert = $check["alert"];
            }
            
            if (isset($check["value"])) {
                $this->value = $check["value"];
            }
            
            if (isset($check["count"])) {
                $this->count = $check["count"];
            }
            
            if (isset($check["delay"])) {
                $this->delay = $check["delay"];
            }
            
            if (isset($check["service"])) {
                $this->service = $check["service"];
            }
        }

        public function getID() {
            return $this->id;
        }

        public function getClientID() {
            return $this->clientID;
        }

        public function getAlert() {
            return $this->alert;
        }

        public function getValue() {
            return $this->value;
        }

        public function getCount() {
            return $this->count;
        }

        public function getDelay() {
            return $this->delay;
        }

        public function getCommand() {
            return $this->commandID;
        }

        public function getDelayFormat() {
            return sprintf("%02d:%02d:%02d", floor($this->delay/3600), ($this->delay/60)%60, $this->delay%60);
        }

        public function getService() {
            return $this->service;
        }
    }
    
    function getAlertOptions($db, $clientID) {
        $stmt = $db->prepare("SELECT * FROM `alert_options` WHERE `client_id`=?");
        $stmt->execute(array($clientID));

        $alerts = array();
        if ($stmt->rowCount() <= 0) {
            return $alerts;
        }

        while($row = $stmt->fetch(\PDO::FETCH_ASSOC)) {
            array_push($alerts, new AlertOptions($row)); 
        }
        return $alerts;
    }

    function getAllAlertOptions($db) {
        $alertOptions = array();
        $stmt = $db->query("SELECT * FROM `alert_options`");

        while ($row = $stmt->fetch(\PDO::FETCH_ASSOC)) {
            array_push($alertOptions, new AlertOptions($row));
        }
        return $alertOptions;
    }
?>