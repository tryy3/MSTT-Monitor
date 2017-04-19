<script src="js/jquery-3.2.0.min.js"></script>
<script src="js/jquery-ui.min.js"></script>
<script src="js/bootstrap.min.js"></script>
<script src="js/bootstrap-switch.min.js"></script>
<script src="js/bootstrap-select.min.js"></script>
<script src="js/bootstrap-timepicker.js"></script>
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