package ts3

import (
	"strings"
)

func Escape(in string) (out string) {
	out = in
	out = strings.Replace(out, "\\", "\\\\", -1)
	out = strings.Replace(out, " ", "\\s", -1)
	out = strings.Replace(out, "|", "\\p", -1)
	out = strings.Replace(out, "\a", "\\a", -1)
	out = strings.Replace(out, "\b", "\\b", -1)
	out = strings.Replace(out, "\f", "\\f", -1)
	out = strings.Replace(out, "\n", "\\n", -1)
	out = strings.Replace(out, "\r", "\\r", -1)
	out = strings.Replace(out, "\t", "\\t", -1)
	out = strings.Replace(out, "\v", "\\v", -1)
	return
}

func Unescape(in string) (out string) {
	out = in
	out = strings.Replace(out, "\\v", "\v", -1)
	out = strings.Replace(out, "\\t", "\t", -1)
	out = strings.Replace(out, "\\r", "\r", -1)
	out = strings.Replace(out, "\\n", "\n", -1)
	out = strings.Replace(out, "\\f", "\f", -1)
	out = strings.Replace(out, "\\b", "\b", -1)
	out = strings.Replace(out, "\\a", "\a", -1)
	out = strings.Replace(out, "\\p", "|", -1)
	out = strings.Replace(out, "\\s", " ", -1)
	out = strings.Replace(out, "\\\\", "\\", -1)
	return
}
