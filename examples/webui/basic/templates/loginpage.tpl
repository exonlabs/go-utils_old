{{define "b_page_head" }}
  <link rel="stylesheet" type="text/css" href="/static/webui/css/webui.min.css">
  <link rel="icon" type="image/png" href="/static/images/favicon.png">
  <meta name="theme-color" content="#343a40">
{{end}}

{{define "b_page_body"}}
  <div id="loginpage-wrapper">
    <div class="container-fluid ">
      <div class="row">
        <div class="col-xs-12 p-4">
          <img src="/static/images/logo.png" height="50px">
        </div>
      </div>
      <div class="row">
        <div class="col-xs-12 col-sm-6 col-lg-4 mx-auto pt-5">
          <div class="card bg-light">
            <div id="formbody" class="card-body pb-3" style="font-size:13px"></div>
          </div>
          {{if .langs}}
            <div class="btn-group-sm float-left">
              {{range $k, $v :=  .langs}}
                <a class="btn btn-link px-2" href="?lang={{$k}}">{{$v}}</a>
              {{end}}
            </div>
          {{end}}
        </div>
      </div>
    </div>
  </div>
{{end}}

{{define "b_page_scripts"}}
  {{if and .doc_lang (ne .doc_lang "en")}}<script type="text/javascript" src="/static/js/locale/{{.doc_lang}}.min.js"></script>{{end}}
  <script type="text/javascript" src="/static/webui/js/i18n/webui.min.js"></script>
  <script type="text/javascript">
    $(document).ready(function(){
      WebUI.loader.load("GET","{{.load_url}}",null,function(r){
        WebUI.doctitle.update(r.doctitle);
        $("#formbody").html(r.payload);
      });
    });
  </script>
{{end}}
