<script src="js/jquery-3.2.0.min.js"></script>
<script src="js/jquery-ui.min.js"></script>
<script src="js/bootstrap.min.js"></script>
<script src="js/bootstrap-switch.min.js"></script>
<script src="js/bootstrap-select.min.js"></script>
<script src="js/bootstrap-timepicker.js"></script>
<script src="js/jquery.bootstrap-touchspin.min.js"></script>
<script src="js/jquery.canvasjs.min.js"></script>
<script src="js/mstt-monitor.js"></script>

<script type="text/javascript">
$(function () {
	<?php
        $options = array();
        if ($page == "client") {
            $options["page"] = "client";
            foreach($checks as $k => $c) {
                echo "createGraph($('#graphCheck[data-check=\"".$k."\"]'), ".
                    json_encode($c->Output()).",".
                    json_encode($options).");\n";
            }

            echo 'createDropdown(".dropdown-monitor-main", ".dropdown-monitor-child",'.json_encode($config["AlertOptions"]).');';
            echo '$(".touchspin").TouchSpin({ step: 1, boostat: 5, maxboostedstep: 10 });';
        }
        if ($page == "start") {
            $options["page"] = "start";
            foreach($graphs as $k => $g) {
                echo "createGraph($('#graphCheck[data-check=\"".$k."\"]'), ".
                    json_encode($g->Output()).",".
                    json_encode($options). ");\n";
            }
        }
    ?>
});
</script>