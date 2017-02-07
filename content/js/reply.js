$(function(){
    $('.reply').click(function () {
        alert("Clicked reply on post id:" + $(this).attr('data-id'));
    });
    $('.edit').click(function () {
        var comment = $(this).parent().prev('.comment');
        var oldContent = comment.html();
        comment.replaceWith($('#edit-box').html())
        $(this).parent().prev().children(':first-child').children('.edit-content').html(oldContent);
        //alert("Clicked edit on post id:" + $(this).attr('data-id'));
    });
});