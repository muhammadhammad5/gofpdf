package gofpdf

import (
	"math"
	//	"strings"
	"unicode"
)

// SplitText splits UTF-8 encoded text into several lines using the current
// font. Each line has its length limited to a maximum width given by w. This
// function can be used to determine the total height of wrapped text for
// vertical placement purposes.
func (f *Fpdf) SplitText(txt string, w float64) (lines []string) {
	cw := f.currentFont.Cw
	if cw == nil || len(cw) == 0 {
		// No width table? Fall back to a single line.
		return []string{txt}
	}

	wmax := int(math.Ceil((w - 2*f.cMargin) * 1000 / f.fontSize))

	// Convert to runes and trim trailing newlines (preserve behavior).
	s := []rune(txt)
	nb := len(s)
	for nb > 0 && s[nb-1] == '\n' {
		nb--
	}
	s = s[:nb]

	// Normalize runes for 8-bit width table:
	// - NBSP -> regular space
	// - any rune outside cw's range -> '?'
	for i, r := range s {
		switch r {
		case '\u00A0': // NBSP
			s[i] = ' '
		default:
			if int(r) < 0 || int(r) >= len(cw) {
				s[i] = '?'
			}
		}
	}

	sep := -1
	i := 0
	j := 0
	l := 0

	for i < len(s) {
		c := s[i]

		// Safe: c is normalized to be within cw's range.
		l += cw[int(c)]

		if unicode.IsSpace(c) || isChinese(c) {
			sep = i
		}

		if c == '\n' || l > wmax {
			if sep == -1 {
				if i == j {
					i++
				}
				sep = i
			} else {
				i = sep + 1
			}
			lines = append(lines, string(s[j:sep]))
			sep = -1
			j = i
			l = 0
		} else {
			i++
		}
	}

	if i != j {
		lines = append(lines, string(s[j:i]))
	}
	return lines
}
