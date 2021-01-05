$(document).ready(function() {
    $("#files").on('click', 'tr', function(e) {
        $("tr").removeClass('selected');
        $(this).addClass('selected');
    });

    // var root = window.location.href.replace(/\/$/, ""); / *// remove end slash
    
    $("#files").on('dblclick', 'tr', function(e) {
        var $this = $(this).find('.fname');
        var name = $this.text();
        var isdir = $this.attr('data-isdir');
        var path = $this.attr('data-path');
        if(path=="/"){
            path = path + name;
        }else{
            path = path + "/" + name;
        }
        // only directory
        if(isdir=="true"){
            $.ajax({
                url: '/',
                type: 'GET',
                data: {
                    isdir: isdir,
                    name: name,
                    method: "ajax",
                    path: path
                },
                success: function(ret) {
                    showFiles(ret);
                    var li = '<li class="breadcrumb-item">'+
                        '<a href="'+ path + '">'+name+'</a>'+
                        '</li>';
                    $('.breadcrumb').append(li);
                }
            }); // end ajax
        }
    });
    
    $(".breadcrumb").on("click", "a", function (e) {
        e.preventDefault();
        var name = $(this).text();
        var url = $(this).attr('href');
        $(this).parent().nextAll().remove();

        $.ajax({
            url: '/',
            type: 'GET',
            data:{
                name: name,
                method: "ajax",
                path: url
            },
            success: function(ret){
                showFiles(ret);
            },
            error: function (x, s, e) {
                console.log(x, s, e);
            }
        }); // end ajax
    });

    // show files from server.
    function showFiles(ret){
        if($(ret).find('tbody').text().trim().length==0){
            $('#files tbody').html("空");
        }else{
            $('#files').html(ret);
        }
    }

    // download
    $('#download').click(function(){
        
    })
});
