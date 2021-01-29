// add
document.querySelector('#add').addEventListener('click', function () {
    window.location = "/roles/new";
});

// edit
document.querySelector('#edit').addEventListener('click', function () {
    var row = selected("角色", "请选择要更改的角色");
    if (row != null) {
        var id = row.cells[0].innerText;
        window.location = "/roles/" + id;
    }
});

// delete
document.querySelector('#delete').addEventListener('click', function () {
    var row = selected("角色", "请选择要删除的角色");
    if (row != null) {
        if (delete_confirm()) {
            var id = row.cells[0].innerText;
            var url = "/roles/" + id;
            ajax("DELETE", url, "", "/roles");
        }
    }
});