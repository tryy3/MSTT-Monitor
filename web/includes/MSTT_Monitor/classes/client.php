<?php
    namespace MSTT_MONITOR\Classes;
    include_once('../common.php');
    
    class Client extends \MSTT_MONITOR\Common {
        private $groups = array();
        private $checks = array();
        private $latest = -1;
        private $id = -1;
        private $ip = "";
        private $name = "";
    }
?>