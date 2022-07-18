//go:generate go run github.com/abice/go-enum --file=padding.go --names --nocase
package enigma

import (
	"bytes"
)

/* ENUM(
NONE
PKCS
)
*/
type Padding int

func (padding Padding) Padder() func([]byte, int) []byte {
	switch padding {
	case PaddingPKCS:
		return func(src []byte, blockSize int) []byte {
			return PKCS7Padding(src, blockSize)
		}
	default:
		return func(src []byte, blockSize int) []byte {
			dst := make([]byte, len(src))
			copy(dst, src)
			return dst
		}

	}

}

func (padding Padding) Unpadder() func(src []byte) (dst []byte) {
	switch padding {
	case PaddingPKCS:
		return func(src []byte) []byte {
			return PKCS7Unpadding(src)
		}
	default:
		return func(src []byte) []byte {
			dst := make([]byte, len(src))
			copy(dst, src)
			return dst
		}
	}
}

func PKCS7Padding(src []byte, blockSize int) []byte {
	padLen := blockSize - len(src)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padding...)
}

func PKCS7Unpadding(src []byte) []byte {
	length := len(src)
	padLen := int(src[length-1])
	return src[:(length - padLen)]
}
