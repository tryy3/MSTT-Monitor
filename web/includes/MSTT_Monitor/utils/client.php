<?php
    namespace MSTT_MONITOR\Utils;
    
    /**
     * Client
     * 
     * Client is the class representation of a Client/Agent,
     * it contains information like what groups the client belongs to,
     * it latest checks, alert settings, ip, name etc.
     */
    class Client {
        /**
         * Client ID
         * 
         * @var int
         */
        private $id = -1;

        /**
         * An array of groupst that the client belong to.
         * @see Group
         * 
         * @var array
         */
        private $groups = array();

        /**
         * An array of the clients latest checks.
         * @see Check
         * 
         * @var array
         */
        private $checks = array();

        /**
         * An array of the clients latest alerts.
         * @see Alert
         * 
         * @var array
         */
        private $alerts = array();
        
        /**
         * An array of alert options that belongs to the client
         * @see AlertOptions
         * 
         * @var array
         */
        private $alertOptions = array();

        /**
         * A timestamp of the latest client check.
         * 
         * @var int
         */
        private $latest = -1;

        /**
         * The ip of the client
         * 
         * @var string
         */
        private $ip = "";
        
        /**
         * The clients name
         * 
         * @var string
         */
        private $name = "";

        /**
         * The construct function turns a raw mysql fetch into a Client Class
         * 
         * @param array $client The output of a \PDOStatement::fetch(\PDO::FETCH_ASSOC)
         * @return \Client
         */
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

        /**
         * A simple wrapper for getting the private values of the client.
         * Used for making configuration easier, shouldn't be used if possible.
         * 
         * @param string $key The key of a property that belongs to the client.
         * @return mixed 
         */
        public function get($key) {
            switch (strtolower($key)) {
                case 'id':
                    return $this->id;
                case 'groups':
                    return $this->groups; // Returns an array of group classes
                case 'groupnames':
                case 'group_names':
                    return $this->getGroupNames(); // Returns an array of group names
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

        /**
         * Add a new check to the client
         * 
         * @param Check $check The check to add to the client
         * @return void
         */
        public function addCheck($check) {
            // Check if the new check is newer then the latest check
            // if so set the latest check to this check.
            if ($check->getTimestamp() != -1 && $check->getTimestamp() > $this->latest) {
                $this->latest = $check->getTimestamp();
            }
            array_push($this->checks, $check);
        }

        /**
         * Add a new alert to the client
         * 
         * @param Alert $alert The alert to add to the client
         * @return void
         */
        public function addAlert($alert) {
            array_push($this->alerts, $alert);
        }

        /**
         * Add a new alert setting to the client
         * 
         * @param Alert $alert The alert setting to add to the client
         * @return void
         */
        public function addAlertOption($alertOption) {
            array_push($this->alertOptions, $alertOption);
        }


        /**
         * Set a new alert array for the client
         * 
         * @param array $alerts An array of alert
         * @return void
         */
        public function setAlerts($alerts) {
            $this->alerts = $alerts;
        }

        /**
         * Set a new alert options array for the client
         * 
         * @param array $alerts An array of alert settings
         * @return void
         */
        public function setAlertOptions($options) {
            $this->alertOptions = $options;
        }

        /**
         * Add a new group to the client
         * 
         * @param Group $group The new group to add to the client
         * @return void
         */
        public function addGroup($group) {
            array_push($this->groups, $group);
        }

        /**
         * Get an array of the names of the groups that the client belongs to
         * 
         * @return array
         */
        public function getGroupNames() {
            $out = "";
            foreach ($this->groups as $group) {
                $out .= $group->getName().", ";
            }
            return trim($out, " ,");
        }

        /**
         * Check if the client belongs to a group with a specific ID,
         * if so it will return it
         * 
         * @param int $id
         * @return Group
         */
        public function getGroupByID($id) {
            foreach ($this->groups as $group) {
                if ($group->getID() == $id) {
                    return $group;
                }
            }
            return NULL;
        }

        /**
         * Get the group
         * 
         * @param int $id
         * @return Group
         */
        public function getGroup($id) {
            if (isset($this->groups[$id])) {
                return $this->groups[$id];
            }
            return NULL;
        }

        /**
         * Undocumented function
         * 
         * @return void
         */
        public function getGroups() {
            return $this->groups;
        }

        public function getAlerts() {
            return $this->alerts;
        }
        
        public function getAlertOptions() {
            return $this->alertOptions;
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

        public function initAlerts($db) {
            foreach ($this->alertOptions as $alertOption) {
                foreach (getAlerts($db, $this->id, $alertOption->getClientID()) as $alert) {
                    $this->addAlert($alert);
                }
            }
        }
    }

    function getAllClients($db, $checks = 10) {
        $clients = array();
        $groups = getAllGroups($db);
        $alertSettings = getAllAlertOptions($db);

        $stmt = $db->query('SELECT * FROM clients');
        if (!$stmt) return $clients;
        while($row = $stmt->fetch(\PDO::FETCH_ASSOC)) {
            $client = new Client($row);
            foreach (explode(",", $row["group_names"]) as $groupName) {
                if (!isset($groups[$groupName])) {
                    continue;
                }
                
                $client->addGroup($groups[$groupName]);

                foreach ($groups[$groupName]->getCommands() as $command) {
                    foreach (getChecks($db, $client->getID(), $command->getCommandID(), $checks) as $check) {
                        $client->addCheck($check);
                    }
                }
            }
            array_push($clients, $client);

            foreach ($alertSettings as $setting) {
                if ($setting->getClientID() == $client->getID()) {
                    $client->addAlert($setting);
                }
            }
            $client->initAlerts($db);
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
        $alertSettings = getAlertOptions($db, $id);

        foreach (explode(",", $row["group_names"]) as $groupName) {
            $group = getGroup($db, $groupName);
            $client->addGroup($group);

            foreach ($group->getCommands() as $command) {
                foreach (getAllChecks($db, $client->getID(), 1) as $check) {
                    $client->addCheck($check);
                }
            }
        }

        $client->setAlertOptions($alertSettings);
        $client->initAlerts($db);
        return $client;
    }
?>