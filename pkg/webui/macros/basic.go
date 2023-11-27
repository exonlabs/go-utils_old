package macros

import (
	"math/rand"
	"strconv"

	"github.com/exonlabs/go-utils/pkg/webui"
)

var (
	basicPath = "templates/macros/basic/"
)

func UiAlert(notifyType, msg string, icon, dismiss bool, styles string) (string, error) {
	tplName := basicPath + "alert.tpl"

	var alertType, alertIcon string
	switch notifyType {
	case "error":
		alertType, alertIcon = "danger", "fa-exclamation-circle"
	case "warn":
		alertType, alertIcon = "warning", "fa-exclamation-circle"
	case "info":
		alertType, alertIcon = "info", "fa-info-circle"
	case "success":
		alertType, alertIcon = "success", "fa-check-circle"
	default:
		alertType, alertIcon = "secondary", "fa-check-circle"
	}

	if !icon {
		alertIcon = ""
	}

	return webui.Render(map[string]any{
		"alert_type":  alertType,
		"alert_icon":  alertIcon,
		"message":     msg,
		"dismissible": dismiss,
		"styles":      styles,
	}, tplName)
}

func RandInt(index int) string {
	min := 1
	max := index
	randVal := rand.Intn(max-min) + min
	return strconv.Itoa(randVal)
}
