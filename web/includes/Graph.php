<?php
    class Graph {
        private $Options = array(); // CanvasJS Graph Options
        private $InternalOptions = array();
        private $GraphData = array(); // array

        public function Parse($options) {
            if (isset($options["Function"])) {
                $this->InternalOptions["Function"] = $options["Function"];
            }
            if (isset($options["Length"])) {
                $this->InternalOptions["Length"] = $options["Length"];
            }
            if (isset($options["From"])) {
                $this->InternalOptions["From"] = $options["From"];
            }
            if (isset($options["To"])) {
                $this->InternalOptions["To"] = $options["To"];
            }
            if (isset($options["Check"])) {
               $this->InternalOptions["Check"] = $options["Check"];
            }
            if (isset($options["CheckType"])) {
               $this->InternalOptions["CheckType"] = $options["CheckType"];
            }
            if(isset($options["ChartOptions"])) {
                $this->Options = $options["ChartOptions"];
            }
            
            if (isset($options["DataOptions"])) {
                foreach ($options["DataOptions"] as $opt) {
                    $graph = new GraphDataOption();
                    $graph->Parse($opt);
                    array_push($this->GraphData, $graph);
                }
            }
        }

        public function FillDataPoints($db, $id = -1) {
            $stmt = "";
            $type = (isset($this->InternalOptions["CheckType"])?$this->InternalOptions["CheckType"]:"all");

            switch($type) {
                case "all":
                    $stmt = $db->prepare("SELECT timestamp, response, checked, client_id, error, finished FROM checks WHERE command_id=? AND timestamp >= FROM_UNIXTIME(?) AND timestamp <= FROM_UNIXTIME(?)");
                    break;
                case "client":
                    $stmt = $db->prepare("SELECT timestamp, response, checked, error, finished FROM checks WHERE client_id=? AND command_id=? AND timestamp >= FROM_UNIXTIME(?) AND timestamp <= FROM_UNIXTIME(?)");
                    break;
                default:
                    error_log("Invalid CheckType please check the documentation.");
                    return;
            }

            foreach ($this->GraphData as $value) {
                $value->FillDataPoints($this->InternalOptions, $id, $stmt);
            }
        }

        public function Output() {
            $data = array();
            foreach ($this->GraphData as $v) {
                array_push($data, $v->Output());
            }
            return array_merge($this->Options, array("data"=>$data));
        }
    }

    class GraphDataOption {
        private $Options; // CanvasJS GraphDataOption
        private $DataPointsOptions; // array

        public function Parse($options) {
            if (isset($options["DataPointsOptions"])) {
                $dataPoints = new GraphDataPointOption();
                $dataPoints->Parse($options["DataPointsOptions"]);
                $this->DataPointsOptions = $dataPoints;
                unset($options["DataPointsOptions"]);
            }
            $this->Options = $options;
        }

        public function FillDataPoints($GraphOpt, $id, $stmt) {
            $this->DataPointsOptions->FillDataPoints($GraphOpt, $id, $stmt);
        }

        public function Output() {
            return array_merge($this->Options, array("dataPoints"=> $this->DataPointsOptions->Output()));
        }
    }

    class GraphDataPointOption {
        private $InternalOptions = array();
        private $Options = array(); // CanvasJS DataPointOption
        private $DataPoints = array();

        public function Parse($options) {
            if (isset($options["Params"])) {
                $this->InternalOptions["Params"] = $options["Params"];
                unset($options["Params"]);
            }
            if (isset($options["YFormat"])) {
                $this->InternalOptions["YFormat"] = $options["YFormat"];
                unset($options["YFormat"]);
            }
            if (isset($options["XFormat"])) {
                $this->InternalOptions["XFormat"] = $options["XFormat"];
                unset($options["XFormat"]);
            }

            $this->Options = $options;
        }

        public function FillDataPoints($GraphOpt, $id, $stmt) {
            if (!isset($GraphOpt["Check"])) {
                error_log("Couldn't find a check, please set the Check option.");
                return;
            }
            $check = $GraphOpt["Check"];

            if (isset($GraphOpt["Function"])) {
                switch ($GraphOpt["Function"]) {
                    case 'average':
                        $days = (isset($GraphOpt["Length"])) ? $GraphOpt["Length"] : 7;

                        if ($days <= 0) {
                            error_log("Days is lower then or equal to 0, please increase the day length.");
                            return;
                        }
                        $dataPoints = array();
                        for ($i = 0; $i < $days; $i++) {
                            $from = strtotime($i." days ago 00:00:00");
                            $to = strtotime(($i-1)." days ago 00:00:00");

                            if ($i == 0) {
                                $to = strtotime("now");
                            }

                            $data = $this->getData($stmt, $id, $check, $from, $to);
                            $prev = array();
                            $d = 0;
                            foreach ($data as $v) {
                                if (isset($v["error"]) && $v["error"]) {
                                    continue;
                                }
                                $json = ParseParams($this->InternalOptions["Params"], json_decode($v["response"], true));

                                if (!isset($prev[$v["client_id"]])) {
                                    $prev[$v["client_id"]] = $json;
                                    continue;
                                }

                                if ($prev[$v["client_id"]] > $json) {
                                    $d += $json;
                                } else {
                                    $d += $json - $prev[$v["client_id"]];
                                }

                                $prev[$v["client_id"]] = $json;
                            }
                            
                            // TODO: Check this stuff more deeply

                            array_push($dataPoints, array("y"=>$d, "label"=> date("m/d", $to)." - ".date("m/d", $from)));
                        }
                        $this->DataPoints = $dataPoints;
                        break;
                    case 'network':
                        $from = (isset($GraphOpt["From"])) ? $GraphOpt["From"] : "1 day ago";
                        $to = (isset($GraphOpt["To"])) ? $GraphOpt["To"] : "now";
                        
                        $from = strtotime($from);
                        $to = strtotime($to);

                        $data = $this->getData($stmt, $id, $check, $from, $to);

                        $prev = 0;
                        $dataPoints = array();
                        foreach ($data as $d) {
                            $YValue = ParseParams($this->InternalOptions["Params"], json_decode($d["response"], true));
                            $v = $YValue;

                            if ($prev <= 0) {
                                $prev = $YValue;
                                continue;
                            }
                            if ($prev < $YValue) {
                                $v = $YValue - $prev;
                            }

                            array_push($dataPoints, array("x"=>strtotime($d["timestamp"]), "y"=>$v));
                            $prev = $YValue;
                        }
                        $this->DataPoints = $dataPoints;
                    default:
                        error_log("Found function, but unsupported function was given, '"+$GraphOpt["Function"]."'.");
                        break;
                }
            } else {
                $from = (isset($GraphOpt["From"])) ? $GraphOpt["From"] : "1 day ago";
                $to = (isset($GraphOpt["To"])) ? $GraphOpt["To"] : "now";
                
                $from = strtotime($from);
                $to = strtotime($to);

                $data = $this->getData($stmt, $id, $check, $from, $to);

                $dataPoints = array();
                foreach ($data as $d) {
                    $v = ParseParams($this->InternalOptions["Params"], json_decode($d["response"], true));
                    array_push($dataPoints, array("x" => strtotime($d["timestamp"]), "y" => $v));
                }
                $this->DataPoints = $dataPoints;
            }
        }

        public function Output() {
            return $this->DataPoints;
        }

        private function format($format, $value) {
            switch($format) {
                case "procent":
                    return "<span style='\"'color: {color};'\"'>{label}:</span> ".number_format($value,2)."%";
                case "GB":
                    return "<span style='\"'color: {color};'\"'>{label}:</span> ".number_format($value/1000/1000/1000,2)." GB";
                case "MB":
                    return "<span style='\"'color: {color};'\"'>{label}:</span> ".number_format($value/1000/1000,2)." MB";
                case "KB":
                    return "<span style='\"'color: {color};'\"'>{label}:</span> ".number_format($value/1000,2)." KB";
                case "B":
                    return "<span style='\"'color: {color};'\"'>{label}:</span> ".number_format($value,2)." B";
            }
        }

        private function getData($stmt, $id, $command_id, $from, $to) {
            if ($id == -1) {
                $stmt->execute(array($command_id, $from, $to));
            } else {
                $stmt->execute(array($id, $command_id, $from, $to));
            }
            
            if ($stmt->rowCount() <= 0) {
                return array();
            }
            return $stmt->fetchAll(PDO::FETCH_ASSOC);
        }
    }
?>