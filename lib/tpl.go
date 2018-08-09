package lib

import (
	"html/template"
)

const indexHTML = `<!DOCTYPE html>
<html>
<head>
<title>Gravitational SSH Certificate Parser</title>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js"></script>
<link rel="stylesheet" href="https://gravitational.com/gravitational/bundle.css" />
<link href="https://gravitational.com/favicon.ico" rel="icon" type="image/x-icon">
<link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
<script>
$(document).ready(function(){
    $("form").submit(function(){
        $("#error").text("");
        $.ajax({
           type: "POST",
           url: "/v1/parsecert",
           data: JSON.stringify({ cert: $("#cert").val() }),
           contentType: "application/json; charset=utf-8",
           dataType: "json",
           success: function(data){
                $("#info").text(data.info);
           },
           error: function(err) {
              $("#error").text(err.responseJSON.message);
           }
       });
       return false;
    });
});
</script>
<style>
label {
  padding-top: 20px;
  font-size: 22px;
}
</style>
</head>
<body class="grv-teleport">
{{.Header}}
<div class="container">
  <div>
    <form>
      <div class="form-group">
        <label for="cert">Paste SSH Certificate here</label>
        <textarea id="cert" class="form-control" rows="15"></textarea>
        <span id="error" class="error text-danger"></span>
      </div>
     <button type="submit" class="btn btn-primary">Parse SSH Certificate</button>
    </form>
  </div>
  <div>
    <label for="info">Certificate Info</label>
    <pre id="info" style="height:400px;"></pre>
  </div>
</div>
</body>
</html>`

type templateParams struct {
	Header template.HTML
}

const headerHTML = `<nav id="top-nav" class="top-nav clearfix">
  <div class="top-nav-body">
  </div>
</nav>

	<header id="hero" class="background-midnight">
  <div class="hero background-globe">
    <h1>SSH Certificate Parser</h1>
    <h2>Built using Teleport: Modern SSH for Teams Managing Distributed Server Clusters</h2>
  </div>
</header>`
