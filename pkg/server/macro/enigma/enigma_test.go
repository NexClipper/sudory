package enigma_test

import (
	"bytes"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/macro/enigma"
)

func TestEnigma_10(t *testing.T) {
	//AES 128 CBC PKCS SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = ""
	crypto_alg.CipherMode = "CBC"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 CBC PKCS SALTY
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 CBC PKCS SALTY
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_101(t *testing.T) {
	//AES 128 CBC PKCS SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "CBC"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 CBC PKCS SALTY
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 CBC PKCS SALTY
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_111(t *testing.T) {
	//AES 128 GCM PKCS SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "GCM"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 GCM PKCS SALTY
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 GCM PKCS SALTY
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_121(t *testing.T) {
	//AES 128 NONE PKCS SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 GCM PKCS SALTY
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 GCM PKCS SALTY
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_131(t *testing.T) {
	//AES 128 GCM NONE SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "NONE"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 GCM NONE SALTY
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 GCM NONE SALTY
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_141(t *testing.T) {
	//AES 128 GCM PKCS NULL
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	// crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 GCM NONE NULL
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 GCM NONE NULL
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_151(t *testing.T) {
	//AES 128 GCM NONE NONE
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "NONE"
	// crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//AES 192 GCM NONE SALTY
	crypto_alg.BlockSize = 192
	EnigmaMachine(t, crypto_alg)

	//AES 256 GCM NONE SALTY
	crypto_alg.BlockSize = 256
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_102(t *testing.T) {
	//DES 64 CBC PKCS SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "DES"
	crypto_alg.BlockSize = 64
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "CBC"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_112(t *testing.T) {
	//DES 64 CBC PKCS SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "DES"
	crypto_alg.BlockSize = 64
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_103(t *testing.T) {
	//DES 64 NONE PKCS SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "DES"
	crypto_alg.BlockSize = 64
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

}

func TestEnigma_1(t *testing.T) {
	//AES 128 CBC PKCS SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "cbc"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_12(t *testing.T) {
	//NONE 128 NONE PKCS SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "NONE"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	//NONE 128 CBC PKCS NONE
	crypto_alg.CipherMode = "CBC"
	EnigmaMachine(t, crypto_alg)

	//NONE 128 GCM PKCS NONE
	crypto_alg.CipherMode = "GCM"
	if false {
		EnigmaMachine(t, crypto_alg)
	}
}

func TestEnigma_13(t *testing.T) {
	//AES 128 CBC PKCS SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "CBC"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_14(t *testing.T) {
	//AES 128 NONE PKCS SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "AES"
	crypto_alg.BlockSize = 128
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	crypto_alg.CipherMode = "CBC"
	EnigmaMachine(t, crypto_alg)

	crypto_alg.CipherMode = "GCM"
	EnigmaMachine(t, crypto_alg)
}

func TestEnigma_15(t *testing.T) {
	//DES 128 NONE PKCS SALTY
	var crypto_alg enigma.ConfigCryptoAlgorithm
	crypto_alg.EncryptionMethod = "DES"
	crypto_alg.BlockSize = 64
	crypto_alg.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_alg.CipherMode = "NONE"
	crypto_alg.Padding = "PKCS"
	crypto_alg.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
	crypto_alg.StrConv = "base64"

	EnigmaMachine(t, crypto_alg)

	crypto_alg.CipherMode = "CBC"
	EnigmaMachine(t, crypto_alg)
}

func EnigmaMachine(t *testing.T, alg enigma.ConfigCryptoAlgorithm) {

	crypto, err := enigma.NewMachine(alg.ToOption())
	if err != nil {
		t.Fatal(err)
	}

	s := `
	???????????? ????????????
	????????????
	????????? ??????
	????????? ?????? ????????? ???????????????
	?????? ???????????? ???????????? ????????? ???????????? ?????? ?????? ?????????
	????????? ??? ?????? ?????? ?????? ????????? ????????? ??????.
	?????? ?????? ?????? ????????? ??????
	?????? ???????????? ????????? ?????????
	???????????? ????????? ?????? ?????? ????????? ?????? ????????? ????????? ??? ????????????.`

	var encripttext, plaintext []byte

	var salt_a, salt_b []byte

	encripttext, err = crypto.Encode([]byte(s))
	if err != nil {
		t.Fatal(err)
	}

	plaintext, err = crypto.Decode(encripttext)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(salt_a, salt_b) {
		t.Fatal("diff salt")
	}

	if s != string(plaintext) {
		t.Fatal("diff text", string(plaintext))
	}

}
