package printer

import "keep/util/color"

func IsNil(i interface{}) interface{} {
	if i == nil {
		return ""
	}
	return i
}

func Status(s string) string {
	if s == "normal" {
		return color.Green("NORMAL")
	} else {
		return color.Red("ERROR")
	}
}
