package color

import "fmt"

const (
	textBlack = iota + 30
	textRed
	textGreen
	textYellow
	textBlue
	textPurple
	textCyan
	textWhite
)

func Red(s string) string {
	return textColor(textRed, s)
}

func Green(s string) string {
	return textColor(textGreen, s)
}

func textColor(c int, s string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", c, s)
}
