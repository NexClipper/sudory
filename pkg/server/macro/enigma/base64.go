package enigma

import (
	"bytes"
	"strings"
)

func RemoveBase64Padding(s string, sep ...string) string {
	sep_ := "="
	for _, sep := range sep {
		sep_ = sep
		break
	}

	return strings.ReplaceAll(s, sep_, "")
}

func RecoverBase64Padding(s string, sep ...string) string {
	sep_ := "="
	for _, sep := range sep {
		sep_ = sep
		break
	}

	paddlen := 3 - (len(s) % 3)

	buf := bytes.Buffer{}
	buf.WriteString(s)
	buf.WriteString(strings.Repeat(sep_, paddlen))

	return buf.String()
}
