package webui

type MenuLink struct {
	Label   string
	URL     string
	Icon    string
	SubMenu map[int]MenuLink
}

func AddMenulink(menuBuffer map[int]MenuLink, index int,
	label, icon, url string, parent int) {
	// standalone link
	if parent == 0 {
		menuBuffer[index] = MenuLink{
			Label:   label,
			Icon:    icon,
			URL:     url,
			SubMenu: make(map[int]MenuLink),
		}
	} else {
		// submenu link
		if _, ok := menuBuffer[parent]; ok {
			menuBuffer[parent].SubMenu[index] = MenuLink{Label: label, Icon: icon, URL: url}
		}
	}
}
