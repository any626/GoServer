$(function(){
    $('.reply').click(function () {
        var $this = $(this);
        $this.parent().parent().after($('#reply-box').html());
        $this.parent().parent().next('.post-reply').children('form').attr('action', function(i, value) {
            return value + $this.attr('data-id');
        });
    });
    $('.post-list').on("click", ".reply-cancel-button", function(){
        $(this).parent().parent().parent().replaceWith('');
    });
    $('.edit').click(function () {
        var comment = $(this).parent().prev('.comment');
        var oldContent = comment.html();
        comment.replaceWith($('#edit-box').html())
        var $this = $(this);
        $this.parent().prev().children(':first-child').children('.edit-content').html(oldContent);
        $this.parent().prev().children(':last-child').children('.edit-cancel-button').attr('data-oldContent',oldContent);
        $this.parent().parent().children('form').attr('action', function(i, value) {
            return value + $this.attr('data-id');
        });
    });
    $('.post-list').on("click", ".edit-cancel-button", function(){
        var oldContent = ($(this).attr('data-oldContent'));
        $(this).parent().parent().replaceWith('<div class="comment">'+oldContent+'</div>');
    });
});