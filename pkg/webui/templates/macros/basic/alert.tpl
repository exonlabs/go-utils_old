<div class="alert-wrapper {{if .dismissible}}alert-dismissible fade show{{end}} {{.styles}}">
  <div class="alert alert-{{.alert_type}} text-left">
    {{ if .alert_icon}}<i class="fa fas fa-ta {{.alert_icon}}"></i>{{end}}
    {{.message}}
    {{ if .dismissible}}<button type="button" class="close" data-dismiss="alert"><span>&times;</span></button>{{end}}
  </div>
</div>
