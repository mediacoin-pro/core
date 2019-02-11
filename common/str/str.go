package str

import (
	"strings"
	"unicode"
)

func InSlice(ss []string, str string) bool {
	return IndexOf(ss, str) >= 0
}

func IndexOf(ss []string, str string) int {
	for i, s := range ss {
		if s == str {
			return i
		}
	}
	return -1
}

func Revert(ss []string) []string {
	for i, j := 0, len(ss)-1; i < j; i++ {
		ss[i], ss[j] = ss[j], ss[i]
		j--
	}
	return ss
}

func Unshift(ss []string, v ...string) []string {
	return append(v, ss...)
}

func Unique(ss []string) (res []string) {
	mm := map[string]bool{}
	for _, s := range ss {
		if !mm[s] {
			mm[s] = true
			res = append(res, s)
		}
	}
	return
}

func Exclude(ss []string, val string) (res []string) {
	for _, s := range ss {
		if s != val {
			res = append(res, s)
		}
	}
	return
}

func Map(ss []string, fn func(string) string) []string {
	for i, s := range ss {
		ss[i] = fn(s)
	}
	return ss
}

func Filter(ss []string, fn func(string) bool) []string {
	res := ss[:0]
	for _, s := range ss {
		if fn(s) {
			res = append(res, s)
		}
	}
	return res
}

func Words(s string) []string {
	return strings.FieldsFunc(s, func(c rune) bool {
		return !unicode.IsLetter(c)
	})
}

func ToTitle(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	return s
}
