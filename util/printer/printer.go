package printer

import "keep/util/color"

func IsNil(i interface{}) interface{} {
	if i == nil {
		return ""
	}
	return i
}

func Status(s string) string {
	switch s {
	case "normal":
		return color.Green(s)
	case "stop":
		return color.Yellow(s)
	default:
		return s
	}
}
