$(function(){
    function update_button(){
        var isValid = namevalid && passvalid;
        $('#submit').prop('disabled',!isValid);
    }
    function nametest(){
        return !(/\/|\?|#/.test($('#name').val()));
    }
    function passtest(){
        return $('#pass').val().length > 9;
    }
    var namevalid = nametest();
    var passvalid = passtest();
    $('#name').keyup(function(event){
        namevalid = nametest();
        update_button();
    });
    $('#pass').keyup(function(event){
        passvalid = passtest();
        update_button();
    });
    update_button();
});