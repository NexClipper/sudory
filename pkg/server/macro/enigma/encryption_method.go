// go-enum 을 사용해서 열거형 데이터를 만들자
// - go-enum 설치 go-enum.install.sh 파일 실행
// - go generate 실행
//go:generate go-enum --file=encryption_method.go --names --nocase
package enigma

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

/* ENUM(
	NONE
	AES
	DES
)
*/
type EncryptionMethod int

func (method EncryptionMethod) BlockFactory() (fn func(key []byte) (cipher.Block, error), err error) {
	switch method {
	case EncryptionMethodNONE:
		fn = func(key []byte) (cipher.Block, error) { return &NoneEncripter{}, nil }
	case EncryptionMethodAES:
		fn = aes.NewCipher // invalid key size [16,24,32]
	case EncryptionMethodDES:
		fn = des.NewCipher // invalid key size [8]
	default:
		return nil, errors.Errorf("invalid encryption method %v",
			logs.KVL(
				"method", method.String(),
			))
	}

	return
}

type NoneEncripter struct{}

func (encripter NoneEncripter) BlockSize() int {
	return 1
}

func (encripter NoneEncripter) Encrypt(dst, src []byte) {
	copy(dst, src)
}

func (encripter NoneEncripter) Decrypt(dst, src []byte) {
	copy(dst, src)
}
