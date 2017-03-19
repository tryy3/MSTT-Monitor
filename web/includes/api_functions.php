<?php
    include_once("api_funcs/api_common.php");
    include_once("api_funcs/api_errors.php");
    include_once("api_funcs/api_client.php");
    include_once("api_funcs/api_command.php");
    include_once("api_funcs/api_group.php");
    include_once("api_funcs/api_server.php");

    class API {
        private $client;
        private $command;
        private $group;
        private $server;

        public function __construct($db) {
            $this->client = new Client($db);
            $this->command = new Command($db);
            $this->group = new Group($db);
            $this->server = new Server($db);
        }

        public function Client() {
            return $this->client;
        }

        public function Command() {
            return $this->command;
        }

        public function Group() {
            return $this->group;
        }

        public function Server() {
            return $this->server;
        }
    }
    
    /* Download API */
    function download() {

    }
?>