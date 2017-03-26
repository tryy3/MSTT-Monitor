<?php
    namespace MSTT_MONITOR\Classes;
    include_once('../common.php');

    class Command extends \MSTT_MONITOR\Common {
        private $id = -1;
        private $name = "";
        private $description = "";
        private $format = "";

        public function __construct($command) {
            $this->id = $command["id"];
            $this->name = $command["namn"];
            $this->description = $command["description"];
            $this->format = $command["format"];
        }

        public function getID() {
            return $this->id;
        }

        public function getName() {
            return $this->name;
        }

        public function getDescription() {
            return $this->description;
        }

        public function getFormat() {
            return $this->format;
        }
    }
?>