{{define "group_items"}}
  {{range $k := .}}
    {{if eq $k.type "text"}}
    <span class="input-group-text">{{$k.value}}</span>
    {{else if eq $k.type "icon" }}
    <span class="input-group-text"><i class="fa fas {{$k.value}}"></i></span>
    {{else if eq $k.type "select" }}
    <select class="custom-select" {{if $k.name}}name="{{$k.name}}"{{end}}>
      {{range $o := $k.options }}<option value="{{$o.value}}" {{if $o.selected}}selected{{end}}>{{$o.label}}</option>{{end}}
    </select>
    {{end}}
  {{end}}
{{end}}

<script type="text/javascript">
  WebUI.loadCssCdn("{{.cdn_url}}","datetimepicker/datetimepicker.min.css");
  WebUI.loadCss("/static/webui/css/webui_inputform.min.css");
</script>
<form id="form_{{.id}}" class="form-wrapper {{.styles}}" method="POST" action="{{.submit_url}}" novalidate>
  {{range $f := .fields }}
    {{if eq $f.type "hidden"}}
      <input type="hidden" {{if $f.name}}name="{{$f.name}}"{{end}} {{if $f.value}}value="{{$f.value}}"{{end}}>
    {{else if eq $f.type "title"}}
      <div class="form-title text-muted py-2">{{$f.label}}</div>
    {{else}}
      <div class="form-row form-group">
        <label class="col-11 col-sm-3 col-form-label text-nowrap {{if and $f.required (ne $f.type "static")}}required{{end}}">{{$f.label}}</label>
        {{if and $f.helpguide (ne $f.type "static")}}
          <div class="col-1 col-sm-1 order-sm-12 mt-2 text-left text-info">
            <a class="btn_helpguide" data-helpguide="{{$f.helpguide}}"><i class="fa fas fa-fw fa-lg fa-question-circle"></i></a>
          </div>
        {{end}}
        <div class="col-12 col-sm-8">
          {{if eq $f.type "custom"}}
            {{$f.value}}
          {{else if eq $f.type "static"}}
            <input type="text" class="form-control-plaintext" value="{{$f.value}}" readonly>
          {{else if or (eq $f.type "checkbox") (eq $f.type "radio")}}
            {{range $i ,$o := $f.options}}
              <div class="form-check {{if eq $i 0}}pt-1{{end}}">
                <input type="{{$f.type}}" class="form-check-input" name="{{if eq $f.type "checkbox"}}{{$o.name}}{{else}}{{$f.name}}{{end}}" value="{{if eq $f.type "checkbox"}}1{{else}}{{$o.value}}{{end}}" {{if $o.selected}}checked{{end}}{{if $o.disabled}}disabled{{end}}{{if $f.required}}required{{end}}>
                <label class="form-check-label px-1">{{$o.label}}</label>
              </div>
            {{end}}
          {{else if eq $f.type "password"}}
            <div class="input-group">
              {{if $f.prepend}}<div class="input-group-prepend">{{template "group_items" $f.prepend}}</div>{{end}}
              <input type="password" class="form-control" name="{{$f.name}}" value="{{$f.value}}" placeholder="{{$f.placeholder}}" {{if $f.strength}}data-plugin="passStrengthify"{{end}} {{if $f.confirm}}data-confirm="{{$f.name}}"{{end}} {{if $f.required}}required{{end}}>
              <div class="input-group-append">
                {{if $f.append}}{{template "group_items" $f.append}}{{end}}
                <span class="input-group-text" data-pwdview="{{$f.name}}"><i class="fa fas fa-eye-slash"></i></span>
              </div>
            </div>
            {{if $f.confirm}}
              <div class="input-group pt-1">
                <div class="input-group-prepend">
                  <span class="input-group-text">Confirm</span>
                </div>
                <input type="password" class="form-control" name="{{$f.name}}_confirm" data-confirm="{{$f.name}}" value="{{$f.value}}" placeholder="{{$f.placeholder}}" {{if $f.required}}required{{end}}>
                <div class="input-group-append">
                  <span class="input-group-text" data-pwdview="{{$f.name}}_confirm"><i class="fa fas fa-eye-slash"></i></span>
                </div>
              </div>
            {{end}}
          {{else if or (eq $f.type "date") (eq $f.type "time") (eq $f.type "datetime")}}
            <div class="input-group">
              {{if $f.prepend}}<div class="input-group-prepend">{{template "group_items" $f.prepend}}</div>{{end}}
              <input type="text" class="form-control" name="{{$f.name}}" value="{{$f.value}}" placeholder="{{$f.placeholder}}" data-plugin="datetimepicker" data-format="{{if $f.format}}{{$f.format}}{{else if eq $f.type "date"}}YYYY-MM-DD{{else if eq $f.type "time"}}HH:mm:00{{else}}YYYY-MM-DD HH:mm:00{{end}}" {{if $f.required}}required{{end}}>
              <div class="input-group-append">
                {{if $f.append}}{{template "group_items" $f.append}}{{end}}
                <span class="input-group-text"><i class="fa fas fa-fw {{if eq $f.type "time"}}fa-clock fa-clock-o{{else}}fa-calendar fa-calendar-alt{{end}}"></i></span>
              </div>
            </div>
          {{else if eq $f.type "file"}}
            <div class="input-group">
              {{if $f.prepend}}<div class="input-group-prepend">{{template "group_items" $f.prepend}}</div>{{end}}
              <div class="custom-file">
                <input type="file" class="custom-file-input" name="{{$f.name}}" {{if $f.format}}accept="{{$f.format}}"{{end}} placeholder="{{if $f.placeholder}}{{$f.placeholder}}{{else}}Select files{{end}}" data-plugin="bsCustomFileInput" data-maxsize="{{$f.maxsize}}" {{if $f.multiple}}multiple{{end}}{{if $f.required}}required{{end}}>
                <label class="custom-file-label text-truncate empty">{{if $f.placeholder}}{{$f.placeholder}}{{else}}Select files{{end}}</label>
              </div>
              <div class="input-group-append">
                <span class="input-group-text text-danger d-none" data-fileerror="{{$f.name}}"></span>
                <button class="input-group-text" data-fileselect="{{$f.name}}">Browse</button>
                {{if $f.append}}{{template "group_items" $f.append}}{{end}}
              </div>
            </div>
          {{else}}
            <div class="input-group">
              {{if $f.prepend}}<div class="input-group-prepend">{{template "group_items" $f.prepend}}</div>{{end}}
              {{if eq $f.type "text"}}
                <input type="text" class="form-control" name="{{$f.name}}" {{if $f.value}}value="{{$f.value}}"{{end}} {{if $f.placeholder}}placeholder="{{$f.placeholder}}"{{end}} {{if $f.required}}required{{end}}>
              {{else if eq $f.type "textarea"}}
                <textarea class="form-control scroll" {{if $f.name}}name="{{$f.name}}"{{end}} {{if $f.rows}}rows="{{$f.rows}}"{{end}} {{if $f.placeholder}}placeholder="{{$f.placeholder}}"{{end}} {{if $f.required}}required{{end}}>{{$f.value}}</textarea>
              {{else if eq $f.type "select"}}
                <select class="custom-select scroll" {{if $f.name}}name="{{$f.name}}"{{end}} {{if $f.multiple}}multiple {{if $f.rows}}size="{{$f.rows}}"{{end}}{{end}} {{if $f.required}}required{{end}}>
                  {{range $o := $f.options}}
                    <option value="{{if $o.value}}{{$o.value}}{{end}}" {{if or $o.selected}}selected{{end}}{{if $o.disabled}}disabled{{end}}>{{$o.label}}</option>
                  {{end}}
                </select>
              {{end}}
              {{if $f.append}}<div class="input-group-append">{{template "group_items" $f.append}}</div>{{end}}
            </div>
          {{end}}
          {{if and $f.help (ne $f.type "static")}}<i class="form-text text-muted">{{$f.help}}</i>{{end}}
        </div>
      </div>
    {{end}}
  {{end}}
