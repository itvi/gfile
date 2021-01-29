    function edit() {
        // get querystring (url) http://localhost:9000/auth/rule?sub=admin&obj=/users&act=GET
        var urlSearch = new URLSearchParams(window.location.search);
        var oSub = urlSearch.get('sub');
        var oObj = urlSearch.get('obj');
        var oAct = urlSearch.get('act');

        var sub = document.getElementById("sub").value;
        var obj = document.getElementById("obj").value;
        var act = document.getElementById("act").value;

        var formData = new FormData();
        formData.append("sub", sub);
        formData.append("obj", obj);
        formData.append("act", act);
        // old value
        formData.append("oSub", oSub);
        formData.append("oObj", oObj);
        formData.append("oAct", oAct);

        var url = "/auth";
        var redirect = "/auth";
        ajax("PUT", url, formData, redirect);
    }