{{define "b_html_head"}}
  <link rel="stylesheet" type="text/css" href="/static/vendor/bootstrap/bootstrap.min.css">
  <link rel="stylesheet" type="text/css" href="/static/vendor/fontawesome/css/fa4.min.css">
  <link rel="stylesheet" type="text/css" href="/static/vendor/pnotify/pnotify.min.css">
  {{block "b_page_head" .}}
    <link rel="stylesheet" type="text/css" href="/static/webui/css/webui.min.css">
  {{end}}
{{end}}

{{define "b_html_body"}}
  {{block "b_page_body" .}}{{end}}
  <script type="text/javascript" src="/static/vendor/jquery/jquery.min.js"></script>
  <script type="text/javascript" src="/static/vendor/jquery/jquery.i18n.min.js"></script>
  <script type="text/javascript" src="/static/vendor/bootstrap/bootstrap.bundle.min.js"></script>
  <script type="text/javascript" src="/static/vendor/pnotify/pnotify.min.js"></script>
  {{block "b_page_scripts" .}}
    <script type="text/javascript" src="/static/webui/js/i18n/webui.min.js"></script>
  {{end}}
{{end}}
