$(function(){
    // start a reply
    $('.reply').click(function () {
        var $this = $(this);
        $this.parent().after($('#reply-box').html());
        $this.parent().next('.post-reply').children('form').children('div').children('textarea').focus();
        $this.parent().next('.post-reply').children('form').attr('action', function(i, value) {
            var type = $this.attr('data-type') == "post" ? "post" : $this.attr('data-parent');
            return value + type + "/" + $this.attr('data-id');
        });
    });
    // cancel replying
    $('.post-list').on("click", ".reply-cancel-button", function(){
        $(this).parent().parent().parent().replaceWith('');
    });
    // start editing
    $('.edit').click(function () {
        var comment = $(this).parent().prev('.comment');
        var oldContent = comment.html();
        comment.replaceWith($('#edit-box').html());
        var $this = $(this);
        $this.parent().parent().children('form').children('div').children('textarea').focus();
        $this.parent().prev().children(':first-child').children('.edit-content').html(oldContent.replace(/<img src="|">/g,""));
        $this.parent().prev().children(':last-child').children('.edit-cancel-button').attr('data-oldContent',oldContent);
        $this.parent().parent().children('form').attr('action', function(i, value) {
            return value + $this.attr('data-type') + "/" + $this.attr('data-id');
        });
    });
    // cancel editing
    $('.post-list').on("click", ".edit-cancel-button", function(){
        var oldContent = ($(this).attr('data-oldContent'));
        $(this).parent().parent().replaceWith('<div class="comment">'+oldContent+'</div>');
    });
    // hides/shows child comments
    $('.hider').click(function(){
        var comment = $(this).parent().parent();
        if (comment.hasClass("mini")){
            $(this).html("[ - ]");
            comment.removeClass("mini");
        } else {
            $(this).html("[ + ]");
            comment.addClass("mini");
        }
    });
    // inlines pictures
    $('.comment').each(function(){
        var content = $(this).html();
        var matches = content.match(/http.+\.(png|jpg|gif)/g);
        if (matches) {
            var unique = Array.from(new Set(matches));            
            unique.forEach(function(element){
                var newelem = "<img src=\"" + element + "\"/>";
                content = content.replace(new RegExp(element, 'g'), newelem);
            });
            $(this).html(content);
        }
    });
});