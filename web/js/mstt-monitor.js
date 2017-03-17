/* Binära prefix */
const KIBI = 1024;
const MEBI = KIBI * 1024;
const GIBI = MEBI * 1024;
const TEBI = GIBI * 1024;

var ipreg = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/

function alert(clazz, title, body) {
    $(".alert").removeClass("alert-danger alert-success");
    $(".alert").show();
    $(".alert").addClass("alert-" + clazz);
    $(".alert-header").text(title);
    $(".alert-text").text(body);
}

function dragg($elem) {
    $($elem).draggable({
        helper: 'clone'
    }).click(function(ev) {
        $(this).draggable({ disabled: false })
    }).dblclick(function(ev) {
        $(this).draggable({ disabled: true })
        $(ev.target).focus()
    })
}

function dropGroup($elem) {
    $elem.droppable({
        drop: function(ev, ui) {
            var command = $(ui.draggable.children()[0]).text();
            var group = $(this).data("check");
            var table = $(this).find("table");

            $.getJSON('/api.php', { "api": "add_command_group", "group": group, "command": command }, function(data) {
                if (!data.error) {
                    $tr = $("<tr>");
                    $tr.append($("<td>").text(data.message[0]))
                    $tr.append($("<td>").text(data.message[1]))
                    $tr.append($("<td>").text(data.message[2]))

                    $td = $("<td>")
                    $select = $("<select>")
                            .data("for", data.message[0])
                            .data("target", "stop_error")
                            .addClass("form-control")
                            .css({"width":"80%", "display": "inline"})
                    $optionTrue = $("<option>").text("True")
                    $optionFalse = $("<option>").text("False")

                    if (data.message[3])
                        $optionTrue.attr("selected", "selected")
                    else 
                        $optionFalse.attr("selected", "selected")

                    $td.append($select.append($optionTrue).append($optionFalse))
                    $td.append($("<i>").addClass("remove-command-group fa fa-close fa-close-red fa-lg pull-right").data("id", data.message[0]))
                    $tr.append($td);
                    $(table).append($tr);
                    dragg($tr);

                    alert("success", "Success! ", "Added the command " + command + " to the group " + group);
                } else {
                    alert("danger", "Error! ", data.message);
                }
            })
        }
    })
}

function createGraph($elem, chartOptions, dataPoints) {
    var opt = {
        "axisX": {
            "labelFormatter": function(e) {
                var d = new Date(e.value * 1000);
                return ("0" + (d.getUTCMonth()+1)).slice(-2) + "/" +
                ("0" + d.getUTCDate()).slice(-2) + " " +
                ("0" + d.getUTCHours()).slice(-2) + ":" +
                ("0" + d.getUTCMinutes()).slice(-2) + ":" +
                ("0" + d.getUTCSeconds()).slice(-2)
            }
        },
        "data": dataPoints
    }
    Object.assign(opt, chartOptions)
    $elem.CanvasJSChart(opt);
}

