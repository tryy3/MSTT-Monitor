<?php
    namespace MSTT_MONITOR\Utils;

    class Command {
        private $id = -1;
        private $commandID = -1;
        private $nextCheck = -1;
        private $stopError = false;
        private $command = "";
        private $name = "";
        private $description = "";
        private $format = "";

        public function addImpl($command) {
            if (isset($command["id"])) {
                $this->id = $command["id"];
            }

            if (isset($command["command_id"])) {
                $this->commandID = $command["command_id"];
            }

            if (isset($command["next_check"])) {
                $this->nextCheck = $command["next_check"];
            }

            if (isset($command["stop_error"])) {
                $this->stopError = $command["stop_error"];
            }
        }

        public function addBase($command) {
            if (isset($command["id"])) {
                $this->commandID = $command["id"];
            }

            if (isset($command["command"])) {
                $this->command = $command["command"];
            }
            
            if (isset($command["namn"])) {
                $this->name = $command["namn"];
            }
            
            if (isset($command["description"])) {
                $this->description = $command["description"];
            }
            
            if (isset($command["format"])) {
                $this->format = $command["format"];
            }
        }

        public function get($key) {
            switch (strtolower($key)) {
                case 'id':
                    return $this->id;
                case 'commandid':
                case 'command_id':
                    return $this->commandID;
                case 'nextcheck':
                case 'next_check':
                    return $this->nextCheck;
                case 'stoperror':
                case 'stop_error':
                    return $this->stopError;
                default:
                    return NULL;
            }
        }

        public function getID() {
            return $this->id;
        }

        public function getCommandID() {
            return $this->commandID;
        }

        public function getNextCheck() {
            return $this->nextCheck;
        }

        public function getStopError() {
            return $this->stopError;
        }

        public function getCommand() {
            return $this->command;
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

    function getAllCommands($db) {
        $commands = array();
        $stmt = $db->query("SELECT * FROM `commands`");
        if (!$stmt || $stmt->rowCount() <= 0) {
            return $commands;
        }
        while ($row = $stmt->fetch(\PDO::FETCH_ASSOC)) {
            $command = new Command();
            $command->addBase($row);
            array_push($commands, $command);
        }
        return $commands;
    }
?>