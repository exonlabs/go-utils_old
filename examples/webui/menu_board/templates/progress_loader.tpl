 <div class="p-3">
  <button id="btnLoad" class="btn btn-primary">Load with Progress update</button>
  <script type="text/javascript">
    $(document).ready(function(){
      $("#btnLoad").bind('click', function(e){
        WebUI.loader.progress("POST","/loader",{},null,null,null,200,1000);
      });
    });
  </script>
</div>