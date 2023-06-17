<form id="frmLogin_{{.id}}" class="loginform-wrapper">
  <div class="form-group mb-2">
    <label class="control-label font-weight-bold mb-1 float-left">Username</label>
    <input type="text" class="form-control" name="username" value="">
  </div>
  <div class="form-group">
    <label class="control-label font-weight-bold mb-1 float-left">Password</label>
    <div class="input-group">
      <input type="password" class="form-control" name="password" value="">
      <div class="input-group-append">
        <span class="input-group-text" data-pwdview="password"><i class="fa fa-fw fa-eye-slash"></i></span>
      </div>
    </div>
  </div>
  <div class="form-group">
    <button type="submit" class="btn {{.btn_style}} float-right font-weight-bold">
      <i class="fa fas fa-ta fa-sign-in fa-sign-in-alt"></i>Login
    </button>
  </div>
</form>
<script type="text/javascript">
  $(document).ready(function(){
    WebUI.loadScriptCdn("{{.cdn_url}}","cryptojs/crypto-js.min.js",function(){
      $("#frmLogin_{{.id}}").submit(function(e){
        e.preventDefault(); WebUI.notify.clear();
        var u=$("#frmLogin_{{.id}} input[name=username]"),p=$("#frmLogin_{{.id}} input[name=password]");
        WebUI.loader.load("POST","{{.submit_url}}",{username:u.val(),authkey:"{{.authkey}}",digest:(p.val())?CryptoJS.SHA256("{{.authkey}}"+CryptoJS.SHA256(p.val())).toString():''},null,null,function(){p.val('');p.focus()},200);
        return false;
      });
    });
    $("#frmLogin_{{.id}} span[data-pwdview=password]")
      .on("mousedown touchstart",function(){
        $("#frmLogin_{{.id}} input[name=password]").attr("type","text");
        $(this).children("i").removeClass("fa-eye-slash").addClass("fa-eye")})
      .on("mouseup mouseleave touchend",function(){
        $("#frmLogin_{{.id}} input[name=password]").attr("type","password");
        $(this).children("i").removeClass("fa-eye").addClass("fa-eye-slash")});
  });
</script>
