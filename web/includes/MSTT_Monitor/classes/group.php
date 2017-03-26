<?php
    namespace MSTT_MONITOR\Classes;
    include_once('../common.php');

    class Group extends \MSTT_MONITOR\Common {
        private $id = -1;
        private $commandID = -1;
        private $name = "";
        private $nextCheck = -1;
        private $stopError = false;

        public function __construct($group) {
            $this->id = $group["id"];
            $this->commandID = $group["command_id"];
            $this->name = $group["group_name"];
            $this->nextCheck = $group["next_check"];
            $this->stopError = $group["stop_error"];
        }

        public function getID() {
            return $this->id;
        }

        public function getCommandID() {
            return $this->commandID;
        }

        public function getName() {
            return $this->name;
        }

        public function getNextCheck() {
            return $this->nextCheck;
        }

        public function getStopError() {
            return $this->stopError;
        }
    }
?>