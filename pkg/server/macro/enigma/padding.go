// go-enum 을 사용해서 열거형 데이터를 만들자
// - go-enum 설치 go-enum.install.sh 파일 실행
// - go generate 실행
//go:generate go-enum --file=padding.go --names --nocase
package enigma

import (
	"bytes"
)

/* ENUM (
	NONE
	PKCS
)
*/
type Padding int

func (padding Padding) Pader() func([]byte, int) []byte {
	switch padding {
	case PaddingPKCS:
		return func(src []byte, blockSize int) []byte {
			return PKCS7Pad(src, blockSize)
		}
	default:
		return func(src []byte, blockSize int) []byte {
			dst := make([]byte, len(src))
			copy(dst, src)
			return dst
		}

	}

}

func (padding Padding) Unpader() func(src []byte) (dst []byte) {
	switch padding {
	case PaddingPKCS:
		return func(src []byte) []byte {
			return PKCS7Unpad(src)
		}
	default:
		return func(src []byte) []byte {
			dst := make([]byte, len(src))
			copy(dst, src)
			return dst
		}
	}
}

func PKCS7Pad(src []byte, blockSize int) []byte {
	padLen := blockSize - len(src)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padding...)
}

func PKCS7Unpad(src []byte) []byte {
	length := len(src)
	padLen := int(src[length-1])
	return src[:(length - padLen)]
}
