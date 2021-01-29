function edit(id) {
    var sn = document.getElementById('sn').value;
    var name = document.getElementById('name').value;
    var email = document.getElementById("email").value;
    var password = document.getElementById("password").value;

    var formData = new FormData();
    formData.append("sn", sn);
    formData.append("name", name);
    formData.append("email", email);
    formData.append("password", password);

    var url = "/users/" + id;
    var redirect = "/users";
    ajax("PUT", url, formData, redirect);
}