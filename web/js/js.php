<script src="/js/jquery-3.2.0.min.js"></script>
<script src="/js/jquery-ui.min.js"></script>
<script src="/js/bootstrap.min.js"></script>
<script src="/js/jquery.canvasjs.min.js"></script>
<script src="/js/mstt-monitor.js"></script>

<script type="text/javascript">
$(function () {
	<?php
        if ($page == "client") {
            foreach($checks as $k => $c) {
                echo "createGraph($('#graphCheck[data-check=\"".$k."\"]'), ".
                    json_encode($c["ChartOptions"]).", ".
                    json_encode($c["dataPoints"]).");\n";
            }
        }
    ?>
});
</script>