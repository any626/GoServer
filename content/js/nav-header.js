// check cookie and show / hide login / logout stuff in the nav
$(function(){
    function getCookie(cname) {
        var name = cname + "=";
        var decodedCookie = decodeURIComponent(document.cookie);
        var ca = decodedCookie.split(';');
        for(var i = 0; i <ca.length; i++) {
            var c = ca[i];
            while (c.charAt(0) == ' ') {
                c = c.substring(1);
            }
            if (c.indexOf(name) == 0) {
                return c.substring(name.length, c.length);
            }
        }
        return "";
    }

    var stuff = getCookie("login-name");
    if (stuff == "") {
        $("#logout").hide();
        $("#account").hide();
        $(".login").show();
    } else {
        $(".login").hide();
    }
}); 