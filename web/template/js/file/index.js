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
        }

        if (isdir == "true") {
            $.ajax({
                url: '/getDirContent',
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
            url: '/getDirContent',
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
            $('#files tbody').html("<span style='margin-left:1rem'>(空)</span>");
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
            notify("下载", "请选择文件");
            console.log("please select one or more rows");
            return;
        }

        $('#spinner').css('display', 'inline');
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
                zipDir(isdir, path, name);
            } else {
                // ajax can't download
                // only working the first location
                //window.location="/dl?name="+name+"&path="+path;

                $('#spinner').css('display', 'none');
                window.open("/dl?name=" + name + "&path=" + path, "_blank");
            }
        }
    });

    // share links
    $('#share').click(function () {
        var selections = $('#files .selected').length;
        if (selections == 0 || selections > 1) {
            notify("分享", "请选择1个文件");
            console.log("please select one or more rows");
            return;
        }
        $(this).after($('#spinner'));
        $('#spinner').css('display', 'inline');
        // share links
        var dlink;
        var ele = $('#files .selected');
        var file = $(ele).find('span');
        var isdir = file.attr('data-isdir');
        var name = file.text();
        var path = file.attr('data-path');

        // directory will be zipped
        if (isdir == "true") {
            if (path == "/") {
                path = "";
            }
            $.ajax({
                url: '/zip',
                type: 'GET',
                data: {
                    isdir: isdir,
                    path: path,
                    name: name
                },
                success: function (ret) {
                    var r = encodeURIComponent(ret);
                    dlink = window.location.host + "/dl?name=" + r + "&isdir=" + isdir;
                    copy(dlink);
                    notify("已复制到剪贴板", dlink, 3600000); // 1 hour
                    $('#spinner').css('display', 'none');
                    $('#download').after($('#spinner'));
                },
                error: function (x, s, e) {
                    console.log(x, s, e);
                }
            });
            // copy(window.location.host + dlink);
        } else {
            dlink = window.location.host + "/dl?name=" + name + "&path=" + path;
            copy(dlink);
            notify("已复制到剪贴板", dlink, 3600000);
            $('#spinner').css('display', 'none');
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
                $('#spinner').css('display', 'none');
                window.open("/dl?name=" + ret + "&isdir=" + isdir, "_blank");
            },
            error: function (x, s, e) {
                console.log(x, s, e);
            }
        });
    }

    // search
    $('#search').click(function (e) {
        var term = $('input[name="q"]')[0].value;
        if (term.trim() == "") {
            e.preventDefault();
            notify("搜索", "请输入要搜索的文件名（支持模糊查询）");
            return;
        }
    }); // end search

}); // end ready
