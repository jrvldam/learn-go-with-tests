package iteration

import "strings"

func Repeat(character string, times int) string {
	if times < 1 {
		times = 1
	}

	return strings.Repeat(character, times)
}
