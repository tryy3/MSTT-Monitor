<?php
    namespace MSTT_Monitor\API;

    include_once(__DIR__."/../common.php");
    include_once(__DIR__."/error.php");
    include_once(__DIR__."/client.php");
    include_once(__DIR__."/command.php");
    include_once(__DIR__."/group.php");
    include_once(__DIR__."/server.php");

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