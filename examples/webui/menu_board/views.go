package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"github.com/exonlabs/go-utils/pkg/webapp"
	"github.com/exonlabs/go-utils/pkg/webui"
	"github.com/exonlabs/go-utils/pkg/webui/macros"
)

const cdnURL = "/static/vendor"

var (
	sharedBuff sync.Map
	count      int
)

var IndexView = &webapp.WebView{
	Name:   "WebUI",
	Routes: []string{"/", "/index"},
	Initialize: func(wv *webapp.WebView) {
		webui.AddMenulink(menuBuffer, 1, "UI Components",
			"fa-cubes", "", 0)
	},
	GET: func(wv *webapp.WebView) (string, error) {
		params := map[string]any{
			"cdn_url": cdnURL,
			"langs": map[string]string{
				"en": "English",
				"ar": "العربية",
				"fr": "Français",
			},
			"doc_title": "WebUI",
			"menu":      menuBuffer,
		}

		return webui.Render(params,
			"templates/html.tpl",
			"templates/menuboard.tpl",
			"templates/mainpage.tpl")
	},
}

var HomeView = &webapp.WebView{
	Name:   "Home",
	Routes: []string{"/home", "/home/"},
	Initialize: func(wv *webapp.WebView) {
		webui.AddMenulink(menuBuffer, 0, "Home",
			"fa-home", "#home", 0)
	},
	GET: func(wv *webapp.WebView) (string, error) {
		params := map[string]any{
			"langs": map[string]string{
				"en": "English",
				"ar": "العربية",
				"fr": "Français",
			},
			"message": "Welcome",
		}

		html, err := webui.Render(params,
			"templates/option_panel.tpl")
		if err != nil {
			return "", err
		}

		return webui.Reply(wv, html, "Home", nil), nil
	},
}

var NotifyView = &webapp.WebView{
	Routes: []string{"/notify", "/notify/"},
	Initialize: func(wv *webapp.WebView) {
		webui.AddMenulink(menuBuffer, 1, "Notifications",
			"", "#notify", 1)
	},
	GET: func(wv *webapp.WebView) (string, error) {
		html, err := macros.UiAlert("general message", "showing notifications",
			true, false, "p-3")
		if err != nil {
			return "", err
		}
		html += "<script>WebUI.board_menu.show_submenu(1)</script>"

		return webui.Reply(wv, html, "Notifications", nil), nil
	},
}

var AlertsView = &webapp.WebView{
	Routes: []string{"/alerts"},
	Initialize: func(wv *webapp.WebView) {
		webui.AddMenulink(menuBuffer, 2, "Alerts",
			"", "#alerts", 1)
	},
	GET: func(wv *webapp.WebView) (string, error) {
		data, err := macros.UiAlert("info", "info message", true, true, "px-3 pt-3")
		if err != nil {
			return "", err
		}
		html := data
		data, err = macros.UiAlert("warn", "warning message", true, true, "px-3")
		if err != nil {
			return "", err
		}
		html += data
		data, err = macros.UiAlert("error", "error message", true, true, "px-3")
		if err != nil {
			return "", err
		}
		html += data
		data, err = macros.UiAlert("success", "success message", true, true, "px-3")
		if err != nil {
			return "", err
		}
		html += data
		data, err = macros.UiAlert("message", "message", true, true, "px-3")
		if err != nil {
			return "", err
		}
		html += data
		html += "<script>WebUI.board_menu.show_submenu(1)</script>"

		return webui.Reply(wv, html, "Alerts", nil), nil
	},
}