</form>
<script type="text/javascript">
  $(document).ready(function(){
    if($("#form_{{.id}} input[data-plugin='passStrengthify']").length){
      WebUI.loadScriptCdn("{{.cdn_url}}","js/passstrength.min.js",function(){
        $("#form_{{.id}} input[data-plugin='passStrengthify']").each(function(){
          $(this).passStrengthify({minimum:1,security:1})});
      });
    };
    if($("#form_{{.id}} input[data-plugin='datetimepicker']").length){
      WebUI.loadScriptCdn("{{.cdn_url}}","js/moment.min.js",function(){
      WebUI.loadScriptCdn("{{.cdn_url}}","datetimepicker/datetimepicker.min.js",function(){
        $("#form_{{.id}} input[data-plugin='datetimepicker']").each(function(){
          $(this).datetimepicker({format:$(this).data('format'),useCurrent:true,showClear:true,showClose:true})});
      })});
    };
    if($("#form_{{.id}} input[data-plugin='bsCustomFileInput']").length){
      WebUI.loadScriptCdn("{{.cdn_url}}","js/bs-custom-file-input.min.js",function(){
        bsCustomFileInput.init("#form_{{.id}} input[data-plugin='bsCustomFileInput']");
        $("#form_{{.id}} input[data-plugin='bsCustomFileInput']").on("change",function(e){
          var m=parseInt($(this).data('maxsize')),fsizechk=true;
          if($(this).val()){
            $(this).next(".custom-file-label").removeClass('empty');
            for(const f of $(this)[0].files) if(m>0 && f.size>m) fsizechk=false;
          }
          else {
            $(this).next(".custom-file-label").addClass('empty').html($(this).attr("placeholder"));
          };
          $("#form_{{.id}} button[data-fileselect="+$(this).attr('name')+"]")
            .html($(this).val()?'<b>&times;</b>':'Browse');
          if(fsizechk) {
            $(this).get(0).setCustomValidity('');
            $("#form_{{.id}} span[data-fileerror="+$(this).attr('name')+"]")
              .html('').addClass('d-none')
          }else{
            $(this).get(0).setCustomValidity('Invalid');
            $("#form_{{.id}} span[data-fileerror="+$(this).attr('name')+"]")
              .html('Too Large File').removeClass('d-none');
          };
        });
        $("#form_{{.id}} button[data-fileselect]").on("click",function(e){
          e.preventDefault();
          var sel=$("#form_{{.id}} input[name="+$(this).data('fileselect')+"]");
          if(sel.val()) sel.val('').trigger('change'); else sel.trigger('click');
          return false;
        });
      });
    };
    $("#form_{{.id}} input[data-confirm]").on("click change keypress keyup keydown touchstart touchend",function(){
      var m=$("#form_{{.id}} input[name="+$(this).data('confirm')+"]"),
        c=$("#form_{{.id}} input[name="+$(this).data('confirm')+"_confirm]");
      c.get(0).setCustomValidity((m.val()==c.val())?'':'Invalid');
    });
    $("#form_{{.id}} span[data-pwdview]")
      .on("mousedown touchstart",function(){
        $("#form_{{.id}} input[name="+$(this).data("pwdview")+"]").attr("type","text");
        $(this).children("i").removeClass("fa-eye-slash").addClass("fa-eye")})
      .on("mouseup mouseleave touchend",function(){
        $("#form_{{.id}} input[name="+$(this).data("pwdview")+"]").attr("type","password");
        $(this).children("i").removeClass("fa-eye").addClass("fa-eye-slash")});
    $("#form_{{.id}} select:not([multiple])").on("change",function(){
      if($(this).find(":selected").val()) $(this).removeClass("empty"); else $(this).addClass("empty");
    }).trigger('change');
    $("#form_{{.id}} a.btn_helpguide").on("click",function(e){
      e.preventDefault();
      WebUI.pagelock.modal(
        '<h5 class="m-0 py-2 text-info"><i class="fa fas fa-fw fa-question-circle"></i>Quick Guide</h5>',
        '<p class="">'+$(this).data('helpguide')+'</p>',
        '<button class="btn btn-primary" onclick="WebUI.pagelock.hide()">Got it</button>');
      return false;
    });
    $("#form_{{.id}}").on('submit',function(e){
      e.preventDefault();
      WebUI.notify.clear();
      WebUI.scrolltop();
      $(this).removeClass('was-validated').find("input,select,textarea").removeClass('is-invalid');
      if(this.checkValidity()===false){
        $(this).addClass('was-validated');
        WebUI.notify.error('Please fill all required fields');
      }else{
        WebUI.loader.formsubmit($(this),function(result){
          if(result.validation){
            for(var i=0;i<result.validation.length;i++){
              $("#form_{{.id}} input[name='"+result.validation[i]+"']").addClass('is-invalid');
              $("#form_{{.id}} select[name='"+result.validation[i]+"']").addClass('is-invalid');
              $("#form_{{.id}} textarea[name='"+result.validation[i]+"']").addClass('is-invalid');
            };
          };
          WebUI.request.success(result);
        },null,null,500);
      };
      return false;
    });
    $("#form_{{.id}}").on('reset',function(e){
      WebUI.notify.clear();
      WebUI.scrolltop();
      $(this).removeClass('was-validated');
      $(this)[0].reset();
      $(this).find("input,select,textarea").removeClass('is-valid is-invalid').trigger('change').trigger('keyup');
    });
  });
</script>
{{.jscript}}