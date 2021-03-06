//go:generate go run github.com/abice/go-enum --file=strconv.go --names --nocase
package enigma

import (
	"encoding/base64"
	"encoding/hex"
)

/* ENUM(
plain
base64
hex
)
*/
type StrConv int

func (conv StrConv) Encoder() func([]byte) string {
	switch conv {
	case StrConvBase64:
		return func(b []byte) string {
			return base64.StdEncoding.EncodeToString(b)
		}
	case StrConvHex:
		return func(b []byte) string {
			return hex.EncodeToString(b)
		}
	default:
		return func(b []byte) string {
			dst := make([]byte, len(b))
			copy(dst, b)
			return string(dst)
		}
	}
}

func (conv StrConv) Decoder() func(string) ([]byte, error) {
	switch conv {
	case StrConvBase64:
		return func(s string) ([]byte, error) {
			return base64.StdEncoding.DecodeString(s)
		}
	case StrConvHex:
		return func(s string) ([]byte, error) {
			return hex.DecodeString(s)
		}
	default:
		return func(s string) ([]byte, error) {
			src := []byte(s)
			dst := make([]byte, len(src))
			copy(dst, src)
			return dst, nil
		}
	}
}
