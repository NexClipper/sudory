//go:generate go run github.com/abice/go-enum --file=encryption_method.go --names --nocase
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

type BlockSize_AES int

const (
	BlockSize_AES128 BlockSize_AES = 128 / 8
	BlockSize_AES192               = 192 / 8
	BlockSize_AES256               = 256 / 8
)

func safeAes(key []byte, blockSize BlockSize_AES) cipher.Block {
	b := make([]byte, blockSize)
	copy(b, key)
	c, _ := aes.NewCipher(b)

	return c
}

type BlockSize_DES int

const (
	BlockSize_DES64 BlockSize_DES = 64 / 8
)

func safeDes(key []byte, blockSize BlockSize_DES) cipher.Block {
	b := make([]byte, blockSize)
	copy(b, key)
	c, _ := des.NewCipher(b)

	return c
}
