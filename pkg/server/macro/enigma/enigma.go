package enigma

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"runtime"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

type Cipher interface {
	Encode(src []byte, callback ...func(key, salt, encripttext []byte)) ([]byte, error)
	Decode(src []byte, callback ...func(key, salt, plaintext []byte)) ([]byte, error)
}

type Encoder func(src []byte) (dst, salt []byte, err error)
type Decoder func(src []byte) (dst, salt []byte, err error)

type Machine struct {
	config             ConfigCryptoAlgorithm
	key                func() []byte
	block              cipher.Block
	FuncSaltRuleEncode //인코드 후 SALT 설정에 따라 인코드된 값에 SALT를 추 유무
	// FuncSaltRuleDecode //인코드 전 SALT 설정에 따라 인코드된 값에서 SALT를 추출 하는지 유무
	Encoder
	Decoder
}

func (machine Machine) Encode(src []byte, callback ...func(key, salt, encripttext []byte)) (dst []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = errors.Wrapf(r, "recovered %s %s encoder",
					machine.config.EncryptionMethod,
					machine.config.CipherMode,
				)
			default:
				err = fmt.Errorf("recovered %s %s encoder: %v",
					machine.config.EncryptionMethod,
					machine.config.CipherMode,
					r,
				)
			}
		}
	}()

	//padding
	padding, _ := ParsePadding(machine.config.CipherPadding)
	src = padding.Pad(src, machine.block.BlockSize())
	// switch strings.ToUpper(machine.config.CipherPadding) {
	// case "PKCS5":
	// 	src = PKCS7Pad(src, len(salt))
	// }

	dst, salt, err := machine.Encoder(src)
	if err != nil {
		return nil, errors.Wrapf(err, "enigma encode %v",
			logs.KVL(
				"src", src,
				"salt", base64.StdEncoding.EncodeToString(salt),
			))
	}

	for _, callback := range callback {
		callback(machine.key(), salt, dst)
	}

	dst = machine.FuncSaltRuleEncode(dst, salt)

	return
}

func (machine Machine) Decode(src []byte, callback ...func(key, salt, plaintext []byte)) (dst []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = errors.Wrapf(r, "recovered %s %s decoder",
					machine.config.EncryptionMethod,
					machine.config.CipherMode,
				)
			default:
				err = fmt.Errorf("recovered %s %s decoder: %v",
					machine.config.EncryptionMethod,
					machine.config.CipherMode,
					r,
				)
			}
		}
	}()

	dst, salt, err := machine.Decoder(src)
	if err != nil {
		return nil, errors.Wrapf(err, "enigma decode %v",
			logs.KVL(
				"src", src,
				"salt", base64.StdEncoding.EncodeToString(salt),
			))
	}

	//padding
	padding, _ := ParsePadding(machine.config.CipherPadding)
	dst = padding.Unpad(dst)
	// switch strings.ToUpper(NullString(machine.config.CipherPadding)) {
	// case "PKCS5":
	// 	dst = PKCS7Unpad(dst)
	// }

	for _, callback := range callback {
		callback(machine.key(), salt, dst)
	}
	return
}

func NewMachine(cfg ConfigCryptoAlgorithm) (m *Machine, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = errors.Wrapf(r, "enigma new machine")
			default:
				err = fmt.Errorf("enigma new machine: %v", r)
			}
		}
	}()

	method, err := ParseEncryptionMethod(cfg.EncryptionMethod)
	if err != nil {
		return nil, errors.Wrapf(err, "parse encryption method %v",
			logs.KVL(
				"encryption-method", cfg.EncryptionMethod,
			))
	}

	newCipher, err := BlockFactory(method)
	if err != nil {
		return nil, errors.Wrapf(err, "block factory %v",
			logs.KVL(
				"encryption-method", method,
			))
	}

	buf, err := base64.StdEncoding.DecodeString(cfg.BlockKey)
	if err != nil {
		return nil, errors.Wrapf(err, "decode key %v",
			logs.KVL(
				"block-key", cfg.BlockKey,
			))
	}
	key := make([]byte, cfg.BlockSize/8)
	copy(key, buf)

	block, err := newCipher(key)
	if err != nil {
		return nil, errors.Wrapf(err, "block factory %v",
			logs.KVL(
				"encryption-method", method,
				"block-key", cfg.BlockKey,
			))
	}

	cipherMode, err := ParseCipherMode(cfg.CipherMode)
	if err != nil {
		return nil, errors.Wrapf(err, "parse cipher mode %v",
			logs.KVL(
				"encryption-method", cfg.EncryptionMethod,
				"cipher-mode", cfg.CipherMode,
			))
	}

	encoder, decoder, err := CipherFactory(
		block,
		cipherMode,
		saltMaker_encode(cfg.CipherSalt),
		saltMaker_decode(cfg.CipherSalt),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "cipher factory %v ",
			logs.KVL(
				"encryption-method", cfg.EncryptionMethod,
				"cipher-mode", cfg.CipherMode,
				"block-key", cfg.BlockKey,
			))
	}

	m = &Machine{
		config:             cfg,
		key:                func() []byte { return key },
		block:              block,
		FuncSaltRuleEncode: saltRule_encode(cfg.CipherSalt),
		// FuncSaltRuleDecode: saltRule_decode(cfg.CipherSalt),
		Encoder: encoder,
		Decoder: decoder,
	}

	return
}

