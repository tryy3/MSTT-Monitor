<?php
    namespace MSTT_MONITOR\Classes;
    include_once('../common.php');

    class Check extends \MSTT_MONITOR\Common {
        private $id = -1; // Mysql ID
        private $clientID = -1; // Client ID
        private $commandID = -1; // Command ID
        private $timestamp = -1; // Timestamp (Date)
        private $checked = false; // Checked (bool)
        private $error = false; // Error (bool)
        private $finished = false; // Finished (bool)
        private $response = array(); // Response (JSON)

        public function __construct($check) {
            $this->id = $check["id"];
            $this->clientID = $check["client_id"];
            $this->commandID = $check["command_id"];
            $this->timestamp = strtotime($check["timestamp"]);
            $this->checked = $check["checked"];
            $this->error = $check["error"];
            $this->finished = $check["finished"];
            $this->response = json_decode($check["response"], true);
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
?>