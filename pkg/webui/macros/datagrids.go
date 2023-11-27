package macros

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/exonlabs/go-utils/pkg/webui"
)

var (
	dataGridPath = "templates/macros/datagrids/"
	stdRender    = map[string]any{"_": "_", "display": "d"}
)

func UiStdDataGrid(options map[string]any, styles string) (string, error) {
	tplName := dataGridPath + "stddatagrid.tpl"

	var columns []string
	if _, ok := options["columns"]; ok {
		for _, k := range options["columns"].([]map[string]any) {
			column := map[string]any{
				"name":  k["id"],
				"title": k["title"],
			}
			if _, ok := k["data"]; !ok {
				column["data"] = k["id"]
			} else {
				column["data"] = k["data"]
			}
			if _, ok := k["render"]; !ok {
				column["render"] = stdRender
			} else {
				column["render"] = k["render"]
			}
			if _, ok := k["type"]; !ok {
				column["type"] = "string"
			} else {
				column["type"] = k["type"]
			}
			if _, ok := k["visible"]; !ok {
				column["visible"] = true
			} else {
				column["visible"] = k["visible"]
			}
			if _, ok := k["searchable"]; !ok {
				column["searchable"] = true
			} else {
				column["searchable"] = k["searchable"]
			}
			if _, ok := k["orderable"]; !ok {
				column["orderable"] = true
			} else {
				column["orderable"] = k["orderable"]
			}

			jsoncol, err := json.Marshal(column)
			if err != nil {
				return "", err
			}
			columns = append(columns, string(jsoncol))
		}
	}

	export := map[string]any{
		"types":              []string{"csv", "xls", "print"},
		"file_title":         "",
		"file_prefix":        "export",
		"csv_fieldSeparator": ",",
		"csv_fieldBoundary":  "",
		"csv_escapeChar":     "\"",
		"csv_extension":      ".csv",
		"xls_sheetName":      "Sheet1",
		"xls_extension":      ".xlsx",
	}
	if _, ok := options["export"]; ok {
		for k, v := range options["export"].(map[string]any) {
			export[k] = v
		}
	}
	// default data
	if _, ok := options["cdn_url"]; !ok {
		options["cdn_url"] = ""
	}
	if _, ok := options["grid_id"]; !ok {
		options["grid_id"] = RandInt(10000)
	}
	if _, ok := options["base_url"]; !ok {
		options["base_url"] = ""
	}
	if _, ok := options["load_url"]; !ok {
		options["load_url"] = ""
	}
	if _, ok := options["length_menu"]; !ok {
		options["length_menu"] = []string{"25", "50", "100", "-1"}
	}
	if _, ok := options["order"]; !ok {
		options["order"] = []any{2, "asc"}
	}
	order, err := json.Marshal(options["order"])
	if err != nil {
		return "", err
	}
	if _, ok := options["single_ops"]; !ok {
		options["single_ops"] = []map[string]any{}
	}
	if _, ok := options["group_ops"]; !ok {
		options["group_ops"] = []map[string]any{}
	}
	if _, ok := options["jscript"]; !ok {
		options["jscript"] = ""
	}

	return webui.Render(map[string]any{
		"id":         options["grid_id"],
		"cdn_url":    options["cdn_url"],
		"baseurl":    options["base_url"],
		"loadurl":    options["load_url"],
		"lenMenu":    options["length_menu"],
		"columns":    strings.Join(columns, ","),
		"order":      string(order),
		"single_ops": options["single_ops"],
		"group_ops":  options["group_ops"],
		"export":     export,
		"styles":     styles,
		"jscript":    options["jscript"],
	}, tplName)
}

func StdDataGridText(value, defaultVal, styles string, render map[string]any) map[string]any {
	if render == nil {
		render = stdRender
	}

	if defaultVal == "" {
		defaultVal = "-"
	}

	var data string
	if value == "'" {
		data = fmt.Sprintf(`<div class=""></div>`)
		value = ""
	} else if value != "" {
		data = fmt.Sprintf(`<div class="%s">%s</div>`, styles, value)
	} else {
		data = fmt.Sprintf(`<span class="text-black-50">%s</span>`, defaultVal)
	}

	return map[string]any{
		render["_"].(string):       value,
		render["display"].(string): data,
	}
}

func StdDataGridLink(value, url, label, styles string, render map[string]any) map[string]any {
	if render == nil {
		render = stdRender
	}

	if label != "" {
		value = label
	}

	return map[string]any{
		render["_"].(string): value,
		render["display"].(string): fmt.Sprintf(
			`<div class="%s"><a class="text-primary" href="%s">%s</a></div>`,
			styles, url, value),
	}
}

func StdDataGridPill(value, true_chk, styles string, render map[string]any) map[string]any {
	if render == nil {
		render = stdRender
	}

	var category string
	if value == true_chk {
		category = "success"
	} else {
		category = "danger"
	}

	return map[string]any{
		render["_"].(string): value,
		render["display"].(string): fmt.Sprintf(
			`<span class="badge badge-%s %s">%s</span>`,
			category, styles, value),
	}
}

func StdDataGridCheck(value bool, styles string, render map[string]any) map[string]any {
	if render == nil {
		render = stdRender
	}
	var textCat, iconCat string
	if value {
		textCat = "text-success"
		iconCat = "fa-check"
	} else {
		textCat = "text-danger"
		iconCat = "fa-times"
	}

	return map[string]any{
		render["_"].(string): value,
		render["display"].(string): fmt.Sprintf(
			`<span class="%s %s"><i class="fa fas fa-fw %s"></i></span>`,
			textCat, styles, iconCat),
	}
}
