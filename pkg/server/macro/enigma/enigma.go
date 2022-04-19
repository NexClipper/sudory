package enigma

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

type Converter interface {
	Encode(src []byte, fn func(key, salt, encripttext []byte)) error
	Decode(src, salt []byte, fn func(key, salt, plaintext []byte)) error
}

type Encoder func(salt []byte, src []byte) (dst []byte, err error)
type Decoder func(salt []byte, src []byte) (dst []byte, err error)

type Machine struct {
	config    CryptoAlgorithm
	key, salt func() []byte
	Encoder
	Decoder
}

func (machine Machine) Encode(src []byte, fn func(key, salt, encripttext []byte)) (err error) {
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

	salt := machine.salt()
	//padding
	switch strings.ToUpper(NullString(machine.config.CipherPadding)) {
	case "PKCS5":
		src = PKCS5Padding(src, len(salt))
	}

	dst, err := machine.Encoder(salt, src)
	if err != nil {
		return errors.Wrapf(err, "enigma encode %v",
			logs.KVL(
				"src", src,
				"salt", base64.StdEncoding.EncodeToString(salt),
			))
	}

	fn(machine.key(), salt, dst)
	return nil
}

func (machine Machine) Decode(src, salt []byte, fn func(key, salt, plaintext []byte)) (err error) {
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

	dst, err := machine.Decoder(salt, src)
	if err != nil {
		return errors.Wrapf(err, "enigma decode %v",
			logs.KVL(
				"src", src,
				"salt", base64.StdEncoding.EncodeToString(salt),
			))
	}

	//padding
	switch strings.ToUpper(NullString(machine.config.CipherPadding)) {
	case "PKCS5":
		dst = PKCS5UnPadding(dst)
	}

	fn(machine.key(), salt, dst)
	return nil
}

var (
	Machines map[string]Converter
)

func init() {
	if Machines == nil {
		Machines = make(map[string]Converter)
	}
}

func LoadConfig(cfg map[string]CryptoAlgorithm) error {
	for k, v := range cfg {
		m, err := NewMachine(v)
		if err != nil {
			return errors.Wrapf(err, "new machine")
		}
		Machines[k] = m
	}

	return nil
}

func GetMachine(k string) Converter {
	return Machines[k]
}

func NewMachine(cfg CryptoAlgorithm) (m *Machine, err error) {
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

	buf, err := base64.StdEncoding.DecodeString(cfg.BlockKey)
	if err != nil {
		return nil, errors.Wrapf(err, "decode key %v",
			logs.KVL(
				"key", cfg.BlockKey,
			))
	}
	key := make([]byte, cfg.BlockSize/8)
	copy(key, buf)

	newCipher, err := BlockFactory(cfg.EncryptionMethod)
	if err != nil {
		return
	}

	block, err := newCipher(key)
	if err != nil {
		return
	}

	e, d, salt, err := CipherFactory(block, cfg.CipherMode, cfg.CipherSalt)
	if err != nil {
		return nil, errors.Errorf("cipher factory %v ",
			logs.KVL(
				"key", base64.StdEncoding.EncodeToString(key),
			))
	}

	m = &Machine{
		config:  cfg,
		key:     func() []byte { return key },
		salt:    salt,
		Encoder: e,
		Decoder: d,
	}

	return
}

func BlockFactory(method string) (fn func(key []byte) (cipher.Block, error), err error) {
	switch strings.ToUpper(method) {
	case "AES":
		fn = aes.NewCipher // invalid key size [16,24,32]
	case "DES":
		fn = des.NewCipher // invalid key size [8]
	default:
		return nil, errors.Errorf("invalid encryption method %v",
			logs.KVL(
				"method", strings.ToUpper(method),
			))
	}

	return
}

func CipherFactory(block cipher.Block, mode string, salt_ *string) (encoder Encoder, decoder Decoder, salt func() []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = errors.Wrapf(r, "recovered")
			default:
				err = fmt.Errorf("recovered %v", r)
			}
		}
	}()

	switch strings.ToUpper(mode) {
	case "CBC":

		salt = func() (b []byte) {
			if salt_ == nil {
				b, err = RandBytes(block.BlockSize())
				if err != nil {
					panic(errors.Wrapf(err, "new salt by randbytes block_size=%d",
						block.BlockSize(),
					))
				}
			} else {
				b, err = NewSalt(*salt_, block.BlockSize())
				if err != nil {
					panic(errors.Wrapf(err, "new salt by base64 string salt='%s' block_size=%d",
						*salt_,
						block.BlockSize(),
					))
				}
			}

			return
		}

		encoder = func(iv, src []byte) (dst []byte, err error) {
			dst = make([]byte, len(src))

			en := cipher.NewCBCEncrypter(block, iv)
			en.CryptBlocks(dst, src)
			return
		}
		decoder = func(iv, src []byte) (dst []byte, err error) {
			dst = make([]byte, len(src))

			de := cipher.NewCBCDecrypter(block, iv)
			de.CryptBlocks(dst, src)
			return
		}

	case "GCM":
		var c cipher.AEAD
		c, err = cipher.NewGCM(block)
		if err != nil {
			err = errors.Wrapf(err, "new gcm cipher %v",
				logs.KVL(
					"type", TypeName(block),
				))

			return
		}
		salt = func() (b []byte) {
			if salt_ == nil {
				b, err = RandBytes(c.NonceSize())
				if err != nil {
					panic(errors.Wrapf(err, "new salt by randbytes nonce_size=%d ",
						c.NonceSize(),
					))
				}
			} else {
				b, err = NewSalt(*salt_, c.NonceSize())
				if err != nil {
					panic(errors.Wrapf(err, "new salt by base64 string salt='%s' nonce_size=%d",
						*salt_,
						c.NonceSize(),
					))
				}
			}
			return
		}

		encoder = func(nonce, src []byte) (dst []byte, err error) {
			dst = c.Seal(nonce, nonce, src, nil)
			dst = dst[len(nonce):]
			return
		}

		decoder = func(nonce, src []byte) (dst []byte, err error) {
			dst, err = c.Open(nil, nonce, src, nil)
			return
		}

	default:
		err = errors.Errorf("invalid cipher mode %v",
			logs.KVL(
				"mode", strings.ToUpper(mode),
			))

		return
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

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func RandBytes(n int) (b []byte, err error) {
	b = make([]byte, n)
	_, err = io.ReadFull(rand.Reader, b)
	return
}

func NewSalt(salt string, n int) (b []byte, err error) {
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

func RemoveBase64Padding(s string, sep ...string) string {
	sep_ := "="
	for _, sep := range sep {
		sep_ = sep
		break
	}

	return strings.ReplaceAll(s, sep_, "")
}

func RecoverBase64Padding(s string, sep ...string) string {
	sep_ := "="
	for _, sep := range sep {
		sep_ = sep
		break
	}

	paddlen := 3 - (len(s) % 3)

	buf := bytes.Buffer{}
	buf.WriteString(s)
	buf.WriteString(strings.Repeat(sep_, paddlen))

	return buf.String()
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
