<!doctype html>
<html lang="{{if .doc_lang}}{{.doc_lang}}{{else}}en{{end}}" dir="{{if .doc_langdir}}{{.doc_langdir}}{{else}}ltr{{end}}" class="scroll">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <title>{{.doc_title}}</title>
  {{block "b_html_head" .}}{{end}}
</head>
<body dir="{{if .doc_langdir}}{{.doc_langdir}}{{else}}ltr{{end}}" class="scroll">
  {{block "b_html_body" .}}{{end}}
  <!--[if lt IE 11 ]>
    <script type="text/javascript">document.body.innerHTML='<h2>You are using an outdated browser.</h2><p>Please <a href="https://browsehappy.com/">upgrade your browser</a> to improve your experience and security.</p>'</script>
  <![endif]-->
</body>
</html>
