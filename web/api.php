<?php
    include_once("includes/db_connect.php");
    include_once("includes/function.php");
    include_once("includes/api_functions.php");
    header('Content-Type: application/json');

    $api = new API($monitorDB);

    $errors = new Error();
    $errors->setMessage("Unknown error");
    
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
                    /* Download API */
                    case "download":
                        break;

                    /* Client API */
                    case "create_client":
                        if (!isset($_GET["ip"])) {
                            $errors->setMessage("IP parameter is not set.");
                            break;
                        }
                        $errors = $API->Client()->create($_GET["ip"]);
                        break;
                    case "delete_client":
                        if (!isset($_GET["id"])) {
                            $errors->setMessage("ID parameter is not set.");
                            break;
                        }
                        $errors = $API->Client()->delete($_GET["id"]);
                        break;
                    case "edit_client":
                        if (!isset($_GET["id"])) {
                            $errors->setMessage("ID parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["key"])) {
                            $errors->setMessage("Key parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["value"])) {
                            $errors->setMessage("Value parameter is not set.");
                            break;
                        }
                        switch($_GET["key"]) {
                            case "ip":
                                $errors = $API->Client()->editIP($_GET["id"], $_GET["value"]);
                                break;
                            case "name":
                                $errors = $API->Client()->editName($_GET["id"], $_GET["value"]);
                                break;
                        }
                        break;
                    case "edit_client_group":
                        if (!isset($_GET["group"])) {
                            $errors->setMessage("Group parameter is not set.");
                            break;
                        }
                        if(!isset($_GET["type"])) {
                            $errors->setMessage("Type parameter is not set.");
                            break;
                        }
                        if(!isset($_GET["id"])) {
                            $errors->setMessage("ID parameter is not set.");
                            break;
                        }
                        switch($_GET["type"]) {
                            case "add":
                                $errors = $API->Client()->addGroup($_GET["id"], $_GET["group"]);
                                break;
                            case "del":
                                $errors = $API->Client()->delGroup($_GET["id"], $_GET["group"]);
                                break;
                        }
                        break;

                    /* Group API */
                    case "group_exists":
                        if(!isset($_GET["group"])) {
                            $errors->setMessage("Group parameter is not set.");
                            break;
                        }
                        $errors = $API->Group()->exists($_GET["group"]);
                        break;
                    case "delete_group":
                        if (!isset($_GET["group"])) {
                            $errors->setMessage("Group parameter is not set.");
                            break;
                        }
                        $errors = $API->Group()->delete($_GET["group"]);
                        break;
                    case "edit_group":
                        if (!isset($_GET["id"])) {
                            $errors->setMessage("ID parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["key"])) {
                            $errors->setMessage("Key parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["value"])) {
                            $errors->setMessage("Value parameter is not set.");
                            break;
                        }
                        switch($_GET["key"]) {
                            case "next_check":
                                $errors = $API->Group()->editNextCheck($_GET["id"], $_GET["value"]);
                                break;
                            case "stop_error":
                                $errors = $API->Group()->editStopError($_GET["id"], $_GET["value"]);
                                break;
                        }
                        break;
                    case "remove_command_group":
                        if (!isset($_GET["id"])) {
                            $errors->setMessage("ID parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["group"])) {
                            $errors->setMessage("Group parameter is not set.");
                            break;
                        }
                        $errors = $API->Group()->removeCommand($_GET["id"], $_GET["group"]);
                        break;
                    case "add_command_group":
                        if (!isset($_GET["group"])) {
                            $errors->setMessage("Group parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["command"])) {
                            $errors->setMessage("Command parameter is not set.");
                            break;
                        }
                        $errors = $API->Group()->addCommand($_GET["group"], $_GET["command"]);
                        break;

                    /* Command API */
                    case "create_command":
                        $errors = $API->Command()->create();
                        break;
                    case "delete_command":
                        if (!isset($_GET["id"])) {
                            $errors->setMessage("ID parameter is not set.");
                            break;
                        }
                        $errors = $API->Command()->delete($_GET["id"]);
                        break;
                    case "edit_command":
                        if (!isset($_GET["id"])) {
                            $errors->setMessage("ID parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["key"])) {
                            $errors->setMessage("Key parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["value"])) {
                            $errors->setMessage("Value parameter is not set.");
                            break;
                        }
                        switch($_GET["key"]) {
                            case "command":
                                $errors = $API->Command()->editCommand($_GET["id"], $_GET["value"]);
                                break;
                            case "name":
                                $errors = $API->Command()->editName($_GET["id"], $_GET["value"]);
                                break;
                            case "description":
                                $errors = $API->Command()->editDescription($_GET["id"], $_GET["value"]);
                                break;
                            case "format":
                                $errors = $API->Command()->editFormat($_GET["id"], $_GET["value"]);
                                break;
                        }
                        break;

                    /* Server API */
                    case "add_server":
                        if (!isset($_GET["ip"])) {
                            $errors->setMessage("IP parameter is not set.");
                            break;
                        }
                        $errors = $API->Server()->create($_GET["ip"]);
                        break;
                    case "del_server":
                        if (!isset($_GET["id"])) {
                            $errors->setMessage("ID parameter is not set.");
                            break;
                        }
                        $errors = $API->Server()->delete($_GET["id"]);
                        break;
                    case "edit_server":
                        if (!isset($_GET["id"])) {
                            $errors->setMessage("ID parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["key"])) {
                            $errors->setMessage("Key parameter is not set.");
                            break;
                        }
                        if (!isset($_GET["value"])) {
                            $errors->setMessage("Value parameter is not set.");
                            break;
                        }
                        switch($_GET["key"]) {
                            case "ip":
                                $errors = $API->Server()->editIP($_GET["id"], $_GET["value"]);
                                return $errors;
                            case "namn":
                                $errors = $API->Server()->editName($_GET["id"], $_GET["value"]);
                                break;
                        }
                        break;

                    default:
                        $errors->setMessage("Invalid api function");
                }
            } else {
                $errors->setMessage("Invalid api function");
            }
            break;
        default:
            $errors->setMessage(>"Invalid api version");
    }

    if ($errors->getBaseURL() != null && $errors->getForm != null) {
        $e = $API->Server()->sendRequest($errors->getBaseURL(), $errors->getForm());
        if ($e->getError()) {
            $errors = $e;
        }
    }

    echo json_encode($errors->out(), JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES | JSON_NUMERIC_CHECK);
?>