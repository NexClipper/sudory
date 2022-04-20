package enigma_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"testing"
)

func TestAesGcm(t *testing.T) {
	text := []byte("My Super Secret Code Stuff")
	// text := []byte("small")
	key := []byte("passphrasewhichneedstobe32bytes!")

	c, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		t.Fatal(err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		t.Fatal(err)
	}

	copy(nonce, key)

	t.Log(nonce)

	encrypt_text := gcm.Seal(nonce, nonce, text, nil)
	t.Log(hex.EncodeToString(encrypt_text))

	nonce, ciphertext := encrypt_text[:gcm.NonceSize()], encrypt_text[gcm.NonceSize():]

	plain_text, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(plain_text))

}

func TestAesCbc(t *testing.T) {
	text := []byte("My Super Secret Code Stuff")
	// text := []byte("small")
	key := []byte("passphrasewhichneedstobe32bytes!")

	c, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	iv := make([]byte, c.BlockSize())
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		t.Fatal(err)
	}

	copy(iv, key)

	t.Log(iv)

	cbc := cipher.NewCBCEncrypter(c, iv)
	if err != nil {
		t.Fatal(err)
	}

	text = PKCS5Padding(text, c.BlockSize(), 0)

	dst := make([]byte, len(text))

	cbc.CryptBlocks(dst, text)

	t.Log(hex.EncodeToString(append(iv, dst...)))

	cbc = cipher.NewCBCDecrypter(c, iv)

	plain_text := make([]byte, len(dst))

	cbc.CryptBlocks(plain_text, dst)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(PKCS5UnPadding(plain_text)))

}

func TestEncode(t *testing.T) {

	text := []byte("My Super Secret Code Stuff")
	// text := []byte("small")
	key := []byte("passphrasewhichneedstobe32bytes!")

	en := Encrypter{
		Block:         SafeAes(key, BlockSize_128),
		StringEncoder: base64.URLEncoding.EncodeToString,
	}

	de := Decrypter{
		Block:         SafeAes(key, BlockSize_128),
		StringDecoder: base64.URLEncoding.DecodeString,
	}

	encrypt_text, err := en.Encrypt(text)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(base64.URLEncoding.EncodeToString(encrypt_text))

	plain_text, err := de.Decrypt(encrypt_text)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(plain_text))
}

type BlockSize int

const (
	BlockSize_128 BlockSize = 128 / 8
	BlockSize_192           = 192 / 8
	BlockSize_256           = 256 / 8
)

type Encrypter struct {
	Block         cipher.Block
	FuncCipher    func(cipher.Block) (interface{}, error)
	StringEncoder func([]byte) string
}

func SafeAes(key []byte, blockSize BlockSize) cipher.Block {
	buf := make([]byte, blockSize)
	copy(buf, key)

	c, _ := aes.NewCipher(buf)

	return c
}

func (encrypter Encrypter) Encrypt(src []byte) (dst []byte, err error) {
	gcm, err := cipher.NewGCM(encrypter.Block)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	dst = gcm.Seal(nonce, nonce, src, nil)

	return
}

func (encrypter Encrypter) EncryptToString(src []byte) (dst string, err error) {
	encrypt, err := encrypter.Encrypt(src)
	if err != nil {
		return "", err
	}
	return encrypter.StringEncoder(encrypt), nil
}

type Decrypter struct {
	Block         cipher.Block
	FuncGcb       func(cipher cipher.Block) (cipher.AEAD, error)
	StringDecoder func(string) ([]byte, error)
}

func (decrypter Decrypter) Decrypt(src []byte) (dst []byte, err error) {
	gcm, err := cipher.NewGCM(decrypter.Block)
	if err != nil {
		return
	}

	dst, err = gcm.Open(nil, src[:gcm.NonceSize()], src[gcm.NonceSize():], nil)

	return
}

func (decrypter Decrypter) DecryptToString(src string) (dst []byte, err error) {
	dst, err = decrypter.StringDecoder(src)
	if err != nil {
		return
	}
	dst, err = decrypter.Decrypt(dst)
	if err != nil {
		return
	}

	return
}

func NewAesCbc(key []byte, blockSize BlockSize) {

}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