var InputForm = &webapp.WebView{
	Routes: []string{"/inputform", "/inputform/"},
	Initialize: func(wv *webapp.WebView) {
		webui.AddMenulink(menuBuffer, 3, "Input Form",
			"", "#inputform", 1)
	},
	GET: func(wv *webapp.WebView) (string, error) {
		options := map[string]any{
			"cdn_url":    cdnURL,
			"form_id":    "1234",
			"submit_url": "/inputform/",
			"fields": []map[string]any{
				{
					"type": "checkbox", "label": "Validation",
					"options": []map[string]any{
						{"label": "Server side validation", "name": "validation"},
					},
				},
				{"type": "title", "label": "Group Label"},
				{"type": "text", "label": "Required Field",
					"name": "field1", "required": true,
					"help":      "* example with input append",
					"helpguide": "Extra detailed long help for fields",
					"append": []map[string]any{
						{"type": "text", "value": ".00"},
						{"type": "text", "value": "$"},
					},
				},
				{"type": "text", "label": "Optional Field",
					"name": "field2", "required": false,
					"help":      "* example with input prepend",
					"helpguide": "Extra detailed long help for fields",
					"prepend": []map[string]any{
						{"type": "icon", "value": "fa-phone"},
						{"type": "text", "value": "+00"},
					},
				},
				{"type": "text", "label": "Optional Field",
					"name": "field3", "required": false,
					"help":      "* help text for field",
					"helpguide": "Extra detailed long help for field",
					"append": []map[string]any{
						{"type": "select", "options": []map[string]any{
							{"label": ".com", "value": ".com"},
							{"label": ".net", "value": ".net", "selected": true},
							{"label": ".org", "value": ".org"}}}},
				},
				{"type": "textarea", "label": "Textarea",
					"name":      "field4",
					"help":      "* help text for field",
					"helpguide": "Extra detailed long help for fields"},

				{"type": "title", "label": "Group Label"},
				{"type": "password", "label": "Password 1",
					"name": "pass1", "required": true, "strength": true},
				{"type": "password", "label": "Password 2",
					"name": "pass2", "required": true, "confirm": true},
				{"type": "password", "label": "Password 3",
					"name": "pass3", "required": true, "strength": true,
					"confirm": true},

				{"type": "title", "label": "Group Label"},
				{"type": "select", "label": "Select",
					"name": "select1", "required": true,
					"options": []map[string]any{{"label": "Select", "value": nil},
						{"label": "Option 1", "value": "01"},
						{"label": "Option 2", "value": "02"},
						{"label": "Option 3", "value": "03"},
						{"label": "Option 4", "value": "04"},
						{"label": "Option 5", "value": "05"}}},
				{"type": "select", "label": "Select multiple",
					"name": "select2", "required": true, "multiple": true,
					"options": []map[string]any{{"label": "Option 1", "value": "01",
						"selected": true},
						{"label": "Option 2", "value": "02",
							"selected": true},
						{"label": "Option 3", "value": "03"},
						{"label": "Option 4", "value": "04"},
						{"label": "Option 5", "value": "05"}}},

				{"type": "title", "label": "Group Label"},
				{"type": "checkbox", "label": "Checkbox",
					"helpguide": "Extra detailed long help for fields",
					"options": []map[string]any{{"label": "Select 1",
						"name": "check1", "selected": true},
						{"label": "Select 2",
							"name": "check2"}}},
				{"type": "radio", "label": "Radio",
					"name": "radio1", "required": true,
					"helpguide": "Extra detailed long help for fields",
					"options": []map[string]any{{"label": "Option 1", "value": "1"},
						{"label": "Option 2", "value": "2"},
						{"label": "Option 3", "value": "3"}}},

				{"type": "title", "label": "Group Label"},
				{"type": "datetime", "label": "Date & Time",
					"name": "date1", "required": true},
				{"type": "date", "label": "Date", "name": "date2"},
				{"type": "time", "label": "Time", "name": "time1"},

				{"type": "title", "label": "Group Label"},
				{"type": "file", "label": "Upload File",
					"name": "files1", "required": true,
					"format": ".txt,.pdf,.png",
					"help":   fmt.Sprintf("* %s: <span dir=\"ltr\">.txt, .pdf, .png</span>", "allowed types")},
				{"type": "file", "label": "Upload Multiple",
					"name": "files2", "required": false, "multiple": true,
					"placeholder": "", "maxsize": 1048576,
					"help":      "* all types allowed, max file size: 1MB",
					"helpguide": "Extra detailed long help for fields"},
			},
		}

		form, err := macros.UiInputForm(options, "")
		if err != nil {
			return "", err
		}
		html, err := webui.Render(map[string]any{"contents": form},
			"templates/input_form.tpl")
		if err != nil {
			return "", err
		}
		html += "<script>WebUI.board_menu.show_submenu(1)</script>"

		return webui.Reply(wv, html, "Input Form", nil), nil
	},
	POST: func(wv *webapp.WebView) (string, error) {
		validation := wv.Env.Request.FormValue("validation")
		if validation == "1" {
			params := map[string]any{
				"validation": []string{"field1", "date1", "files1"},
			}

			return webui.Reply(wv, "", "", params), nil
		}

		msg := fmt.Sprintf(`%s<br><div dir="ltr" style="text-align:left">`,
			"Submited Data:")
		for k := range wv.Env.Request.Form {
			if k == "_csrf_token" {
				continue
			}
			v := wv.Env.Request.FormValue(k)
			msg += fmt.Sprintf("<b>%s:</b> %s<br>", k, v)
		}
		for k := range wv.Env.Request.MultipartForm.File {
			_, fh, err := wv.Env.Request.FormFile(k)
			if err != nil {
				return "", err
			}
			msg += fmt.Sprintf("<b>%s:</b> %s<br>",
				fh.Filename, fh.Header.Get("Content-Type"))
		}

		msg += "</div>"
		msg = strings.ReplaceAll(msg, "'", "")
		return webui.Notify(wv, msg, "success", false, true, nil), nil
	},
}

