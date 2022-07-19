package printer

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"keep/util/color"
	"keep/util/sys"
)

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
	case "stopped":
		return color.Red(s)
	default:
		return s
	}
}

func Confirm() bool {
	fmt.Println()
	p := promptui.Prompt{
		Label:     "return",
		IsConfirm: true,
	}
	result, _ := p.Run()
	if result != "y" {
		sys.Exit()
	}
	return true
}

func PrintError(e string) {
	fmt.Println(color.Red(e))
}