func BlockFactory(method EncryptionMethod) (fn func(key []byte) (cipher.Block, error), err error) {
	switch method {
	case EncryptionMethodNONE:
		// fn = func(key []byte) (cipher.Block, error) { return &NoneEncripter{key: key}, nil }
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

type FuncSaltMakerEncode func(n int) (dst []byte)
type FuncSaltMakerDecode func(src []byte, n int) (dst, salt []byte)

type FuncSaltRuleEncode func(src, salt []byte) (dst []byte)
type FuncSaltRuleDecode func(src, salt []byte) (dst, salt_ []byte)

func CipherFactory(block cipher.Block, mode CipherMode,
	saltMaker_encode FuncSaltMakerEncode, saltMaker_decode FuncSaltMakerDecode) (encoder Encoder, decoder Decoder /*, salt func() []byte*/, err error) {
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
		encoder = func(src []byte) (dst, iv []byte, err error) {
			src = PKCS7Pad(src, block.BlockSize())

			dst = make([]byte, len(src))

			var i int = 0
			for i = 0; i < len(src); i += block.BlockSize() {
				block.Encrypt(dst[i:i+block.BlockSize()], src[i:i+block.BlockSize()])
			}

			return
		}
		decoder = func(src []byte) (dst, iv []byte, err error) {
			dst = make([]byte, len(src))

			for i := 0; i < len(src); i += block.BlockSize() {
				block.Decrypt(dst[i:i+block.BlockSize()], src[i:i+block.BlockSize()])
			}

			dst = PKCS7Unpad(dst)
			return
		}
	case CipherModeCBC:
		encoder = func(src []byte) (dst, iv []byte, err error) {
			iv = saltMaker_encode(block.BlockSize())
			en := cipher.NewCBCEncrypter(block, iv)

			dst = make([]byte, len(src))
			en.CryptBlocks(dst, src)
			return
		}
		decoder = func(src []byte) (dst, iv []byte, err error) {
			src, iv = saltMaker_decode(src, block.BlockSize())

			de := cipher.NewCBCDecrypter(block, iv)

			dst = make([]byte, len(src))
			de.CryptBlocks(dst, src)
			return
		}

	case CipherModeGCM:
		var c cipher.AEAD
		c, err = cipher.NewGCM(block)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "new gcm cipher %v",
				logs.KVL(
					"type", TypeName(block),
				))
		}

		encoder = func(src []byte) (dst, nonce []byte, err error) {
			nonce = saltMaker_encode(c.NonceSize()) //make nonce

			dst = c.Seal(nonce, nonce, src, nil)
			dst = dst[len(nonce):]
			return
		}
		decoder = func(src []byte) (dst, nonce []byte, err error) {
			src, nonce = saltMaker_decode(src, c.NonceSize()) //make nonce

			dst, err = c.Open(nil, nonce, src, nil)
			return
		}

	default:
		return nil, nil, errors.Errorf("invalid cipher mode %v",
			logs.KVL(
				"cipher-mode", mode.String(),
			))
	}

	return
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

func RandBytes(n int) (b []byte, err error) {
	b = make([]byte, n)
	_, err = io.ReadFull(rand.Reader, b)
	return
}

func safeSalt(salt string, n int) (b []byte, err error) {
	salt_, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, errors.Wrapf(err, "decode salt %v",
			logs.KVL(
				"salt", salt,
			))

	}

	b = make([]byte, n)
	copy(b, salt_)

	return
}

func getSaltSize(cipherMode interface{}) (n int, err error) {
	switch cipherMode := cipherMode.(type) {
	case cipher.AEAD:
		n = cipherMode.NonceSize()
	case cipher.BlockMode:
		n = cipherMode.BlockSize()
	default:
		err = errors.Errorf("invalid cipher-mode")
	}

	return
}

func saltRule_encode(salt *string) FuncSaltRuleEncode {
	if salt == nil {
		return func(src, salt []byte) (dst []byte) {
			dst = make([]byte, len(src)+len(salt))

			copy(dst, append(salt, src...))
			return
		}
	} else {
		return func(src, salt []byte) (dst []byte) {
			dst = make([]byte, len(src))

			copy(dst, src)
			return
		}
	}

}

func saltMaker_encode(salt *string) FuncSaltMakerEncode {
	var err error
	if salt == nil {
		return func(n int) (salt_ []byte) {
			salt_, err = RandBytes(n)
			if err != nil {
				panic(errors.Wrapf(err, "new salt by randbytes %v",
					logs.KVL(
						"salt_size", n,
					)))
			}
			return
		}
	} else {
		return func(n int) (salt_ []byte) {
			salt_, err = safeSalt(NullString(salt), n)
			if err != nil {
				panic(errors.Wrapf(err, "new salt by base64 string %v",
					logs.KVL(
						"salt", NullString(salt),
						"salt_size", n,
					)))
			}
			return
		}
	}
}

func saltMaker_decode(salt *string) FuncSaltMakerDecode {
	if salt == nil {
		return func(src []byte, n int) (dst, salt_ []byte) {
			dst, salt_ = src[n:], src[:n]

			return
		}
	} else {
		var err error
		return func(src []byte, n int) (dst, salt_ []byte) {
			dst = make([]byte, len(src))
			copy(dst, src)

			salt_, err = safeSalt(NullString(salt), n)
			if err != nil {
				panic(errors.Wrapf(err, "new salt by base64 string %v",
					logs.KVL(
						"salt", NullString(salt),
						"salt_size", n,
					)))
			}

			return
		}
	}
}

func TypeName(i interface{}) string {
	t := reflect.ValueOf(i).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	}
	return t.String()
}

func NullString(s *string) (r string) {
	r = ""
	if s != nil {
		r = *s
	}

	return
}

// func reflectValueOfPointer(i interface{}) uintptr {
// 	return reflect.ValueOf(i).Pointer()
// }
