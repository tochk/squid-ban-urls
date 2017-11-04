function CheckMACAddress(a) {
    //var b =/(([a-z]+:\/\/(www\.)*)*[a-z0-9\-_]+\.[a-z]+)/igm
    //var b = /[-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?/;
    var b = /^(?:([a-z]+):(?:([a-z]*):)?\/\/)?(?:([^:@]*)(?::([^:@]*))?@)?((?:[a-z0-9_-]+\.)+[a-z]{2,}|localhost|(?:(?:[01]?\d\d?|2[0-4]\d|25[0-5])\.){3}(?:(?:[01]?\d\d?|2[0-4]\d|25[0-5])))(?::(\d+))?(?:([^:\?\#]+))?(?:\?([^\#]+))?(?:\#([^\s]+))?$/i;
    return !!b.test(a)
}
var macAddrs = $("input.mac_valid");
jQuery(document).ready(function () {
    $("input.mac_valid").each(function () {
    }), $("input.mac_valid").on("input", '[data-action="text"]', function () {
        var a = $(this), b = a.val();
        console.log(CheckMACAddress(b))
    }), $("#cancel").click(function () {
        return $(".item:last-child").remove(), $(".item:last-child").is("div") || $("#cancel").css("display", "none"), i--, !1
    })
});
if (i == null) {
    var i = 1;
}
$("#addNew").click(function () {
    addNew("");
});
function addNew(url) {
    $("#cancel").css("display", "inline-block"), $("#sbm").attr("disabled", "disabled"), macAddrs = $("input.mac_valid"), $(".input_forms").append(`
            <div class='item'>
                <h2 class='title'>Сайт #` + i + `<h2>
                <div class='form-group'>
                <label for='inputUrl` + i + `' class='col-md-2 control-label'>URL</label>
                <div class='col-md-10'>
                    <input type='mac' name='url` + i + `' class='form-control mac_valid mac' id='inputUrl` + i + `' placeholder='URL' required>
                    <span class='help-block'>Введите URL</span>
                    <a class="btn btn-primary link">Удалить ссылку</a>
                </div></div></div>`),
        $("input[name =url" + i + "]").val(url),
        i++, console.log("added"), $("#ttl" + i).mask("+7 (999) 999-9999", {autoclear: !0}), $("input.phone, input.mac, input.name").focusin(function () {
        $(this).parent().parent().removeClass("has-error1")
    });
    var b = "#ttl" + i;
    var numb = null;
    $(".link").click(function () {
        $(this).parent().parent().parent().parent().remove();
        console.log("deleted");
        let w = 1;
        $(".item").each(function () {
            $(this).children(".title").text("Сайт#" + w);
            $(this).children("h2:last-child").children(".form-group").children(".col-md-10").children("input").attr("name", "url" + w);
            w++;
        });
        if ($("input.mac").val() == undefined) {
            $("#sbm").attr("disabled", "disabled");
        }
        i = w;
        w = 1;
    });
}
$("input.mac").each(function () {
    0 == CheckMACAddress($(this).val()) ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success")) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"))
}),
    $("input.mac").each(function () {
        if (0 == CheckMACAddress($(this).val()))return $("#sbm").attr("disabled", "disabled"), console.log("mE" + i), c = !1, !1
    }),

    jQuery(document).ready(function () {


        //console.log(111);
        let links = new Array();
        $("#addText").on('click', function (e) {
            var mess = $("#doc").val();
            //console.log(mess);
            var reg = /[-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?/gi;
            pregMatch = mess.match(reg);
            mess = mess.replace(reg, function (s) {
                let str = (/:\/\//.exec(s) === null ? "http://" + s : s );
                links.push(str);
                return null;//"<a href=\""+ str + "\">" + str /*s*/ + "</a>"; 
            });
            links.forEach(function (item) {
                addNew(item);
            });
            let allOn = true;
            $("input.mac").each(function () {
                0 == CheckMACAddress($(this).val()) ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success"), allOn = false) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"))
            });
            if (allOn == true) {
                $("#sbm").removeAttr("disabled");
            }
            if ($("input.mac").val() == undefined) {
                $("#sbm").attr("disabled", "disabled");
            }
            $("input.mac").focusout(function () {
                allOn = true;
                $("input.mac").each(function () {
                    0 == CheckMACAddress($(this).val()) ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success"), allOn = false) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"))
                });
                if (allOn == true) {
                    $("#sbm").removeAttr("disabled");
                }
            });
            $("input.mac").keyup(function () {
                $("input.mac").each(function () {
                    0 == CheckMACAddress($(this).val()) ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success")) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"))
                });
            });
            //console.log(links);
        });
    });