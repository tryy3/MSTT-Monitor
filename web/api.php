<?php
    include_once("includes/db_connect.php");
    include_once("includes/function.php");
    include_once("includes/api_functions.php");
    header('Content-Type: application/json');

    $out = array("error"=>true,"message"=>"Unknown error");
    $version = "1.0";
    if (isset($_POST["version"])) {
        $version = $_POST["version"];
    }
    if (isset($_GET["version"])) {
        $version = $_GET["version"];
    }

    switch ($version) {
        case "1.0":
            if (isset($_GET["api"])) {
                switch(strtolower($_GET["api"])) {
                    case "download":
                        break;
                    case "create_client":
                        break;
                    case "create_group":
                        break;
                    case "create_command":
                        $out = createCommand($monitorDB);
                        break;
                    case "delete_client":
                        break;
                    case "delete_group":
                        break;
                    case "delete_command":
                        if (!isset($_GET["id"])) {
                            $out = array("error"=> true, "message" => "ID parameter is not set.");
                            break;
                        }
                        $out = deleteCommand($monitorDB, $_GET["id"]);
                        break;
                    case "edit_command":
                        if (!isset($_GET["id"])) {
                            $out = array("error" => true, "message" => "ID parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["key"])) {
                            $out = array("error" => true, "message" => "Key parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["value"])) {
                            $out = array("error" => true, "message" => "Value parameter is not set.");
                            break;
                        }
                        $out = editCommand($monitorDB, $_GET["id"], $_GET["key"], $_GET["value"]);
                        break;
                    case "edit_group":
                        if (!isset($_GET["id"])) {
                            $out = array("error" => true, "message" => "ID parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["key"])) {
                            $out = array("error" => true, "message" => "Key parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["value"])) {
                            $out = array("error" => true, "message" => "Value parameter is not set.");
                            break;
                        }
                        $out = editGroup($monitorDB, $_GET["id"], $_GET["key"], $_GET["value"]);
                        break;
                    case "edit_client":
                        break;
                    case "remove_command_group":
                        if (!isset($_GET["id"])) {
                            $out = array("error" => true, "message" => "ID parameter is not set.");
                            break;
                        }
                        $out = removeCommandGroup($monitorDB, $_GET["id"]);
                        break;
                    case "add_command_group":
                        if (!isset($_GET["group"])) {
                            $out = array("error" => true, "message" => "Group parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["command"])) {
                            $out = array("error" => true, "message" => "Command parameter is not set.");
                            break;
                        }
                        $out = addCommandGroup($monitorDB, $_GET["group"], $_GET["command"]);
                        break;
                    default:
                        $out = array("error"=>true,"message"=>"Invalid api function");
                }
            } else {
                $out = array("error"=>true,"message"=>"Invalid api function");
            }
            break;
        default:
            $out = array("error"=>true,"message"=>"Invalid api version");
    }

    echo json_encode($out, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES | JSON_NUMERIC_CHECK);
?>