// Copyright 05-Sep-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"strings"
	"unicode"
)

type statusT int

const (
	code statusT = iota
	quote
	long
)

var status = code
var line = 0
var ok = false
var props map[string]int
var charQuote = ""

func processPoint(c, rest string, ix int) int {
	ch := '\n'
	r := ix
	for ; r < len([]rune(c)); r++ {
		ch = []rune(c)[r]
		if !(unicode.IsDigit(ch) || unicode.IsLetter(ch) ||
			ch == ' ' || ch == '$') {
			break
		}
	}
	if r == len([]rune(c)) && strings.HasPrefix(rest, "/**/") {
		return r
	}
	if ch != '(' {
		prop := strings.TrimSpace(c[ix:r])
		if _, yes := props[prop]; !yes {
			props[prop] = line
		}
		r = ix + 1
	}
	return r
}

func processCode(c, rest string) {
	length := len(c)
	ix := strings.Index(c, ".")
	if ix != -1 {
		ix++
		if strings.HasPrefix(c[ix:], "_") ||
			strings.HasPrefix(c[ix:], "length") {
			processCode(c[ix+1:], rest)
		} else if strings.HasPrefix(c[ix:], "..") {
			processCode(c[ix+2:], rest)
		} else if strings.HasPrefix(c[ix:], "dedeme.") {
			processCode(c[ix+8:], rest)
		} else if unicode.IsDigit([]rune(c)[ix]) {
			processCode(c[ix+1:], rest)
		} else {
			ix2 := processPoint(c, rest, ix)
			if ix2 < length {
				processCode(c[ix2+1:], rest)
			}
		}
	}
}

func processLine(l string) {
	line++

	switch status {
	case quote:
		iBar := strings.Index(l, "\\")
		iQuote := strings.Index(l, charQuote)
		if iBar != -1 && (iQuote == -1 || iBar < iQuote) {
			line--
			processLine(l[iBar+2:])
		} else if iQuote != -1 {
			status = code
			line--
			processLine(l[iQuote+1:])
		} else {
			status = code
		}
	case long:
		iEnd := strings.Index(l, "*/")
		if iEnd != -1 {
			status = code
			line--
			processLine(l[iEnd+2:])
		}
	default:
		iLong := strings.Index(l, "/*")
		iShort := strings.Index(l, "//")
		iQuote := strings.Index(l, "\"")
		iQuote2 := strings.Index(l, "'")
		iQuote3 := strings.Index(l, "`")

		if iQuote == -1 {
			if iQuote2 == -1 {
				if iQuote3 != -1 {
					iQuote = iQuote3
				}
			} else {
				iQuote = iQuote2
				if iQuote3 != -1 && iQuote3 < iQuote2 {
					iQuote = iQuote3
				}
			}
		} else {
			if iQuote2 != -1 && iQuote2 < iQuote {
				iQuote = iQuote3
			}
			if iQuote3 != -1 && iQuote3 < iQuote {
				iQuote = iQuote3
			}
		}
		charQuote = "\""
		if iQuote == iQuote2 {
			charQuote = "'"
		} else if iQuote == iQuote3 {
			charQuote = "`"
		}

		sel := 0
		if iLong == -1 {
			if iShort == -1 {
				if iQuote != -1 {
					sel = 3
				}
			} else {
				sel = 2
				if iQuote != -1 && iQuote < iShort {
					sel = 3
				}
			}
		} else {
			sel = 1
			if iShort != -1 && iShort < iLong {
				sel = 2
				if iQuote != -1 && iQuote < iShort {
					sel = 3
				}
			} else {
				if iQuote != -1 && iQuote < iLong {
					sel = 3
				}
			}
		}

		switch sel {
		case 1:
			processCode(l[:iLong], l[iLong:])
			line--
			status = long
			processLine(l[iLong+2:])
		case 2:
			processCode(l[:iShort], l[iShort:])
		case 3:
			processCode(l[:iQuote], l[iQuote:])
			line--
			status = quote
			processLine(l[iQuote+1:])
		default:
			processCode(l, "")
		}
	}

}
