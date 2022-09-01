package promp

import (
	"github.com/manifoldco/promptui"
	"strings"
)

func Select(o []string, l string, n int) *promptui.Select {
	p := &promptui.Select{
		Label: l,
		Items: o,
		Size:  n,
		Searcher: func(input string, index int) bool {
			pName := o[index]
			name := strings.Replace(strings.ToLower(pName), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)
			return strings.Contains(name, input)
		},
	}
	return p
}
