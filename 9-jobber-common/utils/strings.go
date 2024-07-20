package utils

import (
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func FirstLetterUpperCase(s string) string {
	return cases.Title(language.AmericanEnglish, cases.Compact).String(s)
}

func LowerCase(s string) string {
	return cases.Lower(language.AmericanEnglish, cases.Compact).String(s)
}

func FirstToLower(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s
	}
	lc := unicode.ToLower(r)
	if r == lc {
		return s
	}
	return string(lc) + s[size:]
}
