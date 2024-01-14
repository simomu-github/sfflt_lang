package formatter

import (
	"strings"
)

func FormatOneLine(instructions []string) string {
	return strings.Join(instructions, "")
}

func FormatSquere(instructions []string, column int) string {
	joinedInstructions := strings.Join(instructions, "")
	out := ""
	for i, c := range joinedInstructions {
		out = out + string(c)
		if (i+1)%column == 0 {
			out = out + "\n"
		}
	}

	return out
}

func FormatRaw(instructions []string) string {
	return strings.Join(instructions, "\n")
}
