package enigma

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
	"reflect"
	"runtime"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

type Cipher interface {
	EncodeDetail(src []byte, callback ...func(map[string]interface{})) ([]byte, error)
	Encode(src []byte) ([]byte, error)
	DecodeDetail(src []byte, callback ...func(map[string]interface{})) ([]byte, error)
	Decode(src []byte) ([]byte, error)
}

type Encoder func(src, salt []byte) (dst []byte, err error)
type Decoder func(src, salt []byte) (dst []byte, err error)

type Machine struct {
	method  func() EncryptionMethod
	mode    func() CipherMode
	key     func() []byte
	padding func() Padding
	salt    func() *Salt
	block   func() cipher.Block
	Encoder
	Decoder
}

func (machine *Machine) EncodeDetail(src []byte, callback ...func(map[string]interface{})) (dst []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = errors.Wrapf(r, "recovered %s %s encoder",
					machine.method().String(),
					machine.mode().String(),
				)
			default:
				err = fmt.Errorf("recovered %s %s encoder: %v",
					machine.method().String(),
					machine.mode().String(),
					r,
				)
			}
		}
	}()

	//salt
	// salt, hasSalt := machine.salt().GenSalt(), machine.salt().Has()
	err = machine.salt().Scope(func(ss *ScopeSalt) error {
		//padding
		src = machine.padding().Pad(src, machine.block().BlockSize())
		//encode
		dst, err = machine.Encoder(src, ss.GenSalt())
		if err != nil {
			return errors.Wrapf(err, "enigma encode %v",
				logs.KVL(
					"src", src,
					"salt", base64.StdEncoding.EncodeToString(ss.GenSalt()),
				))
		}

		for _, callback := range callback {
			callback(map[string]interface{}{
				"encript":     dst,
				"method":      machine.method().String(),
				"block_size":  machine.block().BlockSize(),
				"block_key":   machine.key(),
				"cipher_mode": machine.mode().String(),
				"cipher_salt": ss.GenSalt(),
				"padding":     machine.padding().String(),
			})
		}

		//encode rule
		dst = EncodeRule(dst, ss.GenSalt(), ss.Has())

		return nil
	})

	return
}

func (machine *Machine) Encode(src []byte) ([]byte, error) {
	return machine.EncodeDetail(src)
}

func (machine *Machine) DecodeDetail(src []byte, callback ...func(map[string]interface{})) (dst []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = errors.Wrapf(r, "recovered %s %s decoder",
					machine.method().String(),
					machine.mode().String(),
				)
			default:
				err = fmt.Errorf("recovered %s %s decoder: %v",
					machine.method().String(),
					machine.mode().String(),
					r,
				)
			}
		}
	}()

	//salt
	err = machine.salt().Scope(func(ss *ScopeSalt) error {
		//decode rule
		src, salt_ := DecodeRule(src, ss.GenSalt(), ss.Has())
		//decode
		dst, err = machine.Decoder(src, salt_)
		if err != nil {
			return errors.Wrapf(err, "enigma decode %v",
				logs.KVL(
					"src", src,
					"salt", base64.StdEncoding.EncodeToString(salt_),
				))
		}

		//unpadding
		dst = machine.padding().Unpad(dst)

		for _, callback := range callback {
			callback(map[string]interface{}{
				"encript":     dst,
				"method":      machine.method().String(),
				"block_size":  machine.block().BlockSize(),
				"block_key":   machine.key(),
				"cipher_mode": machine.mode().String(),
				"cipher_salt": salt_,
				"padding":     machine.padding().String(),
			})
		}

		return nil
	})

	return
}

func (machine *Machine) Decode(src []byte) ([]byte, error) {
	return machine.DecodeDetail(src)
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
				"encryption_method", cfg.EncryptionMethod,
			))
	}

	cipherMode, err := ParseCipherMode(cfg.CipherMode)
	if err != nil {
		return nil, errors.Wrapf(err, "parse cipher mode %v",
			logs.KVL(
				"cipher_mode", cfg.CipherMode,
			))
	}

	padding, err := ParsePadding(cfg.Padding)
	if err != nil {
		return nil, errors.Wrapf(err, "parse cipher padding %v",
			logs.KVL(
				"cipher_padding", cfg.Padding,
			))
	}

	buf, err := base64.StdEncoding.DecodeString(cfg.BlockKey)
	if err != nil {
		return nil, errors.Wrapf(err, "decode key %v",
			logs.KVL(
				"block_key", cfg.BlockKey,
			))
	}
	blockKey := make([]byte, cfg.BlockSize/8)
	copy(blockKey, buf)

	var salt Salt
	if cfg.CipherSalt != nil {
		b, err := base64.StdEncoding.DecodeString(*cfg.CipherSalt)
		if err != nil {
			return nil, errors.Wrapf(err, "decode salt %v",
				logs.KVL(
					"cipher_salt", cfg.CipherSalt,
				))

		}
		salt.SetValue(b)
	}

	newCipher, err := method.BlockFactory()
	if err != nil {
		return nil, errors.Wrapf(err, "block factory %v",
			logs.KVL(
				"encryption_method", method,
			))
	}
	block, err := newCipher(blockKey)
	if err != nil {
		return nil, errors.Wrapf(err, "block factory %v",
			logs.KVL(
				"encryption_method", method,
				"block_key", cfg.BlockKey,
			))
	}

	encoder, decoder, err := cipherMode.CipherFactory(block, &salt)
	if err != nil {
		return nil, errors.Wrapf(err, "cipher builder %v ",
			logs.KVL(
				"encryption_method", cfg.EncryptionMethod,
				"block_key", cfg.BlockKey,
				"cipher_mode", cfg.CipherMode,
			))
	}

	m = &Machine{
		method:  func() EncryptionMethod { return method },
		mode:    func() CipherMode { return cipherMode },
		key:     func() []byte { return blockKey },
		padding: func() Padding { return padding },
		salt:    func() *Salt { return &salt },
		block:   func() cipher.Block { return block },
		Encoder: encoder,
		Decoder: decoder,
	}

	return
}

type FuncSaltMakerEncode func(n int) (salt Salt)
type FuncSaltMakerDecode func(src []byte, n int) (dst []byte, salt Salt)

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

func typeName(i interface{}) string {
	t := reflect.ValueOf(i).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	}
	return t.String()
}
