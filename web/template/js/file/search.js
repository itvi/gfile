var urlSearch = new URLSearchParams(window.location.search);
var term = urlSearch.get("q");
// document.getElementsByName("q")[0].value;
document.querySelector('input[name="q"]').value = term;