var DatagridView = &webapp.WebView{
	Routes: []string{"/datagrid", "/datagrid/"},
	Initialize: func(wv *webapp.WebView) {
		webui.AddMenulink(menuBuffer, 4, "Datagrid",
			"", "#datagrid", 1)
	},
	GET: func(wv *webapp.WebView) (string, error) {
		options := map[string]any{
			"cdn_url":     cdnURL,
			"grid_id":     "1234",
			"base_url":    "/datagrid",
			"load_url":    "/datagrid/loaddata",
			"length_menu": []string{"10", "50", "100", "250", "-1"},
			"columns": []map[string]any{
				{"id": "field1", "title": "Field Name 1"},
				{"id": "field2", "title": "Field Name 2"},
				{"id": "field3.item1", "title": "Field3_1"},
				{"id": "field3.item2", "title": "Field3_2"},
				{"id": "field4", "title": "Other Field 4"},
				{"id": "field5", "title": "Extra Field 5",
					"visible": false},
				{"id": "field6", "title": "Data Field 6",
					"visible": false},
				{"id": "field7", "title": "Field Header 7",
					"visible": false},
				{"id": "field8", "title": "Field8",
					"visible": false},
				{"id": "field9", "title": "Field Number 9",
					"visible": false},
			},
			"export": map[string]any{
				"types":              []string{"csv", "xls", "print"},
				"file_title":         "Example Data",
				"file_prefix":        "export",
				"csv_fieldSeparator": ";",
				"csv_fieldBoundary":  "",
			},
			"single_ops": []map[string]any{
				{"label": "Single Operation 1", "action": "single_op1"},
				{"label": "Single Op 2 with confirm", "action": "single_op2",
					"confirm": "Are you sure you want to do this operation?"},
			},
			"group_ops": []map[string]any{
				{"label": "Group Operation 1", "action": "group_op1"},
				{"label": "Group Op 2 with confirm", "action": "group_op2",
					"confirm": "Are you sure?"},
				{"label": "Op 3 with Reload", "action": "group_op3"},
			},
		}
		dataGrid, err := macros.UiStdDataGrid(options, "")
		if err != nil {
			return "", err
		}
		html, err := webui.Render(map[string]any{"contents": dataGrid},
			"templates/data_grid.tpl")
		if err != nil {
			return "", err
		}
		html += "<script>WebUI.board_menu.show_submenu(1)</script>"

		return webui.Reply(wv, html, "Datagrid", nil), nil
	},
	POST: func(wv *webapp.WebView) (string, error) {
		action := strings.Split(
			strings.TrimPrefix(wv.Env.Request.URL.Path, "/"), "/")[1]

		if action == "loaddata" {
			data := []map[string]any{}
			for k := 0; k < 228; k++ {
				_k := fmt.Sprintf("%03d", k)
				d := map[string]any{
					"DT_RowId": fmt.Sprintf("rowid_%s", _k),
					"field1": macros.StdDataGridLink(
						fmt.Sprintf("master_%s", _k), "#datagrid", "", "", nil),
					"field2": macros.StdDataGridText(
						fmt.Sprintf("field2 %s", _k), "", "", nil),
				}

				item := map[string]any{}
				if rand.Intn(3) == 0 {
					item["item1"] = macros.StdDataGridText(
						fmt.Sprintf("field3.1 %s", _k), "", "", nil)
				} else {
					item["item1"] = macros.StdDataGridText("'", "", "", nil)
				}
				if rand.Intn(4) == 0 {
					item["item2"] = macros.StdDataGridText(
						fmt.Sprintf("field3.2 %s", _k), "", "", nil)
				} else {
					item["item2"] = macros.StdDataGridText("'", "", "", nil)
				}
				d["field3"] = item

				var pill string
				if rand.Intn(2) == 0 {
					pill = "Yes"
				} else {
					pill = "No"
				}
				d["field4"] = macros.StdDataGridPill(pill, "Yes", "", nil)

				d["field5"] = macros.StdDataGridCheck(rand.Intn(2) == 0, "", nil)

				if rand.Intn(2) == 0 {
					d["field6"] = macros.StdDataGridText(
						fmt.Sprintf("field6 %s", _k), "", "", nil)
				} else {
					d["field6"] = macros.StdDataGridText("'", "", "", nil)
				}

				if rand.Intn(2) == 0 {
					d["field7"] = macros.StdDataGridText(
						fmt.Sprintf("field7 %s", _k), "", "", nil)
				} else {
					d["field7"] = macros.StdDataGridText("'", "", "", nil)
				}

				if rand.Intn(2) == 0 {
					d["field8"] = macros.StdDataGridText(
						fmt.Sprintf("field8 %s", _k), "", "", nil)
				} else {
					d["field8"] = macros.StdDataGridText("'", "", "", nil)
				}

				if rand.Intn(2) == 0 {
					d["field9"] = macros.StdDataGridText(
						fmt.Sprintf("field9 %s", _k), "", "", nil)
				} else {
					d["field9"] = macros.StdDataGridText("'", "", "", nil)
				}

				data = append(data, d)
			}
			return webui.Reply(wv, "", "", map[string]any{"payload": data}), nil
		}

		if slices.Contains([]string{"single_op1", "single_op2",
			"group_op1", "group_op2", "group_op3"}, action) {
			wv.Env.Request.ParseForm()
			rows := wv.Env.Request.Form["items[]"]
			if len(rows) > 20 {
				rows = rows[:21]
				rows[20] = "..."
			}
			msg := `<span dir="ltr" style="text-align:left">`
			msg += fmt.Sprintf("Operation: %s<br>", action)
			msg += fmt.Sprintf("Rows: %s", rows)
			msg += "</span>"
			msg = strings.ReplaceAll(msg, "'", "")
			return webui.Notify(wv, msg, "success", false, true, nil), nil
		}

		return webui.Notify(wv, "Invalid request", "error", false, true, nil), nil
	},
}

