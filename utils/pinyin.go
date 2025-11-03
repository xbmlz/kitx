package utils

import (
	"strings"

	"github.com/mozillazg/go-pinyin"
)

func GetPinyinFirstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}

	a := pinyin.NewArgs()
	a.Style = pinyin.FirstLetter
	a.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
	// [[z] [g] [r]] => zgr
	arr := pinyin.Pinyin(s, a)
	var builder strings.Builder
	for _, inner := range arr {
		if len(inner) > 0 {
			builder.WriteString(inner[0])
		}
	}
	return builder.String()
}

func GetPinyin(s string) string {
	if len(s) == 0 {
		return ""
	}

	a := pinyin.NewArgs()
	a.Style = pinyin.Normal
	a.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
	arr := pinyin.Pinyin(s, a)
	//  [[zhong] [guo] [ren]] => zhong guo ren
	var builder strings.Builder
	for _, inner := range arr {
		if len(inner) > 0 {
			builder.WriteString(inner[0])
			builder.WriteString(" ")
		}
	}
	return builder.String()
}
