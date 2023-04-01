package helper

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func ExtLess(value string) string {
	return value[:len(value)-len(filepath.Ext(value))]
}

func Capitalize(value string) string {
	c := strings.ToUpper(value[:1])
	t := value[1:]
	return fmt.Sprintf("%s%s", c, t)
}

func Colorize(value string, col string) string {
	var cl color.Attribute
	switch strings.ToLower(col) {
	case "blue":
		cl = color.FgBlue
	case "red":
		cl = color.FgRed
	case "yellow":
		cl = color.FgYellow
	case "green":
		cl = color.FgGreen
	case "cyan":
		cl = color.FgCyan
	case "magenta":
		cl = color.FgMagenta
	default:
		cl = color.FgWhite
	}

	fg := color.New(cl, color.Bold)
	return fg.Sprintf(value)
}
