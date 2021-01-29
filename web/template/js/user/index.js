// add
document.querySelector('#add').addEventListener('click', function () {
    window.location = "/users/new"
});

// edit (GET)
document.querySelector('#edit').addEventListener('click', function () {
    var row = selected("用户", "请选择要更改的用户");
    if (row != null) {
        var id = row.cells[0].innerText;
        window.location = "/users/" + id;
    }
});

// delete
document.querySelector('#delete').addEventListener('click', function () {
    var row = selected("用户", "请选择要删除的用户");
    if (row != null) {
        if (delete_confirm()) {
            var id = row.cells[0].innerText;
            var url = "/users/" + id;
            ajax("DELETE", url, "", "/users");
        }
    }
});

// add roles for user
document.querySelector('#addRoles').addEventListener('click', function () {
    var row = selected("用户", "请选择要分配角色的用户");
    if (row != null) {
        var userID = row.cells[0].innerText;
        window.location = "/auth/roles4user/" + userID;
    }
});