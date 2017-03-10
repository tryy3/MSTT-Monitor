/* Binära prefix */
const KIBI = 1024;
const MEBI = KIBI * 1024;
const GIBI = MEBI * 1024;
const TEBI = GIBI * 1024;

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

$(document).ready(function() {
    $('.button-convert-size').click(function() {
        var target = $(this).data("target")
        var convert = $(".convert-size[data-identifier='" + target + "']")
        var format = convert.data('format')
        var prefix = $(this).text()
        var convert_format = prefix.slice(-1)
        var value = parseInt(convert.data('value'))

        if (convert_format === 'b' && format == 'bytes') {
            value = value * 8;
        } else if (convert_format === 'B' && format == 'bits') {
            value = value / 8;
        }

        if (prefix.length == 1) {
            convert.text(value);
            return;
        }

        switch (prefix[0]) {
            case "K":
                convert.text(value / KIBI);
                break;
            case "M":
                convert.text(value / MEBI);
                break;
            case "G":
                convert.text(value / GIBI);
                break;
            case "T":
                convert.text(value / TEBI);
                break;
        }
    })

    dragg($('.drag'));

    $('.drop-group').droppable({
        drop: function(ev, ui) {
            var command = $(ui.draggable.children()[0]).text();
            var group = $(this).data("check");
            var table = $(this).find("table");

            $.getJSON('/api.php', { "api": "add_command_group", "group": group, "command": command }, function(data) {
                if (!data.error) {
                    $tr = $("<tr>");
                    $.each(data.message, function(k, v) {
                        $tr.append($("<td>").text(v));
                    })
                    $(table).append($tr);
                    dragg($tr);

                    $($tr.children()[3]).append($("<i>").addClass("remove-command-group fa fa-close fa-close-red fa-lg").css("padding-left", "60%").data("id", data.message[0]))

                    alert("success", "Success! ", "Added the command " + command + " to the group " + group);
                } else {
                    alert("danger", "Error! ", data.message);
                }
            })
        }
    })

    $('.drop-command').droppable({
        drop: function(ev, ui) {
            var drag = $(ui.draggable);
            var id = $(drag.children()[0]).text();

            $.getJSON("/api.php", { "api": "remove_command_group", "id": id }, function(data) {
                if (!data.error) {
                    drag.remove();
                    alert("success", "Success! ", data.message);
                } else {
                    alert("danger", "Error! ", data.message);
                }
            })
        }
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
                        .append($("<option>").text("Bytes"))
                        .append($("<option>").text("Bits")))
                    .append($("<i>").addClass("delete-command fa fa-close fa-close-red fa-lg").css("padding-left", "60%").data("id", data.message)));

                $(".table-commands").append($tr);
                dragg($tr);

                alert("success", "Success! ", "Succesfully created a new command with the id " + data.message);
            } else {
                alert("danger", "Error! ", data.message);
            }
        })
    })

    $(document).on("click", ".remove-command-group", function() {
        var parent = $(this).parent().parent()
        var id = $(this).data("id");

        $.getJSON('/api.php', { "api": "remove_command_group", "id": id }, function(data) {
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

    $(document).on("blur", "td", function(ev) {
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
        window.location = $(this).data("href");
    })
});