<?php
    class Graph {
        private Options; // CanvasJS Graph Options
        private Func; // string
        private Length; // int
        private Check; // ID
        private From; // strtotime
        private To; // strtotime

        private GraphData; // array

        public function Parse($options) {
            if (isset($options["Function"])) {
                $this->Func = $options["Function"];
            }
            if (isset($options["Length"])) {
                $this->Length = $options["Length"];
            }
            if (isset($options["From"])) {
                $this->From = $options["From"];
            }
            if (isset($options["To"])) {
                $this->To = $options["To"];
            }
            if (isset($options["Check"])) {
                $this->Check = $options["Check"];
            }
            if(isset($options["ChartOptions"])) {
                $this->Options = $options["ChartOptions"];
            }

            if (isset($options["DataOptions"])) {
                $this->GraphData = array();

                foreach ($options["DataOptions"] as $opt) {
                    
                }
            }
        }
    }

    class GraphDataOption {
        private Options; // CanvasJS GraphDataOption
        private DataPoints; // array

        public function Parse($options) {

        }
    }

    class GraphDataPointOption {
        private Params; // array
        private YFormat; // string
        private XFormat; // string
        private Options; // CanvasJS DataPointOption

        private DataPoints; // array

        public function Parse($options) {
            if (isset($options["Params"])) {
                $this->Params = $options["Params"];
                unset($options["Params"]);
            }
            if (isset($options["YFormat"])) {
                $this->YFormat = $options["YFormat"];
                unset($options["YFormat"]);
            }
            if (isset($options["XFormat"])) {
                $this->XFormat = $options["XFormat"];
                unset($options["XFormat"]);
            }

            $this->Options = $options
        }
    }
?>