<div class="p-3">
  <div class="card">
    <div class="card-body text-left">
      {{.message}}
      {{if .langs}}
        <div class="btn-group-sm mt-3">
          {{range $k, $v := .langs}}
            <a class="btn btn-light px-3 pagelink" href="/?lang={{$k}}">{{$v}}</a>
            &nbsp;
          {{end}}
        </div>
      {{end}}
    </div>
  </div>
  <div class="card mt-2">
    <div class="card-body text-left">
        <a class="btn btn-secondary px-3" href="#?toogleboard=1">Toogle boards</a>
    </div>
  </div>
</div>
