<?php
    namespace MSTT_MONITOR\Utils;

    class Group {
        private $commands = array();
        private $name = "";

        public function __construct($name) {
            $this->name = $name;
        }

        public function addCommand($command) {
            array_push($this->commands, $command);
        }

        public function get($key) {
            switch (strtolower($key)) {
                case 'name':
                case 'namn':
                    return $this->name;
                case 'commands':
                    return $this->commands;
                default:
                    return NULL;
            }
        }

        public function getCommands() {
            return $this->commands;
        }

        public function getName() {
            return $this->name;
        }
    }

    function getGroup($db, $groupName) {
        $group = new Group($groupName);
        $groupStmt = $db->prepare("SELECT * FROM `groups` WHERE `group_name`=?");
        $commandStmt = $db->prepare("SELECT * FROM `commands` WHERE `id`=?");

        $groupStmt->execute(array($groupName));
        while ($row = $groupStmt->fetch(\PDO::FETCH_ASSOC)) {
            $command = new Command();
            $command->addImpl($row);
            
            $commandStmt->execute(array($command->getCommandID()));
            if ($commandStmt->rowCount() > 0) {
                $r = $commandStmt->fetch(\PDO::FETCH_ASSOC);
                $command->addBase($r);
            }
            $group->addCommand($command);
        }
        return $group;
    }

    function getAllGroups($db) {
        $groups = array();
        $stmt = $db->query("SELECT * FROM `groups`");
        while($row = $stmt->fetch(\PDO::FETCH_ASSOC)) {
            if (!isset($groups[$row["group_name"]])) {
                $groups[$row["group_name"]] = new Group($row["group_name"]);
            }
            $command = new Command();
            $command->addImpl($row);
            $groups[$row["group_name"]]->addCommand($command);
        }
        return $groups;
    }
?>