var QueryBuilderView = &webapp.WebView{
	Routes: []string{"/qbuilder"},
	Initialize: func(wv *webapp.WebView) {
		webui.AddMenulink(menuBuffer, 5, "Query Builder",
			"", "#qbuilder", 1)
	},
	GET: func(wv *webapp.WebView) (string, error) {
		options := map[string]any{
			"cdn_url": cdnURL,
			"form_id": "1234",
			"filters": []map[string]any{
				{"id": "field1", "label": "Field 1", "type": "string",
					"operators": []string{"equal", "not_equal", "contains"}},
				{"id": "field2", "label": "Field Name 2", "type": "string",
					"input": "textarea"},
				{"id": "field3_1", "label": "Integer 1", "type": "integer",
					"input": "text"},
				{"id": "field3_2", "label": "Integer 2", "type": "integer",
					"input": "number"},
				{"id": "field3_3", "label": "Double 1", "type": "double",
					"input": "number"},
				{"id": "field4", "label": "Select", "type": "integer",
					"input":     "select",
					"values":    map[int]string{1: "Option 1", 2: "Option 2", 3: "Option 3"},
					"operators": []string{"equal", "not_equal", "in", "not_in"}},
				{"id": "field5", "label": "Checkbox", "type": "integer",
					"input": "radio", "values": []map[int]any{{1: "Yes"}, {0: "No"}},
					"operators": []string{"equal"}},
				{"id": "field6", "label": "Choose", "type": "integer",
					"input": "checkbox",
					"values": []map[int]any{{1: "Opt 1"}, {2: "Opt 2"}, {3: "Opt 3"},
						{4: "Opt 4"}, {5: "Opt 5"}}},
			},
			"initial_rules": map[string]any{
				"not":       true,
				"condition": "AND",
				"rules": []map[string]any{
					{"id": "field1", "operator": "equal", "value": "value"},
					{"id": "field3_2", "operator": "less", "value": 10},
					{
						"condition": "OR",
						"rules": []map[string]any{
							{"id": "field1", "operator": "equal",
								"value": "value2"},
							{"id": "field6", "operator": "equal",
								"value": 2},
						},
					},
					{
						"not":       true,
						"condition": "OR",
						"rules": []map[string]any{
							{"id": "field4", "operator": "equal",
								"value": 3},
							{"id": "field2", "operator": "not_equal",
								"value": "text value"},
						},
					},
				},
			},
		}
		queryBuilder, err := macros.UiQBuilder(options, "")
		if err != nil {
			return "", err
		}
		html, err := webui.Render(map[string]any{"contents": queryBuilder},
			"templates/query_builder.tpl")
		if err != nil {
			return "", err
		}
		html += "<script>WebUI.board_menu.show_submenu(1)</script>"

		return webui.Reply(wv, html, "Query Builder", nil), nil
	},
}

