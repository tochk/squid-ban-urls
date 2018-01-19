function CheckMACAddress(a) {
    var b = /^(?:([a-z]+):(?:([a-z]*):)?\/\/)?(?:([^:@]*)(?::([^:@]*))?@)?((?:[a-z0-9_-]+\.)+[a-z]{2,}|localhost|(?:(?:[01]?\d\d?|2[0-4]\d|25[0-5])\.){3}(?:(?:[01]?\d\d?|2[0-4]\d|25[0-5])))(?::(\d+))?(?:([^:\?\#]+))?(?:\?([^\#]+))?(?:\#([^\s]+))?$/i;
    return !!b.test(a)
}

var macAddrs = $("input.mac_valid");
jQuery(document).ready(function () {
    $("input.mac_valid").each(function () {
    });
    $("input.mac_valid").on("input", '[data-action="text"]', function () {
        var a = $(this), b = a.val();
    });
    $("#cancel").click(function () {
        $(".item:last-child").remove();
        if (!$('input.mac')[0]) $("#sbm").attr("disabled","disabled");
        $("input.mac").each(function () {
            0 == CheckMACAddress($(this).val()) ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success"), allOn = false) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"), $("#sbm").removeAttr("disabled"))
        });
        return $(".item:last-child").is("div") || $("#cancel").css("display", "none"), i--, !1
    });
});
if (i == null) {
    var i = 1;
}
$("#addNew").click(function () {
    addNew("");
    $("#sbm").attr("disabled","disabled");
});

function addNew(url) {
    $("#cancel").css("display", "inline-block"), macAddrs = $("input.mac_valid"), $(".input_forms").append(`
            <div class='item'>
                <h2 class='title'>Сайт #` + i + `<h2>
                <div class='form-group'>
                <label for='inputUrl` + i + `' class='col-md-2 control-label'>URL</label>
                <div class='col-md-10'>
                    <input type='mac' name='url` + i + `' class='form-control mac_valid mac' id='inputUrl` + i + `' placeholder='URL' required>
                    <span class='help-block'>Введите URL</span>
                    <a class="btn btn-primary link">Удалить ссылку</a>
                </div></div></div>`),
        $("input[name =url" + i + "]").val(url);
    i++;
    var b = "#ttl" + i;
    var numb = null;
    $(".link").click(function () {
        $(this).parent().parent().parent().parent().remove();
        let w = 1;
        $(".item").each(function () {
            $(this).children(".title").text("Сайт#" + w);
            $(this).children("h2:last-child").children(".form-group").children(".col-md-10").children("input").attr("name", "url" + w);
            w++;
        });
        if (!$('input.mac')[0]) {
            $("#sbm").attr("disabled","disabled");
            $("#cancel").css("display", "none");
        }
        $("input.mac").each(function () {
            0 == CheckMACAddress($(this).val()) ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success"), allOn = false) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"), $("#sbm").removeAttr("disabled"))
        });
        i = w;
        w = 1;
    });
    $("#inputUrl" + (i - 1)).bind('focusin' , function () {
        $(this).parent().parent().removeClass("has-error1")
    });
    $("#inputUrl" + (i - 1)).bind('focusout', function () {
        allOn = true;
        $("input.mac").each(function () {
            0 == CheckMACAddress($(this).val()) ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success"), allOn = false) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"))
        });
        if (allOn === true) {
            $("#sbm").removeAttr("disabled");
        }
    });
    $("#inputUrl" + (i - 1)).bind('keyup', function () {
        $("input.mac").each(function () {
            0 == CheckMACAddress($(this).val()) ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success"), $("#sbm").attr("disabled","disabled")) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"), $("#sbm").removeAttr("disabled"))
        });
    });
}

$("input.mac").each(function () {
    0 == CheckMACAddress($(this).val()) ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success")) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"))
});
$("input.mac").each(function () {
        if (0 == CheckMACAddress($(this).val())) return $("#sbm").attr("disabled", "disabled"), c = !1, !1
});
jQuery(document).ready(function () {
    var links = new Array();
    $("#addText").on('click', function (e) {
        var mess = $("#doc").val();
        var reg = /[-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?/gi;
        pregMatch = mess.match(reg);
        mess = mess.replace(reg, function (s) {
            var str = (/:\/\//.exec(s) === null ? "http://" + s : s);
            links.push(str);
            return null;
        });
        links.forEach(function (item) {
            addNew(item);
        });
        links = [];
        $("input.mac").each(function () {
            0 == CheckMACAddress($(this).val()) ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success"), allOn = false) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"), $("#sbm").removeAttr("disabled"))
        });
    });
});