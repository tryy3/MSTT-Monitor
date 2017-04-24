<?php
    include_once(__DIR__."/includes/db_connect.php");
    include_once(__DIR__."/includes/function.php");
    include_once(__DIR__."/includes/MSTT_Monitor/api/api.php");

    use MSTT_MONITOR\API as API;

    header('Content-Type: application/json');

    $API = new API\API($monitorDB);

    $errors = new API\ErrorAPI();
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
                            case "namn":
                                $errors = $API->Client()->editName($_GET["id"], $_GET["value"]);
                                break;
                            default:
                                $errors->setMessage("Invalid key");
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
                            case "remove":
                                $errors = $API->Client()->delGroup($_GET["id"], $_GET["group"]);
                                break;
                            default:
                                $errors->setMessage("Invalid type");
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
                            default:
                                $errors->setMessage("Invalid key");
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
                        $errors = $API->Group()->removeCommand(intval($_GET["id"]), $_GET["group"]);
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
                        $errors = $API->Command()->delete(intval($_GET["id"]));
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
                                $errors = $API->Command()->editCommand(intval($_GET["id"]), $_GET["value"]);
                                break;
                            case "namn":
                                $errors = $API->Command()->editName(intval($_GET["id"]), $_GET["value"]);
                                break;
                            case "description":
                                $errors = $API->Command()->editDescription(intval($_GET["id"]), $_GET["value"]);
                                break;
                            case "format":
                                $errors = $API->Command()->editFormat(intval($_GET["id"]), $_GET["value"]);
                                break;
                            default:
                                $errors->setMessage("Invalid key");
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
                        $errors = $API->Server()->delete(intval($_GET["id"]));
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
                                $errors = $API->Server()->editIP(intval($_GET["id"]), $_GET["value"]);
                                break;
                            case "namn":
                                $errors = $API->Server()->editName(intval($_GET["id"]), $_GET["value"]);
                                break;
                            default:
                                $errors->setMessage("Invalid key");
                                break;
                        }
                        break;
                    case "manual_check":
                        if (!isset($_GET["command"])) {
                            $errors->setMessage("Command parameter is not set.");
                            break;
                        }
                        if (isset($_GET["save"])) {
                            $_GET["save"] = toBool($_GET["save"]);
                        }
                        if (isset($_GET["id"])) {
                            $_GET["id"] = intval($_GET["id"]);
                        }
                        if (isset($_GET["command_id"])) {
                            $_GET["command_id"] = intval($_GET["command_id"]);
                        }
                        $errors = $API->Server()->sendRequest("/check", $_GET, true);
                        break;

                    /* Alert API */
                    case "add_alert_option":
                        if (!isset($_GET["client_id"])) {
                            $errors->setMessage("Client ID parameter is not set.");
                            break;
                        }
                        $errors = $API->Alerts()->create(intval($_GET["client_id"]));
                        break;
                    case "delete_alert":
                        if (!isset($_GET["id"])) {
                            $errors->setMessage("ID parameter is not set.");
                            break;
                        }
                        $errors = $API->Alerts()->delete(intval($_GET["id"]));
                        break;
                    case "edit_alert":
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
                            case "alert":
                                $errors = $API->Alerts()->editAlert(intval($_GET["id"]), $_GET["value"]);
                                break;
                            case "value":
                                $errors = $API->Alerts()->editValue(intval($_GET["id"]), $_GET["value"]);
                                break;
                            case "count":
                                if (!is_numeric($_GET["value"])) {
                                    $errors->setMessage("Value needs to be numeric.");
                                    break;
                                }
                                $errors = $API->Alerts()->editCount(intval($_GET["id"]), intval($_GET["value"]));
                                break;
                            case "command":
                                if (!is_numeric($_GET["value"])) {
                                    $errors->setMessage("Value needs to be numeric.");
                                    break;
                                }
                                $errors = $API->Alerts()->editCommand(intval($_GET["id"]), intval($_GET["value"]));
                                break;
                            case "delay":
                                if (!is_numeric($_GET["value"])) {
                                    $errors->setMessage("Value needs to be numeric.");
                                    break;
                                }
                                $errors = $API->Alerts()->editDelay(intval($_GET["id"]), intval($_GET["value"]));
                                break;
                            case "service":
                                $errors = $API->Alerts()->editService(intval($_GET["id"]), $_GET["value"]);
                                break;
                            default:
                                $errors->setMessage("Invalid key");
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
            $errors->setMessage("Invalid api version");
    }

    if ($errors->getBaseURL() != null && $errors->getForm() != null) {
        $e = $API->Server()->sendRequest($errors->getBaseURL(), $errors->getForm());
        if ($e->getError()) {
            $errors = $e;
        }
    }

    echo json_encode($errors->out(), JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES | JSON_NUMERIC_CHECK);
?>