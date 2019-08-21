package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

// buildBalloon takes a slice of strings of max width maxwidth
// prepends/appends margins on first and last line, and at start/end of each line
// and returns a string with the contents of the balloon
func buildBalloon(lines []string, maxWidth int) string {
	var borders []string
	var ret []string
	count := len(lines)

	borders = []string{"/", "\\", "\\", "/", "|", "<", ">"}

	top := " " + strings.Repeat("_", maxWidth+2)
	bottom := " " + strings.Repeat("-", maxWidth+2)

	ret = append(ret, top)
	if count == 1 {
		ret = append(ret, fmt.Sprintf("%s %s %s", borders[5], lines[0], borders[6]))
	} else {
		ret = append(ret, fmt.Sprintf(`%s %s %s`, borders[0], lines[0], borders[1]))
		i := 1
		for ; i < count-1; i++ {
			ret = append(ret, fmt.Sprintf(`%s %s %s`, borders[4], lines[i], borders[4]))
		}
		ret = append(ret, fmt.Sprintf(`%s %s %s`, borders[2], lines[i], borders[3]))
	}

	ret = append(ret, bottom)
	return strings.Join(ret, "\n")
}

// tabsToSpaces converts all tabs found in the strings
// found in the `lines` slice to 4 spaces, to prevent misalignment in
// counting the runes
func tabsToSpaces(lines []string) []string {
	var ret []string
	for _, l := range lines {
		l = strings.Replace(l, "\t", "    ", -1)
		ret = append(ret, l)
	}
	return ret
}

// normalizeStringsLength takes a slice of strings and appends
// to each one a number of spaces needed to have them all the same number
// of runes
func normalizeStringsLength(lines []string, maxwidth int) []string {
	var ret []string
	for _, l := range lines {
		s := l + strings.Repeat(" ", maxwidth-utf8.RuneCountInString(l))
		ret = append(ret, s)
	}
	return ret
}

// calculatemaxwidth given a slice of strings returns the length of the
// string with max length
func calculateMaxWidth(lines []string) int {
	w := 0
	for _, l := range lines {
		len := utf8.RuneCountInString(l)
		if len > w {
			w = len
		}
	}

	return w
}

// setNewLine take a slice of strings and split the string
// based on the maximum width
// then return the slice of splitted string
func setNewLine(lines []string, maxwidth int) []string {
	var newLines []string
	for _, l := range lines {
		if len(l) > maxwidth {
			sub := ""

			runes := bytes.Runes([]byte(l))
			l := len(runes)

			for i, r := range runes {
				sub = sub + string(r)
				if (i+1)%maxwidth == 0 {
					if runes[i+1] != rune(' ') && runes[i] != rune(' ') {
						sub = sub + "-"
					}

					newLines = append(newLines, sub)
					sub = ""
				} else if (i + 1) == l {
					newLines = append(newLines, sub)
				}
			}
		} else {
			newLines = append(newLines, l)
		}
	}
	return newLines
}

// Generate the cow
func generateCow(text string) string {
	lines := []string{fmt.Sprint(text)}

	var cow = `     \  ^__^
      \ (oo)\_______
        (__)\       )\/\
	        ||----w |
	        ||     ||
		`

	lines = setNewLine(lines, 30)
	lines = tabsToSpaces(lines)
	maxWidth := calculateMaxWidth(lines)
	messages := normalizeStringsLength(lines, maxWidth)
	balloon := buildBalloon(messages, maxWidth)
	return fmt.Sprintf("```\n%s\n%s\n```", balloon, cow)
}