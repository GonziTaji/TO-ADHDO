package funcmap

import (
	"slices"
	"strconv"
	"strings"
)

func format_clp(value int) string {
	source_str := strings.TrimSpace(strconv.Itoa(value))

	source_runes := reverseRunes([]rune(source_str))
	out_runes := []rune{}

	for i, char := range source_runes {
		if i != 0 && i%3 == 0 {
			out_runes = append(out_runes, '.')
		}

		out_runes = append(out_runes, char)
	}

	out_runes = append(out_runes, '$')

	return string(reverseRunes(out_runes))
}

func reverseRunes(values []rune) []rune {
	runes := slices.Clone(values)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return runes
}
