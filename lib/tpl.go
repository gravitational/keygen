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
  <div class="row">
    <form>
      <div class="form-group">
        <label for="cert">Paste SSH Certificate here</label>
        <textarea id="cert" class="form-control" rows="15"></textarea>
        <span id="error" class="error text-danger"></span>
      </div>
     <button type="submit" class="btn btn-primary">Parse SSH Certificate</button>
    </form>
  </div>
  <div class="row">
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
  <div class="top-nav-logo is-inverse">
    <a href="https://gravitational.com"><img alt="Gravitational" src="https://gravitational.com/gravitational/images/logos/logo-gravitational-inverse.svg"></a>
  </div>
  <button id="top-nav-mobile-trigger" class="top-nav-mobile-trigger"><i class="material-icons">menu</i></button>
  <ul id="top-nav-menu">    
    <li>
      <a class="top-nav-button has-dropdown is-hidden-mobile" href="#">Solutions</a>
      <h4 class="top-nav-divider is-visible-mobile">Solutions</h4>
      <div class="top-nav-dropdown-overlay is-hidden"></div>
      <ul class="top-nav-dropdown is-hidden">        
        <li class="top-nav-dropdown-arrow">
          <span class="top-nav-dropdown-arrow-content"/>
        </li>
        <li class="top-nav-dropdown-item">
          <a class="has-error" href="https://gravitational.com/solutions/private-saas/">
            <strong>Private SaaS</strong>
            <em>Deliver your SaaS offering to Enterprise Buyers on Private Clouds or On-Premise</em>
          </a>
        </li>
        <li class="top-nav-dropdown-item">
          <a class="has-error" href="https://gravitational.com/solutions/msp/">
            <strong>Managed Services Providers</strong>
            <em>Scale your Managed Services across your distributed customer base</em>
          </a>
        </li>
        <li class="top-nav-dropdown-item">
          <a class="has-error" href="https://gravitational.com/solutions/edge/">
            <strong>Edge Computing Systems</strong>
            <em>Unified access patterns and ops automation across edge computing systems</em>
          </a>
        </li>		
      </ul>
    </li>    
	  
    <li>
      <a class="top-nav-button has-dropdown is-hidden-mobile" href="#">Products</a>
      <h4 class="top-nav-divider is-visible-mobile">Products</h4>
      <div class="top-nav-dropdown-overlay is-hidden"></div>
      <ul class="top-nav-dropdown is-hidden">
        <li class="top-nav-dropdown-arrow">
          <span class="top-nav-dropdown-arrow-content"/>
        </li>
        <li class="top-nav-dropdown-item">
          <a href="https://gravitational.com/teleport/">
            <strong>Teleport</strong>
            <em>Multi-Region SSH for Distributed Server Clusters</em>
          </a>
        </li>
        <li class="top-nav-dropdown-item">
          <a href="https://gravitational.com/telekube/">
            <strong>Telekube</strong>
            <em>Multi-Region Automation for Distributed Kubernetes Deployments</em>
          </a>
        </li>
      </ul> 
    </li>      
    
    <li class="top-nav-item-demo">
      <a class="se" href="https://gravitational.com/demo/">Schedule a Demo</a>
    </li>
  </ul>
  </div>
</nav>

	<header id="hero" class="background-midnight">
  <div class="hero background-globe">
    <h1>SSH Certificate Parser</h1>
    <h2>Built using Teleport: Modern SSH for Teams Managing Distributed Server Clusters</h2>
  </div>
</header>`
