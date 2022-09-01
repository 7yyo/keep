package printer

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"keep/util/color"
)

func Welcome() {
	fmt.Println("\n" +
		" __   ___  _______   _______    _______   \n" +
		"|/\"| /  \")/\"     \"| /\"     \"|  |   __ \"\\  \n" +
		"(: |/   /(: ______)(: ______)  (. |__) :) \n" +
		"|    __/  \\/    |   \\/    |    |:  ____/  \n" +
		"(// _  \\  // ___)_  // ___)_   (|  /      \n" +
		"|: | \\  \\(:      \"|(:      \"| /|__/ \\     \n" +
		"(__|  \\__)\\_______) \\_______)(_______)    \n" +
		"                                          ")
}

func Bye() {
	fmt.Println("\nbye!")
}

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
	p := promptui.Prompt{
		Label:     "return",
		IsConfirm: true,
		Default:   "y",
		Validate: func(input string) error {
			if input != "y" && input != "Y" {
				return fmt.Errorf("invalid input, please enter `y` or `Y`")
			}
			return nil
		},
	}
	_, _ = p.Run()
	return true
}

func PrintError(e string) {
	fmt.Printf(color.Red(e))
}

func Return() string {
	return color.Green("<")
}
