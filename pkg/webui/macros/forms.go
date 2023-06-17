package macros

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/exonlabs/go-utils/pkg/webui"
)

var (
	formPath = "templates/macros/forms/"
)

func UiInputForm(options map[string]any, styles string) (string, error) {
	tplName := formPath + "inputform.tpl"

	fields := []map[string]any{}
	if _, ok := options["fields"]; ok {
		for i, k := range options["fields"].([]map[string]any) {
			field := map[string]any{}
			if _, ok := k["type"]; !ok {
				field["type"] = "text"
			} else {
				field["type"] = k["type"]
			}
			if _, ok := k["label"]; !ok {
				field["label"] = fmt.Sprintf("Field %d", (i + 1))
			} else {
				field["label"] = k["label"]
			}
			if _, ok := k["name"]; !ok {
				field["name"] = fmt.Sprintf("field_%d", (i + 1))
			} else {
				field["name"] = k["name"]
			}
			if _, ok := k["value"]; !ok {
				field["value"] = ""
			} else {
				field["value"] = k["value"]
			}
			if _, ok := k["default"]; !ok {
				field["default"] = ""
			} else {
				field["default"] = k["default"]
			}
			if _, ok := k["format"]; !ok {
				field["format"] = ""
			} else {
				field["format"] = k["format"]
			}
			if _, ok := k["options"]; !ok {
				field["options"] = []map[string]any{}
			} else {
				field["options"] = k["options"]
			}
			if _, ok := k["multiple"]; !ok {
				field["multiple"] = false
			} else {
				field["multiple"] = k["multiple"]
			}
			if _, ok := k["rows"]; !ok {
				field["rows"] = 4
			} else {
				field["rows"] = k["rows"]
			}
			if _, ok := k["maxsize"]; !ok {
				field["maxsize"] = 0
			} else {
				field["maxsize"] = k["maxsize"]
			}
			if _, ok := k["placeholder"]; !ok {
				if _, ok := k["label"]; ok {
					field["placeholder"] = k["label"]
				} else {
					field["placeholder"] = ""
				}
			} else {
				field["placeholder"] = k["placeholder"]
			}
			if _, ok := k["required"]; !ok {
				field["required"] = false
			} else {
				field["required"] = k["required"]
			}
			if _, ok := k["confirm"]; !ok {
				field["confirm"] = false
			} else {
				field["confirm"] = k["confirm"]
			}
			if _, ok := k["strength"]; !ok {
				field["strength"] = false
			} else {
				field["strength"] = k["strength"]
			}
			if _, ok := k["help"]; !ok {
				field["help"] = ""
			} else {
				field["help"] = k["help"]
			}
			if _, ok := k["helpguide"]; !ok {
				field["helpguide"] = ""
			} else {
				field["helpguide"] = k["helpguide"]
			}
			if _, ok := k["prepend"]; !ok {
				field["prepend"] = []map[string]any{}
			} else {
				field["prepend"] = k["prepend"]
			}
			if _, ok := k["append"]; !ok {
				field["append"] = []map[string]any{}
			} else {
				field["append"] = k["append"]
			}

			fields = append(fields, field)
		}
	}

	// default data
	if _, ok := options["cdn_url"]; !ok {
		options["cdn_url"] = ""
	}
	if _, ok := options["form_id"]; !ok {
		options["form_id"] = webui.RandInt(10000)
	}
	if _, ok := options["submit_url"]; !ok {
		options["submit_url"] = ""
	}
	if _, ok := options["jscript"]; !ok {
		options["jscript"] = ""
	}

	return webui.Render(map[string]any{
		"id":         options["form_id"],
		"cdn_url":    options["cdn_url"],
		"submit_url": options["submit_url"],
		"fields":     fields,
		"styles":     styles,
		"jscript":    options["jscript"],
	}, tplName)
}

