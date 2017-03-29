<?php
    namespace MSTT_MONITOR\Utils;
    
    class Client {
        private $id = -1;
        private $groups = array();
        private $checks = array();
        private $latest = -1;
        private $ip = "";
        private $name = "";

        public function __construct($client) {
            if (isset($client["id"])) {
                $this->id = $client["id"];
            }
            if (isset($client["ip"])) {
                $this->ip = $client["ip"];
            }
            if (isset($client["namn"])) {
                $this->name = $client["namn"];
            }
        }

        public function get($key) {
            switch (strtolower($key)) {
                case 'id':
                    return $this->id;
                case 'groups':
                    return $this->groups;
                case 'groupnames':
                case 'group_names':
                    return $this->getGroupNames();
                case 'checks':
                    return $this->checks;
                case 'name':
                case 'namn':
                    return $this->name;
                case 'latest':
                    return $this->getLatest();
                case 'ip':
                    return $this->ip;
                default:
                    return NULL;
            }
        }

        public function addCheck($check) {
            if ($check->getTimestamp() != -1 && $check->getTimestamp() > $this->latest) {
                $this->latest = $check->getTimestamp();
            }
            array_push($this->checks, $check);
        }

        public function addGroup($group) {
            array_push($this->groups, $group);
        }

        public function getGroupNames() {
            $out = "";
            foreach ($this->groups as $group) {
                $out .= $group->getName().", ";
            }
            return trim($out, " ,");
        }

        public function getGroupByID($id) {
            foreach ($this->groups as $group) {
                if ($group->getID() == $id) {
                    return $group;
                }
            }
            return NULL;
        }

        public function getGroup($id) {
            if (isset($this->groups[$id])) {
                return $this->groups[$id];
            }
            return NULL;
        }

        public function getGroups() {
            return $this->groups;
        }

        public function getCommand($id) {
            foreach ($this->groups as $group) {
                foreach ($group->getCommands() as $command) {
                    if ($command->getCommandID() == $id) {
                        return $command;
                    }
                }
            }
            return NULL;
        }

        public function getCommands() {
            $commands = array();
            foreach ($this->groups as $group) {
                foreach ($group->getCommands() as $command) {
                    array_push($commands, $command);
                }
            }
            return $commands;
        }

        public function getCheckByID($id) {
            foreach ($this->checks as $check) {
                if ($check->getID() == $id) {
                    return $check;
                }
            }
            return NULL;
        }

        public function getChecksByCommandID($commandID) {
            $checks = array();
            foreach ($this->checks as $check) {
                if ($check->getCommandID() == $commandID) {
                    array_push($checks, $check);
                }
            }
            return $checks;
        }

        public function getCheck($id) {
            if (isset($this->checks[$id])) {
                return $this->checks[$id];
            }
            return NULL;
        }

        public function getID() {
            return $this->id;
        }

        public function getLatest() {
            return date("Y/m/d H:i:s", $this->latest);
        }

        public function getName() {
            return $this->name;
        }

        public function getIP() {
            return $this->ip;
        }
    }

    function getAllClients($db, $checks = 10) {
        $clients = array();
        $groups = getAllGroups($db);

        $stmt = $db->query('SELECT * FROM clients');
        if (!$stmt) return $clients;
        while($row = $stmt->fetch(\PDO::FETCH_ASSOC)) {
            $client = new Client($row);
            foreach (explode(",", $row["group_names"]) as $groupName) {
                $client->addGroup($groups[$groupName]);

                foreach ($groups[$groupName]->getCommands() as $command) {
                    foreach (getChecks($db, $client->getID(), $command->getCommandID(), $checks) as $check) {
                        $client->addCheck($check);
                    }
                }
            }
            array_push($clients, $client);
		}
		return $clients;
    }

    function getClient($db, $id) {
        $stmt = $db->prepare("SELECT * FROM clients WHERE `id`=?");
        $stmt->execute(array($id));
        if ($stmt->rowCount() <= 0) {
            return NULL;
        }
        $row = $stmt->fetch(\PDO::FETCH_ASSOC);
        $client = new Client($row);

        foreach (explode(",", $row["group_names"]) as $groupName) {
            $group = getGroup($db, $groupName);
            $client->addGroup($group);

            foreach ($group->getCommands() as $command) {
                foreach (getAllChecks($db, $client->getID(), 1) as $check) {
                    $client->addCheck($check);
                }
            }
        }
        return $client;
    }
?>