var LoaderView = &webapp.WebView{
	Routes: []string{"/loader"},
	Initialize: func(wv *webapp.WebView) {
		webui.AddMenulink(menuBuffer, 2, "Page Loader",
			"", "#loader", 0)
	},
	GET: func(wv *webapp.WebView) (string, error) {
		html, err := macros.UiAlert("message", "loaded after delay",
			false, false, "p-3")
		if err != nil {
			return "", err
		}
		render, err := webui.Render(nil,
			"templates/progress_loader.tpl")
		if err != nil {
			return "", err
		}
		html += render

		// simulate delay
		time.Sleep(time.Second * 1)

		return webui.Reply(wv, html, "Page Loader", nil), nil
	},
	POST: func(wv *webapp.WebView) (string, error) {
		var wView webapp.WebView
		if count == 0 {
			wView = *wv
		}
		count++

		wv.Env.Request.ParseForm()
		// get loading progress status
		if wv.Env.Request.Form.Get("get_progress") != "" {
			res, _ := sharedBuff.Load("loader_progress")
			return webui.Reply(wv, strconv.Itoa(res.(int)), "", nil), nil
		}

		// simulate long delay
		t := 5
		for i := 1; i < t; i++ {
			sharedBuff.Store("loader_progress", i*100/t)
			time.Sleep(time.Duration(float64(time.Second) * 1))
		}
		sharedBuff.Delete("loader_progress")

		count = 0
		return webui.Notify(&wView, "success", "success",
			false, false, nil), nil
	},
}

var LoginView = &webapp.WebView{
	Routes: []string{"/loginpage", "/loginpage/"},
	Initialize: func(wv *webapp.WebView) {
		webui.AddMenulink(menuBuffer, 3, "Login Page",
			"", "loginpage", 0)
	},
	GET: func(wv *webapp.WebView) (string, error) {
		action := strings.Split(
			strings.TrimPrefix(wv.Env.Request.URL.Path, "/"), "/")[1]

		if action == "load" {
			html, err := macros.UiLoginForm(map[string]any{
				"cdn_url":    cdnURL,
				"submit_url": "/loginpage/",
				"authkey":    "123456",
			}, "text-white bg-secondary")
			if err != nil {
				return "", err
			}

			return webui.Reply(wv, html, "Loginpage", nil), nil
		} else {
			params := map[string]any{
				"cdn_url":   cdnURL,
				"doc_title": "WebUI",
				"langs": map[string]string{
					"en": "English",
					"ar": "العربية",
					"fr": "Français",
				},
				"load_url": "loginpage/load",
			}

			html, err := webui.Render(params,
				"templates/html.tpl",
				"templates/simplepage.tpl",
				"templates/loginpage.tpl")
			if err != nil {
				return "", err
			}

			return html, nil
		}
	},
	POST: func(wv *webapp.WebView) (string, error) {
		wv.Env.Request.ParseForm()
		username := wv.Env.Request.Form.Get("username")
		authdigest := wv.Env.Request.Form.Get("digest")

		var errMsg string
		if username == "" || authdigest == "" {
			errMsg = "Please enter username and password"
		} else {
			if username == "admin" {
				return webui.Notify(wv, "Welcome, admin", "success", false, true, nil), nil
			} else {
				errMsg = "Invalid username or password"
			}
		}
		return webui.Notify(wv, errMsg, "error", false, true, nil), nil
	},
}
