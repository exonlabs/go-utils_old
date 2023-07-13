{{define "b_board_head"}}
    <link rel="stylesheet" type="text/css" href="/static/webui/css/webui.min.css">
    <link rel="stylesheet" type="text/css" href="/static/webui/css/webui_menuboard.min.css">
    <link rel="icon" type="image/png" href="/static/images/favicon.png">
    <meta name="theme-color" content="#343a40">
{{end}}

{{define "b_board_menuhead"}}
  <div class="p-3">
    <img class="img-fluid" src="/static/images/logo.png">
  </div>
{{end}}

{{define "b_board_pagehead"}}
    <div id="pagehead-bar" class="d-flex">
        <div id="pagehead-title" class="flex-grow-1 text-truncate text-left">
        {{block "b_pagehead_title" .}}{{end}}
        </div>
        <div id="pagehead-widgets" class="text-right d-print-none">
        {{block "b_pagehead_widgets" .}}{{end}}
        </div>
    </div>
    <div class="px-3 text-right text-info" style="line-height:18px; margin-bottom:-20px">
        <small>[info bar]</small>
    </div>
{{end}}

{{define "b_pagehead_title"}}
  <span>Sample Portal</span>
{{end}}

{{define "b_pagehead_widgets"}}
  <div class="btn-group">
    <button type="button" class="btn btn-sm dropdown-toggle" data-toggle="dropdown">
      <i class="fa fas fa-user"></i>
    </button>
    <div class="dropdown-menu dropdown-menu-right">
      <a class="dropdown-item btn-sm px-3" href="#?toogleboard=1">Toogle boards</a>
      {{if .langs}}
        <div class="dropdown-divider"></div>
        {{range $k, $v := .langs}}
          <a class="dropdown-item btn-sm px-3 pagelink" href="/?lang={{$k}}">{{$v}}</a>
        {{end}}
      {{end}}
      <div class="dropdown-divider"></div>
      <a class="dropdown-item btn-sm px-3" href="/loginpage">
        <i class="fa fas fa-flip fa-sign-out fa-sign-out-alt"></i>Logout
      </a>
    </div>
  </div>
{{end}}

{{define "b_board_pagebody"}}
  <div id="pagebody-contents" class="pt-3"></div>
{{end}}

{{define "b_board_scripts" }}
  {{if and .doc_lang (ne .doc_lang "en")}}<script type="text/javascript" src="/static/js/locale/{{.doc_lang}}.min.js"></script>{{end}}
  <script type="text/javascript" src="/static/webui/js/i18n/webui.min.js"></script>
  <script type="text/javascript" src="/static/webui/js/i18n/webui_menuboard.min.js"></script>
{{end}}