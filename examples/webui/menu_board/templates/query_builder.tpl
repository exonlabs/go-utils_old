<div class="container-fluid pt-3 pb-5">
  <div class="row">
    <div class="col-xs-12 col-sm-12">
      <div id="QBuilderForm">
        {{.contents}}
      </div>
      <div class="text-right py-3">
        <button id="btn_reset" class="btn btn-sm btn-secondary mx-1">Reset</button>
        <button id="btn_submit" class="btn btn-sm btn-primary mx-1">Submit</button>
      </div>
      <pre id="Results" class="card p-3" style="display:none"></pre>
    </div>
  </div>
</div>
<script type="text/javascript">
  $(document).ready(function() {
    $("#btn_submit").on("click", function(){
      var result = "";
      var res=$("#QBuilderForm .qbuilder").queryBuilder('getRules');
      result += res?("JSON Rules:<br>"+JSON.stringify(res,null,2)):"";
      var res=$("#QBuilderForm .qbuilder").queryBuilder('getSQL');
      result += res?("<hr>SQL Rules:<br>"+JSON.stringify(res,null,2)):"";
      $("#Results").html(result?result:"Empty Result").show();
    });
    $("#btn_reset").on("click", function(){
      $("#QBuilderForm .qbuilder").queryBuilder('reset');
      $("#Results").html('').hide();
    });
  });
</script>