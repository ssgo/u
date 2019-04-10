package u

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type AttribCode int

const (
	AttrNone AttribCode = iota
	AttrBold
	AttrDim
	AttrItalic
	AttrUnderline
	AttrBlink
	AttrFastBlink
	AttrReverse
	AttrHidden
	AttrCrossedOut
)

type TextColor int

const (
	TextBlack TextColor = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
	TextNone TextColor = 0
)

type BgColor int

const (
	BgBlack BgColor = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
	BgNone  BgColor = 0
)

func Black(s string, attribs ...AttribCode) string {
	return Color(s, TextBlack, BgNone, attribs...)
}
func Red(s string, attribs ...AttribCode) string {
	return Color(s, TextRed, BgNone, attribs...)
}
func Green(s string, attribs ...AttribCode) string {
	return Color(s, TextGreen, BgNone, attribs...)
}
func Yellow(s string, attribs ...AttribCode) string {
	return Color(s, TextYellow, BgNone, attribs...)
}
func Blue(s string, attribs ...AttribCode) string {
	return Color(s, TextBlue, BgNone, attribs...)
}
func Magenta(s string, attribs ...AttribCode) string {
	return Color(s, TextMagenta, BgNone, attribs...)
}
func Cyan(s string, attribs ...AttribCode) string {
	return Color(s, TextCyan, BgNone, attribs...)
}
func White(s string, attribs ...AttribCode) string {
	return Color(s, TextWhite, BgNone, attribs...)
}
func Dim(s string) string {
	return Color(s, TextWhite, BgNone, AttrDim)
}
func Italic(s string) string {
	return Color(s, TextWhite, BgNone, AttrItalic)
}

func BBlack(s string, attribs ...AttribCode) string {
	return Color(s, TextWhite, BgBlack, attribs...)
}
func BRed(s string, attribs ...AttribCode) string {
	return Color(s, TextWhite, BgRed, attribs...)
}
func BGreen(s string, attribs ...AttribCode) string {
	return Color(s, TextBlack, BgGreen, attribs...)
}
func BYellow(s string, attribs ...AttribCode) string {
	return Color(s, TextBlack, BgYellow, attribs...)
}
func BBlue(s string, attribs ...AttribCode) string {
	return Color(s, TextWhite, BgBlue, attribs...)
}
func BMagenta(s string, attribs ...AttribCode) string {
	return Color(s, TextWhite, BgMagenta, attribs...)
}
func BCyan(s string, attribs ...AttribCode) string {
	return Color(s, TextBlack, BgCyan, attribs...)
}
func BWhite(s string, attribs ...AttribCode) string {
	return Color(s, TextBlack, BgWhite, attribs...)
}

func Color(s string, textColor TextColor, bgColor BgColor, attribs ...AttribCode) string {
	if runtime.GOOS == "windows" {
		return s
	}

	sets := make([]string, 0)
	for _, attr := range attribs {
		sets = append(sets, strconv.Itoa(int(attr)))
	}
	if textColor != TextNone {
		sets = append(sets, strconv.Itoa(int(textColor)))
	}
	if bgColor != BgNone {
		sets = append(sets, strconv.Itoa(int(bgColor)))
	}

	return fmt.Sprint("\033[", strings.Join(sets, ";"), "m", s, "\033[0m")
}
