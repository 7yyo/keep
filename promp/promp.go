package promp

import "github.com/manifoldco/promptui"

func Select(option []string, label string, n int) *promptui.Select {
	return &promptui.Select{
		Label: label,
		Items: option,
		Size:  n,
	}
}
