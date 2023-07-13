<script type="text/javascript">
  WebUI.loadCssCdn("{{.cdn_url}}","datatables/datatables.min.css");
  WebUI.loadCss("/static/webui/css/webui_datagrid.min.css");
</script>
<table id="tblDataGrid_{{.id}}" class="table table-sm table-hover text-nowrap w-100 {{.styles}}"></table>
<script type="text/javascript">
  $(document).ready(function(){
    WebUI.loadScriptCdn("{{.cdn_url}}","datatables/datatables.min.js",function(){
      function dataTables_expfname(){var tz=(new Date()).getTimezoneOffset()*60000,ts=(new Date(Date.now()-tz)).toISOString().slice(0,19).replace(/T/g,'').replace(/-/g,'').replace(/:/g,'');return '{{.export.file_prefix}}_'+ts};
      $.fn.dataTable.ext.errMode='none';
      $.fn.DataTable.ext.pager.numbers_length=5;
      $("#tblDataGrid_{{.id}}").on('error.dt',function(e,s,t,msg){WebUI.notify.error(msg,null,true)});
      var dt=$("#tblDataGrid_{{.id}}").DataTable({
        dom:'<"float-right"<"float-left"f><"float-left"l>><"#DataGrid_Ops_{{.id}}.float-left"><"table-responsive scroll"t><"float-left"i><"float-right pt-2"p>',
        colReorder:true,
        pagingType:"full_numbers",
        stateSave:true,
        search:{smart:false,regex:true,caseInsensitive:true},
        lengthMenu:[["{{join .lenMenu "\",\""}}"],["{{replaceAll (join .lenMenu "\",\"") "-1" "All"}}"]],
        select:{style:'multi+shift',selector:'td.row-select>input[type=checkbox]'},
        columns:[
          {defaultContent:'',searchable:false,orderable:false,className:'dtctrl table-index px-0'},
          {defaultContent:'',searchable:false,orderable:false,className:'dtctrl row-select px-3'},
          {{.columns}}
        ],
        order:[{{.order}}],
        buttons:{
          dom:{button:{className:'dropdown-item pl-3 pr-4 py-0 text-left tick-select'}},
          buttons:[
            {extend:'columnsToggle',columns:':not(.dtctrl)'},
            {extend:'csvHtml5',className:'d-none',
             filename:dataTables_expfname,exportOptions:{columns:':visible:not(.dtctrl)',orthogonal:'_'},
             fieldSeparator:'{{.export.csv_fieldSeparator}}',fieldBoundary:'{{.export.csv_fieldBoundary}}',
             escapeChar:'{{.export.csv_escapeChar}}',extension:'{{.export.csv_extension}}'},
            {extend:'excelHtml5',className:'d-none',title:'{{.export.file_title}}',sheetName:'{{.export.xls_sheetName}}',
             filename:dataTables_expfname,exportOptions:{columns:':visible:not(.dtctrl)',orthogonal:'_'},
             extension:'{{.export.xls_extension}}'},
            {extend:'print',className:'d-none',title:'{{.export.file_title}}',
             exportOptions:{columns:':visible:not(.dtctrl)',orthogonal:'_'}},
          ]
        },
        language:{
          search:'<i title="Search" class="fa fas fa-search"></i>',
          lengthMenu:'_MENU_',
          info:"Showing _START_-_END_ of _TOTAL_",
          infoEmpty:"",
          infoFiltered:"",
          emptyTable:"No data available",
          zeroRecords:"No matching records found",
          paginate:{
            first:'<i class="fa fas fa-angle-double-left"></i>',
            last:'<i class="fa fas fa-angle-double-right"></i>',
            previous:'<i class="fa fas fa-angle-left"></i>',
            next:'<i class="fa fas fa-angle-right"></i>'},
          select:{rows:'',columns:'',cells:''},
        }
      });
      $("#tblDataGrid_{{.id}}_length").prepend('<div class="dataTables_custombtns mb-2"><button id="btnReload_{{.id}}" class="btn btn-sm border px-3" title="Reload"><i class="fa fas fa-fw fa-refresh fa-sync-alt"></i></button>'+
        {{if .export.types}}
          '<div class="dropdown"><button class="btn btn-sm border dropdown-toggle" data-toggle="dropdown" title="Export"><i class="fa fas fa-fw fa-download fa-file-export"></i></button><div class="dropdown-menu dropdown-menu-right p-0 pb-1" style="min-width:100px"><h6 class="dropdown-header px-3">Export</h6>'+
          {{if inStr .export.types "csv"}}'<button id="btnExpCSV_{{.id}}" class="dropdown-item pl-3 pr-4 py-1"><i class="fa fas fa-fw fa-file-text-o fa-file-csv"></i> csv</button>'+{{end}}
          {{if inStr .export.types "xls"}}'<button id="btnExpXLS_{{.id}}" class="dropdown-item pl-3 pr-4 py-1"><i class="fa fas fa-fw fa-file-excel-o fa-file-excel"></i> excel</button>'+{{end}}
          {{if inStr .export.types "print"}}'<button id="btnPRINT_{{.id}}" class="dropdown-item pl-3 pr-4 py-1"><i class="fa fas fa-fw fa-print"></i> print</button>'+{{end}}
          '</div></div>'+
        {{end}}
        '<div class="dropdown"><button class="btn btn-sm border dropdown-toggle" data-toggle="dropdown" title="Columns"><i class="fa fas fa-fw fa-th-list"></i></button><div id="btnDataGridVis_{{.id}}" class="dropdown-menu dropdown-menu-right p-0" style="min-width:100px"><h6 class="dropdown-header px-3">Show / Hide</h6></div></div></div>');
      dt.buttons().container().appendTo('#btnDataGridVis_{{.id}}');
      $(dt.column(1).header()).html('<input type="checkbox">');
      $('input[type=checkbox]',dt.column(1).header()).on("click",function(){
        if($(this).is(':checked')) dt.rows({search:'applied'}).select(); else dt.rows().deselect();
      });
      var serchindxs=[]; dt.settings()[0].aoColumns.forEach(function(col){if(col.bSearchable)serchindxs.push(col.idx)});
      function draw_rows(){
        dt.draw().rows().deselect();
        dt.column(1).nodes().each(function(cell,i){cell.innerHTML='<input type="checkbox">'});
      };
      function search_update(){
        var visindxs=dt.columns(':visible:not(.dtctrl)').indexes().toArray();
        dt.settings()[0].aoColumns.forEach(function(col){
          col.bSearchable=(serchindxs.indexOf(col.idx)!=-1 && visindxs.indexOf(col.idx)!=-1)});
        dt.rows().invalidate(); draw_rows();
      };
      dt.on('draw',function(){
        var e=$(".dataTables_wrapper .dataTables_paginate, .dataTables_wrapper .dataTables_info");
        if(dt.data().length) e.removeClass("d-none"); else e.addClass("d-none");
      });
      dt.on('order.dt search.dt',function(){
        dt.column(0,{search:'applied',order:'applied'}).nodes().each(function(cell,i){cell.innerHTML=i+1});
        dt.rows({search:'removed'}).deselect();
      });
      dt.on('select deselect',function(){
        var w=$("#DataGrid_Ops_{{.id}}"),c=dt.rows({selected:true,search:'applied'}).count();
        $('input[type=checkbox]',dt.column(1).header()).prop('checked',c>0);
        dt.column(1).nodes().flatten().to$().each(function(){
          $(this).find('input[type=checkbox]').prop("checked",$(this).parents('tr').hasClass('selected'));
        });
        w.html('');
        if(c>0){
          if(c==1){ {{range $o := .single_ops}}w.append('<button class="dropdown-item" {{if $o.action}}data-op="{{$o.action}}"{{end}} {{if $o.confirm}}data-confirm="{{$o.confirm}}"{{end}}>{{if $o.label}}{{$o.label}}{{end}}</button>');{{end}} }
          else{ {{range $o := .group_ops}}w.append('<button class="dropdown-item" {{if $o.action}}data-op="{{$o.action}}"{{end}} {{if $o.confirm}}data-confirm="{{$o.confirm}}"{{end}}>{{if $o.label}}{{$o.label}}{{end}}</button>');{{end}} };
          if(w.children().length){
            w.wrapInner('<div class="dropdown-menu '+($('body').attr('dir')=='rtl'?'dropdown-menu-right':'')+'"></div>');
            w.prepend('<button class="btn btn-sm border dropdown-toggle" data-toggle="dropdown"><span class="title">Select Operation</span></button>');
            w.wrapInner('<div class="dropdown dataTables_opbtns"></div>');
            $("#DataGrid_Ops_{{.id}} button[data-op]").off('click').on('click',function(){
              var op=$(this).data('op'),confirm=$(this).data('confirm');
              var run_op = function(){WebUI.loader.load(
                "POST","{{.baseurl}}/"+op,{items:dt.rows({selected:true,search:'applied'}).ids().toArray()},
                function(r){if(r.reload)$('#btnReload_{{.id}}').trigger('click');WebUI.request.success(r)},null,null,200)};
              if(confirm){
                WebUI.pagelock.modal(
                  '<h5 class="m-0 py-2 text-info">Confirm</h5>','<p>'+confirm+'</p>',
                  '<button class="btn btn-sm btn-secondary px-3" onclick="WebUI.pagelock.hide()">No</button><button id="btnConfirmOp_{{.id}}" class="btn btn-sm btn-primary px-3">Yes</button>');
                $('#btnConfirmOp_{{.id}}').on('click',function(){WebUI.pagelock.hide();run_op()});
              }else run_op();
            });
          };
        };
      });
      dt.on('column-visibility.dt',search_update);
      search_update();
      $('#btnReload_{{.id}}').on("click",function(){
        WebUI.loader.cancel();dt.rows().deselect().clear().draw();
        $(dt.table().body()).html('<tr><td colspan="100" class="loading"></td></tr>');
        WebUI.loader.lock_timer=setTimeout(function(){
          WebUI.loader.req_xhr=WebUI.request("POST","{{.loadurl}}",{},
            function(r){
              if(r.payload){dt.clear().rows.add(r.payload);draw_rows()}
              else{dt.clear().draw().rows().deselect()};
              WebUI.request.success(r)},
            function(e){
              $(dt.table().body()).html('<tr><td colspan="100" class="text-center">Failed loading data</td></tr>');
              WebUI.request.error(e)},
            function(){WebUI.loader.reset()});
        },100);
      }).trigger('click');
      $('#btnExpCSV_{{.id}}').on("click",function(){dt.button('.buttons-csv').trigger()});
      $('#btnExpXLS_{{.id}}').on("click",function(){dt.button('.buttons-excel').trigger()});
      $('#btnPRINT_{{.id}}').on("click",function(){dt.button('.buttons-print').trigger()});
      $('#btnDataGridVis_{{.id}}').on("click",function(e){e.stopPropagation()});
    });
  });
</script>
{{.jscript}}
