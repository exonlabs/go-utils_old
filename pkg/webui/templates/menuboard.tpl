{{define "b_html_head" }}
  <link rel="stylesheet" type="text/css" href="/static/vendor/bootstrap/bootstrap.min.css">
  <link rel="stylesheet" type="text/css" href="/static/vendor/fontawesome/css/fa4.min.css">
  <link rel="stylesheet" type="text/css" href="/static/vendor/pnotify/pnotify.min.css">
  <link rel="stylesheet" type="text/css" href="/static/vendor/metismenu/metisMenu.min.css">
  {{block "b_board_head" .}}
      <link rel="stylesheet" type="text/css" href="/static/webui/css/webui.min.css">
      <link rel="stylesheet" type="text/css" href="/static/webui/css/webui_menuboard.min.css">
  {{end}}
{{end}}

{{define "b_html_body"}}
  {{block "b_board_body" .}}
    <div id="board-wrapper" class="ease h-100" style="display:none">
      <div id="board-menu" class="ease h-100 d-print-none">
        <div id="board-menutoggle" class="text-right">
          <a title="Toggle Menu"><i class="fa fas fa-bars"></i></a>
        </div>
        <div id="board-menubody" class="h-100 scroll">
          <div id="board-menuhead">
            {{block "b_board_menuhead" .}}{{end}}
          </div>
          {{block "board_menubody" .}}
            <ul class="metismenu">
              {{range $i, $m := .menu}}
                {{if $m.SubMenu}}
                  <li>
                    <a id="submenu_{{$i}}" href="#" class="has-arrow">{{if $m.Icon}}<i class="fa fas fa-fw fa-ta {{$m.Icon}}"></i>{{end}}{{$m.Label}}</a>
                    <ul>
                      {{range $j, $sm := $m.SubMenu}}
                        <li><a class="pagelink" href="{{$sm.URL}}">{{if $sm.Icon}}<i class="fa fas fa-fw fa-ta {{$sm.Icon}}"></i>{{end}}{{$sm.Label}}</a></li>
                      {{end}}
                    </ul>
                  </li>
                {{else if (ne $m.URL "#")}}
                  <li>
                    <a class="pagelink" href="{{$m.URL}}">{{if $m.Icon}}<i class="fa fas fa-fw fa-ta {{$m.Icon}}"></i>{{end}}{{$m.Label}}</a>
                  </li>
                {{end}}
              {{end}}
              <li><a></a></li>
            </ul>
          {{end}}
        </div>
      </div>
      <div id="board-page" class="h-100 scroll">
        <div id="board-pagehead" class="sticky-top">
          {{block "b_board_pagehead" .}}
            <div id="pagehead-bar" class="d-flex">
              <div id="pagehead-title" class="flex-grow-1 text-truncate text-left">
                {{block "b_pagehead_title" .}}{{end}}
              </div>
              <div id="pagehead-widgets" class="text-right d-print-none">
                {{block "b_pagehead_widgets" .}}{{end}}
              </div>
            </div>
          {{end}}
        </div>
        <div id="board-pagebody" class="h-100">
          {{block "b_board_pagebody" .}}
            <div id="pagebody-contents"></div>
          {{end}}
        </div>
      </div>
    </div>
    <div id="board-backdrop" class="d-print-none"></div>
  {{end}}
  <script type="text/javascript" src="/static/vendor/jquery/jquery.min.js"></script>
  <script type="text/javascript" src="/static/vendor/jquery/jquery.i18n.min.js"></script>
  <script type="text/javascript" src="/static/vendor/bootstrap/bootstrap.bundle.min.js"></script>
  <script type="text/javascript" src="/static/vendor/pnotify/pnotify.min.js"></script>
  <script type="text/javascript" src="/static/vendor/metismenu/metisMenu.min.js"></script>
  {{block "b_board_scripts" .}}
    <script type="text/javascript" src="/static/webui/js/i18n/webui.min.js"></script>
    <script type="text/javascript" src="/static/webui/js/i18n/webui_menuboard.min.js"></script>
  {{end}}
{{end}}

