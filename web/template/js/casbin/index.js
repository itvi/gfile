// add
document.querySelector('#add').addEventListener('click', function () {
    window.location = "/auth/new";
});

// edit
document.querySelector('#edit').addEventListener('click', function () {
    var row = selected("权限", "请选择要更改的权限");
    if (row != null) {
        var sub = row.cells[0].innerText;
        var obj = row.cells[1].innerText;
        var act = row.cells[2].innerText;

        var qstring = "?sub=" + sub + "&obj=" + obj + "&act=" + act;
        var url = "/auth/rule" + qstring;
        window.location = url;
    }
});

// delete
document.querySelector('#delete').addEventListener('click', function () {
    var row = selected("权限", "请选择要删除的权限");
    if (row != null) {
        if (delete_confirm()) {
            var sub = row.cells[0].innerText;
            var obj = row.cells[1].innerText;
            var act = row.cells[2].innerText;

            var formData = new FormData();
            formData.append("sub", sub);
            formData.append("obj", obj);
            formData.append("act", act);

            var url = "/auth";
            ajax("DELETE", url, formData, url);
        }
    }
});