$(document).ready(function() {
    dragg($('.drag'));
    dropGroup($('.drop-group'));

    $('.button-convert-size').click(function() {
        var target = $(this).data("target")
        var convert = $(".convert-size[data-identifier='" + target + "']")
        var format = convert.data('format')
        var prefix = $(this).text()
        var convert_format = prefix.slice(-1)
        var value = parseInt(convert.data('value'))
        var format_output = $(".convert-size-format[data-identifier='"+target+"']")

        if (format == 'disc' || format == 'memory' || format == 'network') {
            f = 'bytes'
        } else {
            f = 'bytes'
        }

        if (convert_format === 'b' && f == 'bytes') {
            value = value * 8;
        } else if (convert_format === 'B' && f == 'bits') {
            value = value / 8;
        }
        
        format_output.text(prefix)

        if (prefix.length == 1) {
            convert.text(value);
            return;
        }

        switch (prefix[0]) {
            case "K":
                convert.text((value / KIBI).toFixed(2));
                break;
            case "M":
                convert.text((value / MEBI).toFixed(2));
                break;
            case "G":
                convert.text((value / GIBI).toFixed(2));
                break;
            case "T":
                convert.text((value / TEBI).toFixed(2));
                break;
        }
    })

    $('.drop-command').droppable({
        drop: function(ev, ui) {
            var drag = $(ui.draggable);
            var id = $(drag.children()[0]).text();
            var group = $(drag).closest(".panel").data("check")

            $.getJSON("/api.php", { "api": "remove_command_group", "id": id, "group": group }, function(data) {
                if (!data.error) {
                    drag.remove();
                    alert("success", "Success! ", data.message);
                } else {
                    alert("danger", "Error! ", data.message);
                }
            })
        }
    })

    $(".group-dd").droppable({
        drop: function(ev, ui) {
            var drag = $(ui.draggable);
            var group = drag.data("target")
            var id = drag.data("id")
            var $this = $(this)
            var type = $this.data("type")

            $.getJSON("/api.php", { "api": "edit_client_group", "type": type, "id": id, "group": group }, function(data) {
                if (!data.error) {
                    switch(type.toLowerCase()) {
                        case "remove":
                            drag.remove()
                            break;
                        case "add":
                            $d = drag.clone()
                            dragg($d)
                            $this.append($d)
                            break;
                    }
                    alert("success", "Success! ", data.message)
                } else {
                    alert("danger", "Error! ", data.message)
                }
            })
        }
    })

    $(".delete-client").click(function() {
        var id = $(this).data("id")
        $.getJSON('/api.php', { "api": "delete_client", "id": id }, function(data) {
            if (!data.error) {
                window.location = "?page=clients"
            } else {
                alert("danger", "Error! ", data.message)
            }
        })
    })

    $(".add-group").click(function() {
        var group = $(".group-name").val()
        $.getJSON('/api.php', { "api": "group_exists", "group": group }, function(data) {
            if (!data.error) {
                if (!data.message) {
                    $panel = $("<div>").addClass("panel panel-primary checks-item drop-group").attr("data-check", group)
                    $heading = $("<div>").addClass("panel-heading").append($("<h3>").addClass("panel-title").text(group))
                    $table = $("<table>").addClass("table table-groups table-hover table-bordered table-condensed")
                    $tbody = $("<tbody>").append($("<tr>")
                        .append($("<th>").width("5%").text("ID"))
                        .append($("<th>").width("45%").text("Namn"))
                        .append($("<th>").width("25%").text("Nästa Check"))
                        .append($("<th>").width("25%").text("Stop Error")))
                    $panel.append($heading).append($table.append($tbody))
                    $(".group-list").append($panel)

                    dropGroup($panel)

                    $button = $("<button>").attr("type", "button").addClass("list-group-item checks").attr("data-target", group)
                    $h4 = $("<h4>").addClass("list-group-item-heading").text(group)
                    $i = $("<i>").addClass("delete-group fa fa-close fa-close-red fa-lg pull-right");
                    $(".list-group").append($button.append($h4.append($i)))
                } else {
                    alert("danger", "Error! ", "Group already exists")
                }
            } else {
                alert("danger", "Error! ", data.message)
            }
        })
    })

    $('.add-command').click(function() {
        $.getJSON('/api.php', { "api": "create_command" }, function(data) {
            if (!data.error) {
                $tr = $("<tr>");
                $tr.append($("<td>").text(data.message));
                $tr.append($("<td>").attr("contenteditable", "").data("for", "cmd").data("id", data.message).data("target", "namn"));
                $tr.append($("<td>").attr("contenteditable", "").data("for", "cmd").data("id", data.message).data("target", "command"));
                $tr.append($("<td>").attr("contenteditable", "").data("for", "cmd").data("id", data.message).data("target", "description"));
                $tr.append($("<td>")
                    .append($("<select>").data("for", "cmd").data("id", data.message).data("target", "format").addClass("form-control").css({ "width": "80%", "display": "inline" })
                        .append($("<option>").text("Nothing"))
                        .append($("<option>").text("Memory"))
                        .append($("<option>").text("Disc"))
                        .append($("<option>").text("Procent"))
                        .append($("<option>").text("Date"))
                        .append($("<option>").text("Seconds"))
                        .append($("<option>").text("Network")))
                    .append($("<i>").addClass("delete-command fa fa-close fa-close-red fa-lg pull-right").data("id", data.message)));

                $(".table-commands").append($tr);
                dragg($tr);

                alert("success", "Success! ", "Succesfully created a new command with the id " + data.message);
            } else {
                alert("danger", "Error! ", data.message);
            }
        })
    })

    $(document).on("click", ".delete-group", function() {
        group = $(this).closest(".checks").data("target")
        $.getJSON('/api.php', { "api": "delete_group", "group": group }, function(data) {
            console.log(data)
            if (!data.error) {
                console.log($(".checks[data-target='"+group+"']"))
                $(".checks[data-target='"+group+"']").remove()
                console.log($(".checks-item[data-check='"+group+"']"))
                $(".checks-item[data-check='"+group+"']").remove()
                alert("success", "Success! ", data.message)
            } else { 
                alert("danger", "Error! ", data.message)
            }
        })
    })

    $(document).on("click", ".remove-command-group", function() {
        var parent = $(this).parent().parent()
        var id = $(this).data("id");
        var group = $(this).closest(".panel").data("check")
        console.log(group);

        $.getJSON('/api.php', { "api": "remove_command_group", "id": id, "group": group }, function(data) {
            if (!data.error) {
                alert("success", "Success! ", data.message);
                parent.remove();
            } else {
                alert("danger", "Error! ", data.message);
            }
        })
    })

    $(document).on("click", ".delete-command", function() {
        var id = $(this).data("id");
        var parent = $(this).parent().parent();
        var name = $(parent.children()[1]).text();

        $.getJSON('/api.php', { "api": "delete_command", "id": id }, function(data) {
            if (!data.error) {
                alert("success", "Success! ", data.message);
                parent.remove();
                $(".table-groups").each(function(k, v) {
                    $(v).find("tr").each(function(key, value) {
                        if ($($(value).children()[1]).text() == name) {
                            $(value).remove();
                        }
                    })
                })
            } else {
                alert("danger", "Error! ", data.message);
            }
        })
    })

    $(document).on("change", "select", function(ev) {
        var $this = $(this);
        var id = $this.data("id");
        var target = $this.data("target");
        var f = $this.data("for");
        var opt = $this.val();

        var api = "";

        if (f == "cmd") {
            api = "edit_command";
        } else if (f == "group") {
            api = "edit_group";
        } else {
            return;
        }

        $.getJSON('/api.php', { "api": api, "id": id, "key": target, "value": opt.toLowerCase() }, function(data) {
            if (!data.error) {
                alert("success", "Success! ", data.message);
            } else {
                alert("danger", "Error! ", data.message);
            }
        })
    })

    $(document).on("blur", "[contenteditable]", function(ev) {
        var $this = $(this);
        var f = $this.data("for");
        var id = $this.data("id");
        var target = $this.data("target");
        var previous = $this.data("previous");

        if ($this.text() == previous) {
            return;
        }

        if (f == "cmd") {
            api = "edit_command";
        } else if (f == "group") {
            api = "edit_group";
        } else if (f == "client") {
            api = "edit_client"
        } else if (f == "server") {
            api = "edit_server"
        } else {
            return;
        }

        $.getJSON("/api.php", { "api": api, "id": id, "key": target, "value": $this.text() }, function(data) {
            if (!data.error) {
                alert("success", "Success! ", data.message);
                $this.data("previous", $this.text());
            } else {
                alert("danger", "Error! ", data.message);
                $this.text(previous);
            }
        })
    })

    $('.checks').click(function() {
        $('.checks-item').each(function(k, v) { $(v).hide() });
        $('.checks-item[data-check="' + $(this).data("target") + '"]').show()
    })

    $(".search").keyup(function() {
        var searchTerm = $(".search").val();
        var listItem = $('.results tbody').children('tr');
        var searchSplit = searchTerm.replace(/ /g, "'):containsi('")

        $.extend($.expr[':'], {
            'containsi': function(elem, i, match, array) {
                return (elem.textContent || elem.innerText || '').toLowerCase().indexOf((match[3] || "").toLowerCase()) >= 0;
            }
        });

        $(".results tbody tr").not(":containsi('" + searchSplit + "')").each(function(e) {
            $(this).attr('visible', 'false');
        });

        $(".results tbody tr:containsi('" + searchSplit + "')").each(function(e) {
            $(this).attr('visible', 'true');
        });

        var jobCount = $('.results tbody tr[visible="true"]').length;
        $('.counter').text(jobCount + ' item');

        if (jobCount == '0') { $('.no-result').show(); } else { $('.no-result').hide(); }
    });

    $(".clickable-row").click(function() {
        if (!$(this).attr("contenteditable")) {
            window.location = $(this).data("href");
        }
    })

    $(".btn-create-client").click(function() {
        var ip = $(".create-client").val()
        if (!ip) {
            alert("danger", "Error! ", "Var god skriv in ett IP-nummer.")
            return
        }
        if (!ipreg.test(ip)) {
            alert("danger", "Error! ", "Var god skriv in ett korrekt IP-nummer.")
            return
        }

        $.getJSON("/api.php", { "api": "create_client", "ip": ip }, function(data) {
            if (!data.error) {
                setTimeout(function() {
                    window.location = "?page=client&id="+data.message
                }, 5);
            } else {
                alert("danger", "Error! ", data.message);
            }
        })
    })

    $(document).on("click", ".delete-server", function() {
        var $this = $(this);
        var id = $this.data("id");

        $.getJSON("/api.php", { "api": "del_server", "id": id }, function(data) {
            if(!data.error) {
                $this.parent().parent().remove()
                alert("success", "Success! ", data.message);
            } else {
                alert("danger", "Error! ", data.message);
            }
        })
    })

    $(".add-server").click(function() {
        var ip = $(".server-ip").val()
        if (!ip) {
            alert("danger", "Error! ", "Var god skriv in ett IP-nummer.")
            return
        }
        if (!ipreg.test(ip)) {
            alert("danger", "Error! ", "Var god skriv in ett korrekt IP-nummer.")
            return
        }

        $.getJSON("/api.php", { "api": "add_server", "ip": ip }, function(data) {
            if (!data.error) {
                $span = $("<span>").text(ip).attr("contenteditable", true).data("id", data.message).data("for", "server").data("target", "ip").data("previous", ip)
                $i = $("<i>").data("id", data.message).addClass("delete-server fa fa-close fa-close-red fa-lg pull-right")
                $tr = $("<tr>")
                    .append($("<td>").text(data.message))
                    .append($("<td>").attr("contenteditable", true).data("id", data.message).data("for", "server").data("target", "namn"))
                    .append($("<td>").append($span).append($i))
                $(".table-servers").append($tr)
                alert("success", "Success! ", "Successfully added a server.");
            } else {
                alert("danger", "Error! ", data.message);
            }
        })
    })
});