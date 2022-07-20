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

var (
	base64Encoding = base64.StdEncoding
)

func (conv StrConv) Encoder() func([]byte) []byte {
	switch conv {
	case StrConvBase64:
		return func(src []byte) []byte {
			if false {
				dst := base64Encoding.EncodeToString(src)
				return []byte(dst)
			} else {
				dst := make([]byte, base64Encoding.EncodedLen(len(src)))
				base64Encoding.Encode(dst, src)
				return dst
			}
		}
	case StrConvHex:
		return func(src []byte) []byte {
			dst := make([]byte, hex.EncodedLen(len(src)))
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
			if false {
				dst, err := base64Encoding.DecodeString(string(src))
				return dst, err
			} else {
				dst := make([]byte, base64Encoding.DecodedLen(len(src)))
				n, err := base64Encoding.Decode(dst, src)
				return dst[:n], err
			}
		}
	case StrConvHex:
		return func(src []byte) ([]byte, error) {
			// dst := make([]byte, hex.DecodedLen(len(src)))
			// _, err := hex.Decode(dst, src)
			// return dst, err
			n, err := hex.Decode(src, src)
			return src[:n], err
		}
	default:
		return func(src []byte) ([]byte, error) {
			dst := make([]byte, len(src))
			copy(dst, src)
			return dst, nil
		}
	}
}
