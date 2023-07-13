<script type="text/javascript">
  WebUI.loadCssCdn("{{.cdn_url}}","querybuilder/custom-query-builder.min.css");
</script>
<div id="qbuilder_{{.id}}" class="qbuilder qbuilder-wrapper {{.styles}}"></div>
<script type="text/javascript">
  $(document).ready(function(){
    WebUI.loadScriptCdn("{{.cdn_url}}","js/extendext.min.js",function(){
    WebUI.loadScriptCdn("{{.cdn_url}}","js/doT.min.js",function(){
    WebUI.loadScriptCdn("{{.cdn_url}}","js/moment.min.js",function(){
    WebUI.loadScriptCdn("{{.cdn_url}}","querybuilder/custom-query-builder.min.js",function(){
      $('#qbuilder_{{.id}}').queryBuilder({
        plugins:['invert','not-group'],
        filters:[{{.filters}}],
        rules:{{.rules}},
        allow_groups:{{.allow_groups}},
        allow_empty:{{.allow_empty}},
        default_condition:'{{.default_condition}}',
        inputs_separator:'{{.inputs_separator}}',
        lang:{
          add_rule:"Add rule",
          add_group:"Add group",
          delete_rule:"Delete",
          delete_group:"Delete",
          invert:"Invert",
          NOT:"NOT",
          conditions:{
            AND:"AND",
            OR:"OR"
          },
          operators: {
            equal:"equal",
            not_equal:"not equal",
            in:"in",
            not_in:"not in",
            less:"less",
            less_or_equal:"less or equal",
            greater:"greater",
            greater_or_equal:"greater or equal",
            between:"between",
            not_between:"not between",
            begins_with:"begins with",
            not_begins_with:"does not begin with",
            contains:"contains",
            not_contains:"does not contain",
            ends_with:"ends with",
            not_ends_with:"does not end with",
            is_empty:"is empty",
            is_not_empty:"is not empty",
            is_null:"is null",
            is_not_null:"is not null"
          },
          errors: {
            no_filter:"No filter selected",
            empty_group:"The group is empty",
            radio_empty:"No value selected",
            checkbox_empty:"No value selected",
            select_empty:"No value selected",
            string_empty:"Empty value",
            string_exceed_min_length:"Must contain at least {0} characters",
            string_exceed_max_length:"Must not contain more than {0} characters",
            string_invalid_format:"Invalid format ({0})",
            number_nan:"Not a number",
            number_not_integer:"Not an integer",
            number_not_double:"Not a real number",
            number_exceed_min:"Must be greater than {0}",
            number_exceed_max:"Must be lower than {0}",
            number_wrong_step:"Must be a multiple of {0}",
            number_between_invalid:"Invalid values, {0} is greater than {1}",
            datetime_empty:"Empty value",
            datetime_invalid:"Invalid date format ({0})",
            datetime_exceed_min:"Must be after {0}",
            datetime_exceed_max:"Must be before {0}",
            datetime_between_invalid:"Invalid values, {0} is greater than {1}",
            boolean_not_valid:"Not a boolean",
            operator_not_multiple:"Operator ({1}) cannot accept multiple values"
          }
        }
      });

    })})})});
  });
</script>
