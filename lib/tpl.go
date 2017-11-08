package lib

const indexHTML = `<!DOCTYPE html>
<html>
<head>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js"></script>
<link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
<script>
$(document).ready(function(){
    $("form").submit(function(){
        $.ajax({
           type: "POST",
           url: "/v1/keypairs",
           data: JSON.stringify({ comment: $("#comment").val() }),
           contentType: "application/json; charset=utf-8",
           dataType: "json",
           success: function(data){
                $("#pub").text(data.pub);
                $("#priv").text(data.priv);
           },
           failure: function(errMsg) {
             alert(errMsg);
           }
       });
       return false;
    });
});
</script>
</head>
<body>
<div class="container">
  <div class="row">
    <h1>SSH Key Generator</h1>
  </div>
  <div class="row">
    <form>
      <div class="form-group">
        <label for="comment">Your Email</label>
        <input id="comment" type="text"/>
      </div>
     <button type="submit" class="btn btn-primary">Generate Keys</button>
    </form>
  </div>
  <div class="row">
    <label for="pub">Public Key</label>
    <textarea id="pub" class="form-control" rows="5"></textarea>
  </div>
  <div class="row">
    <label for="priv">Private Key</label>
    <textarea id="priv" class="form-control" rows="20"></textarea>
  </div>
</div>
</body>
</html>`
