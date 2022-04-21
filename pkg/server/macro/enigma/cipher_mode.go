// go-enum 을 사용해서 열거형 데이터를 만들자
// - go-enum 설치 go-enum.install.sh 파일 실행
// - go generate 실행
//go:generate go-enum --file=cipher_mode.go --names --nocase
package enigma

import (
	"crypto/cipher"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

/* ENUM(
	NONE
	CBC
	GCM
)
*/
type CipherMode int

func (mode CipherMode) CipherFactory(block cipher.Block, salt *Salt) (encoder Encoder, decoder Decoder, err error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		switch r := r.(type) {
	// 		case error:
	// 			err = errors.Wrapf(r, "recovered")
	// 		default:
	// 			err = fmt.Errorf("recovered %v", r)
	// 		}
	// 	}
	// }()

	switch mode {
	case CipherModeNONE:

		switch block.BlockSize() {
		case 1:
			salt.SetLen(0) //set salt len

			encoder = func(src, salt []byte) (dst []byte, err error) {
				dst = make([]byte, len(src))
				for i := 0; i < len(src); i += block.BlockSize() {
					block.Encrypt(dst[i:i+block.BlockSize()], src[i:i+block.BlockSize()])
				}

				return
			}
			decoder = func(src, salt []byte) (dst []byte, err error) {
				dst = make([]byte, len(src))
				for i := 0; i < len(src); i += block.BlockSize() {
					block.Decrypt(dst[i:i+block.BlockSize()], src[i:i+block.BlockSize()])
				}

				return
			}
		default:
			salt.SetLen(block.BlockSize()) //set salt len

			encoder = func(src, salt []byte) (dst []byte, err error) {
				src = PKCS7Pad(src, block.BlockSize())

				dst = make([]byte, len(src))
				for i := 0; i < len(src); i += block.BlockSize() {
					block.Encrypt(dst[i:i+block.BlockSize()], src[i:i+block.BlockSize()])
				}

				return
			}
			decoder = func(src, salt []byte) (dst []byte, err error) {
				dst = make([]byte, len(src))
				for i := 0; i < len(src); i += block.BlockSize() {
					block.Decrypt(dst[i:i+block.BlockSize()], src[i:i+block.BlockSize()])
				}

				dst = PKCS7Unpad(dst)

				return
			}
		}

	case CipherModeCBC:
		salt.SetLen(block.BlockSize()) //set salt len

		encoder = func(src, iv []byte) (dst []byte, err error) {
			dst = make([]byte, len(src))
			cipher.NewCBCEncrypter(block, iv).CryptBlocks(dst, src)
			return
		}
		decoder = func(src, iv []byte) (dst []byte, err error) {
			dst = make([]byte, len(src))
			cipher.NewCBCDecrypter(block, iv).CryptBlocks(dst, src)
			return
		}

	case CipherModeGCM:
		c, err := cipher.NewGCM(block)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "new gcm cipher %v",
				logs.KVL(
					"type", typeName(block),
				))
		}

		salt.SetLen(c.NonceSize()) //set salt len

		encoder = func(src, nonce []byte) (dst []byte, err error) {
			dst = c.Seal(nonce, nonce, src, nil)
			dst = dst[len(nonce):] //remove nonce
			return
		}
		decoder = func(src, nonce []byte) (dst []byte, err error) {
			dst, err = c.Open(nil, nonce, src, nil)
			return
		}

	default:
		err = errors.Errorf("invalid cipher mode %v",
			logs.KVL(
				"cipher_mode", mode.String(),
			))
	}

	return
}