func UiQBuilder(options map[string]any, styles string) (string, error) {
	tplName := formPath + "qbuilder.tpl"

	var filtersJson string
	if _, ok := options["filters"].([]map[string]any); ok {
		filters := []string{}
		for _, k := range options["filters"].([]map[string]any) {
			filter := map[string]any{}
			filter["id"] = k["id"]
			filter["label"] = k["label"]
			if _, ok := k["type"]; !ok {
				filter["type"] = "string"
			} else {
				filter["type"] = k["type"]
			}
			if _, ok := k["input"]; !ok {
				filter["input"] = "text"
			} else {
				filter["input"] = k["input"]
			}
			if _, ok := k["operators"]; !ok {
				filter["operators"] = nil
			} else {
				filter["operators"] = k["operators"]
			}
			if _, ok := k["values"]; !ok {
				filter["values"] = nil
			} else {
				filter["values"] = k["values"]
			}
			if _, ok := k["default_value"]; !ok {
				filter["default_value"] = nil
			} else {
				filter["default_value"] = k["default_value"]
			}
			if _, ok := k["maxsize"]; !ok {
				filter["size"] = nil
			} else {
				filter["size"] = k["maxsize"]
			}
			if _, ok := k["rows"]; !ok {
				filter["rows"] = 3
			} else {
				filter["rows"] = k["rows"]
			}
			if _, ok := k["multiple"]; !ok {
				filter["multiple"] = false
			} else {
				filter["multiple"] = k["multiple"]
			}
			if _, ok := k["validation"]; !ok {
				filter["validation"] = nil
			} else {
				filter["validation"] = k["validation"]
			}

			j, err := json.Marshal(filter)
			if err != nil {
				return "", err
			}

			filters = append(filters, string(j))
		}

		filtersJson = strings.Join(filters, ",")
	}

	var rules []byte
	if val, ok := options["initial_rules"].(map[string]any); ok {
		var err error
		rules, err = json.Marshal(val)
		if err != nil {
			return "", err
		}
	}

	// default data
	if _, ok := options["cdn_url"]; !ok {
		options["cdn_url"] = ""
	}
	if _, ok := options["form_id"]; !ok {
		options["form_id"] = webui.RandInt(10000)
	}
	if _, ok := options["allow_groups"]; !ok {
		options["allow_groups"] = "1"
	}
	if _, ok := options["allow_empty"]; !ok {
		options["allow_empty"] = "true"
	}
	if _, ok := options["default_condition"]; !ok {
		options["default_condition"] = "AND"
	}
	if _, ok := options["inputs_separator"]; !ok {
		options["inputs_separator"] = ","
	}

	return webui.Render(map[string]any{
		"id":                options["form_id"],
		"cdn_url":           options["cdn_url"],
		"filters":           filtersJson,
		"rules":             string(rules),
		"allow_groups":      options["allow_groups"],
		"allow_empty":       options["allow_empty"],
		"default_condition": options["default_condition"],
		"inputs_separator":  options["inputs_separator"],
		"styles":            styles,
	}, tplName)
}

func UiLoginForm(options map[string]any, styles string) (string, error) {
	tplName := formPath + "loginform.tpl"

	// default data
	if _, ok := options["cdn_url"]; !ok {
		options["cdn_url"] = ""
	}
	if _, ok := options["form_id"]; !ok {
		options["form_id"] = webui.RandInt(10000)
	}
	if _, ok := options["submit_url"]; !ok {
		options["submit_url"] = ""
	}
	if _, ok := options["authkey"]; !ok {
		options["authkey"] = ""
	}
	if _, ok := options["btn_style"]; !ok {
		options["btn_style"] = "btn-primary"
	}

	return webui.Render(map[string]any{
		"id":         options["form_id"],
		"cdn_url":    options["cdn_url"],
		"submit_url": options["submit_url"],
		"authkey":    options["authkey"],
		"btn_style":  options["btn_style"],
		"styles":     styles,
	}, tplName)
}
