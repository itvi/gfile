function edit(id) {
    var name = document.getElementById('name').value;
    var desc = document.getElementById("desc").value;

    var formData = new FormData();
    formData.append("name", name);
    formData.append("desc", desc);

    var url = "/roles/" + id;
    var redirect = "/roles";
    ajax("PUT", url, formData, redirect);
}