package web

import (
	"os"
	"strings"
)

var wdg string

func GetK2Path() string {
	if wdg != "" {
		return wdg
	}
	wd, _ := os.Getwd()
	if strings.Contains(wd, "\\k2web") || strings.Contains(wd, "\\k2web") {
		wdg = wd[:len(wd)-6]
	} else {
		wdg = wd
	}
	return wdg
}
