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

func (conv StrConv) Encoder() func([]byte) []byte {
	switch conv {
	case StrConvBase64:
		return func(src []byte) []byte {
			var dst []byte
			base64.StdEncoding.Encode(dst, src)
			return dst
		}
	case StrConvHex:
		return func(src []byte) []byte {
			var dst []byte
			hex.Encode(dst, src)
			return dst
		}
	default:
		return func(src []byte) []byte {
			dst := make([]byte, len(src))
			copy(dst, src)
			return dst
		}
	}
}

func (conv StrConv) Decoder() func([]byte) ([]byte, error) {
	switch conv {
	case StrConvBase64:
		return func(src []byte) ([]byte, error) {
			var dst []byte
			_, err := base64.StdEncoding.Decode(dst, src)
			return dst, err
		}
	case StrConvHex:
		return func(src []byte) ([]byte, error) {
			var dst []byte
			_, err := hex.Decode(dst, src)
			return dst, err
		}
	default:
		return func(src []byte) ([]byte, error) {
			dst := make([]byte, len(src))
			copy(dst, src)
			return dst, nil
		}
	}
}
