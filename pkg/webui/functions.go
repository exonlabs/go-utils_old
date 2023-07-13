package webui

import (
	"strings"
	"text/template"

	"golang.org/x/exp/slices"
)

var tplFunc = template.FuncMap{
	"join":       strings.Join,
	"replaceAll": strings.ReplaceAll,
	"inStr":      inStr,
}

func inStr(s []string, v string) bool {
	return slices.Contains(s, v)
}
