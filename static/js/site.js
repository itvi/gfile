$(document).ready(function () {
    var multipleSelect = false;

    $("#files").on('click', 'tr', function () {
        if (!multipleSelect) {
            $("tr").removeClass('selected');
            $(this).addClass('selected');
        } else {
            $(this).toggleClass('selected');
        }
    });

    // var root = window.location.href.replace(/\/$/, ""); / *// remove end slash

    // open directory
    $("#files").on('dblclick', 'tr', function () {
        var $this = $(this).find('.fname');
        var name = $this.text();
        var isdir = $this.attr('data-isdir');
        var path = $this.attr('data-path');

        if (path == "/") {
            path = path + name;
        } else {
            path = path + "/" + name;
        }

        if (isdir == "true") {
            $.ajax({
                url: '/',
                type: 'GET',
                data: {
                    isdir: isdir,
                    name: name,
                    method: "ajax",
                    path: path
                },
                success: function (ret) {
                    showFiles(ret);
                    var li = '<li class="breadcrumb-item">' +
                        '<a href="' + path + '">' + name + '</a>' +
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
            data: {
                name: name,
                method: "ajax",
                path: url
            },
            success: function (ret) {
                showFiles(ret);
            },
            error: function (x, s, e) {
                console.log(x, s, e);
            }
        }); // end ajax
    });

    // show files from server.
    function showFiles(ret) {
        if ($(ret).find('tbody').text().trim().length == 0) {
            $('#files tbody').html("空");
        } else {
            $('#files').html(ret);
        }
    }

    // multiple selected
    $('#mulSel').click(function () {
        $(this).toggleClass('multiple-select');
        // change global variable
        multipleSelect = multipleSelect == false ? true : false;
    });

    // download
    $('#download').click(function () {
        // selected or not
        var selections = $('#files .selected').length;
        if (selections == 0) {
            notify("please select");
            console.log("please select one or more rows");
            return;
        }

        // file or directory
        for (var i = 0; i < selections; i++) {
            var ele = $('#files .selected')[i];
            var file = $(ele).find('span');
            var isdir = file.attr('data-isdir');
            var name = file.text();
            var path = file.attr('data-path');
            // console.log(isdir,name,path)

            // directory will be zipped
            if (isdir == "true") {
                if (path == "/") {
                    path = "";
                }
                path = path + "/" + name;
                zipDir(isdir, path, name);
            } else {
                // ajax can't download
                // only working the first location
                //window.location="/dl?name="+name+"&path="+path;

                window.open("/dl?name=" + name + "&path=" + path, "_blank")
            }
        }
    });

    // return file path for download
    function zipDir(isdir, path, name) {
        $.ajax({
            url: '/zip',
            type: 'GET',
            data: {
                isdir: isdir,
                path: path,
                name: name
            },
            success: function (ret) {
                // if ret like A+B.zip then "+" will disappear
                ret = encodeURIComponent(ret);
                window.open("/dl?name=" + ret + "&isdir=" + isdir, "_blank")
            },
            error: function (x, s, e) {
                console.log(x, s, e);
            }
        });
    }

}); // end ready

function notify(message) {
    $.notify({
        icon: 'fa fa-info-circle',
        message: message
    }, {
        type: "info",
        allow_dismiss: true,
        delay: 2000, // 2 seconds
        placement: {
            from: "top",
            align: "center"
        },
        animate: {
            enter: "animate__animated animate__fadeInDown",
            exit: "animate__animated animate__fadeOutUp"
        }
    });
}