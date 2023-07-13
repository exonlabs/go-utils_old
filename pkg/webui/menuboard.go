package webui

// type MenuLink struct {
// 	Label   string
// 	URL     string
// 	Icon    string
// 	SubMenu map[int]MenuLink
// }

// func AddMenulink(menuBuffer map[int]MenuLink, index int,
// 	label, icon, url string, parent int) {
// 	// standalone link
// 	if parent == 0 {
// 		menuBuffer[index] = MenuLink{
// 			Label:   label,
// 			Icon:    icon,
// 			URL:     url,
// 			SubMenu: make(map[int]MenuLink),
// 		}
// 	} else {
// 		// submenu link
// 		if _, ok := menuBuffer[parent]; ok {
// 			menuBuffer[parent].SubMenu[index] = MenuLink{Label: label, Icon: icon, URL: url}
// 		}
// 	}
// }

// // //go:embed templates/* templates/macros/*
// // var Tpl embed.FS

// func Redirect(wv *webapp.WebView, url string, blank bool) string {
// 	redirect := map[string]any{
// 		"redirect": url,
// 	}

// 	if blank {
// 		redirect["blank"] = blank
// 	}

// 	b, err := json.Marshal(redirect)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	wv.Env.Response.SetHeader("Content-Type", "application/json")
// 	return string(b)
// }

// func Reply(wv *webapp.WebView, response string, doctitle string, params any) string {
// 	if !wv.IsJsRequest() {
// 		return response
// 	}

// 	jsonMap := make(map[string]any)
// 	if response != "" {
// 		jsonMap["payload"] = response
// 	}
// 	if doctitle != "" {
// 		jsonMap["doctitle"] = doctitle
// 	}

// 	var err error
// 	var b []byte
// 	if len(jsonMap) > 0 {
// 		b, err = json.Marshal(jsonMap)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}

// 	if params != nil {
// 		b, err = json.Marshal(params)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}

// 	wv.Env.Response.SetHeader("Content-Type", "application/json")
// 	return string(b)
// }

// func Notify(wv *webapp.WebView, massage, category string,
// 	unique, sticky bool, params any) string {
// 	if params == nil {
// 		newParams := map[string]any{
// 			"notifications": []any{
// 				[]any{category, massage, unique, sticky},
// 			},
// 		}

// 		return Reply(wv, "", "", newParams)
// 	}

// 	return Reply(wv, "", "", nil)
// }

// func RandInt(index int) string {
// 	min := 1
// 	max := index
// 	randVal := rand.Intn(max-min) + min
// 	return strconv.Itoa(randVal)
// }
