function CheckMACAddress(a) {
    var b = /[-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?/;
    return !!b.test(a)
}
var macAddrs = $("input.mac_valid");
jQuery(document).ready(function () {
    $(".phone").mask("+7 (999) 999-9999", {autoclear: !1}), $("input.mac_valid").each(function () {
    }), $("input.mac_valid").on("input", '[data-action="text"]', function () {
        var a = $(this), b = a.val();
        console.log(CheckMACAddress(b))
    }), $("#cancel").click(function () {
        return $(".item:last-child").remove(), $(".item:last-child").is("div") || $("#cancel").css("display", "none"), i--, !1
    })
});
var i = 2;
$("#addNew").click(function () {
$("#cancel").css("display", "inline-block"), $("#sbm").attr("disabled", "disabled"), macAddrs = $("input.mac_valid"), $(".input_forms").append("<div class='item'><h2 class='title'>Сайт #" + i + "<h2><div class='form-group'><label for='inputUrl' class='col-md-2 control-label'>URL</label><div class='col-md-10'><input type='mac' name='url"+i + "' class='form-control mac_valid mac' id='inputUrl'           placeholder='URL' required>  <span class='help-block'>Введите URL</span></div></div></div>"), console.log("added"), $("#ttl" + i).mask("+7 (999) 999-9999", {autoclear: !0}), $("input.phone, input.mac, input.name").focusin(function () {
        $(this).parent().parent().removeClass("has-error1")
    });
    var b = "#ttl" + i;
    return a(document.querySelector(b)), $("input.phone").focusout(function () {
        $("input.phone").each(function () {
            17 !== $("input.phone").val().length || parseInt($("input.phone").val().indexOf("_")) !== -1 ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success")) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"))
        })
    }), $("input.mac").focusout(function () {
        $("input.mac").each(function () {
            0 == CheckMACAddress($(this).val()) ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success")) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"))
        })
    }), $("input.phone, input.mac, input.name").keyup(function () {
        var a = !0, b = !0, c = !0;
        $("input.phone").each(function () {
            if (17 !== $(this).val().length || parseInt($(this).val().indexOf("_")) !== -1)return $("#sbm").attr("disabled", "disabled"), console.log($(this).val()), console.log("pE" + i), a = !1, !1
        }), $("input.mac").each(function () {
            if (0 == CheckMACAddress($(this).val()))return $("#sbm").attr("disabled", "disabled"), console.log("mE" + i), c = !1, !1
        }), $("input.name").each(function () {
            if ("" == $(this).val())return $("#sbm").attr("disabled", "disabled"), console.log("nE" + i), b = !1, !1
        }), console.log(b, c, a), b && c && a ? $("#sbm").removeAttr("disabled") : $("#sbm").attr("disabled", "disabled")
    }), i++, !1
}), $("input.phone, input.mac, input.name").focusin(function () {
    $(this).parent().parent().removeClass("has-error1")
}), $("input.phone").focusout(function () {
    17 !== $("input.phone").val().length || parseInt($("input.phone").val().indexOf("_")) !== -1 ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success")) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"))
}), $("input.mac").focusout(function () {
    0 == CheckMACAddress($("input.mac").val()) ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success")) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"))
}), $("input.name").focusout(function () {
    "" == $("input.name").val() ? ($(this).parent().parent().addClass("has-error1"), $(this).parent().parent().removeClass("has-success")) : ($(this).parent().parent().addClass("has-success"), $(this).parent().parent().removeClass("has-error1"))
}), $("input.phone, input.mac, input.name").keyup(function () {0 == CheckMACAddress($("input.mac").val()) ? $("#sbm").attr("disabled", "disabled") : $("#sbm").removeAttr("disabled")
}), jQuery(document).ready(function () {

});