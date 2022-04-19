package enigma_test

import (
	"testing"

	"github.com/NexClipper/sudory/pkg/server/macro/enigma"
)

func TestEnigmaMachineAesCbc(t *testing.T) {

	var crypto_1 enigma.CryptoAlgorithm
	crypto_1.EncryptionMethod = "aes"
	crypto_1.BlockSize = 128
	crypto_1.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_1.CipherMode = "cbc"
	crypto_1.CipherPadding = NewString("PKCS5")
	crypto_1.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")

	EnigmaMachine(t, crypto_1)

	var crypto_2 enigma.CryptoAlgorithm
	crypto_2.EncryptionMethod = "aes"
	crypto_2.BlockSize = 256
	crypto_2.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_2.CipherMode = "gcm"
	crypto_2.CipherSalt = nil

	EnigmaMachine(t, crypto_2)

	var crypto_3 enigma.CryptoAlgorithm
	crypto_3.EncryptionMethod = "des"
	crypto_3.BlockSize = 64
	crypto_3.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_3.CipherMode = "cbc"
	crypto_3.CipherPadding = NewString("PKCS5")
	crypto_3.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")

	EnigmaMachine(t, crypto_3)

	var crypto_4 enigma.CryptoAlgorithm
	crypto_4.EncryptionMethod = "des"
	crypto_4.BlockSize = 64
	crypto_4.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_4.CipherMode = "gcm"
	crypto_2.CipherSalt = nil

	if false {
		EnigmaMachine(t, crypto_4)
	}
}

func EnigmaMachine(t *testing.T, alg enigma.CryptoAlgorithm) {

	crypto, err := enigma.NewMachine(alg)
	t.Error(err)

	s := `
	세종어제 훈민정음
	나랏말이
	중국과 달라
	문자와 서로 통하지 아니하므로
	이런 까닭으로 어리석은 백성이 이르고자 하는 바가 있어도
	마침내 제 뜻을 능히 펴지 못하는 사람이 많다.
	내가 이를 위해 불쌍히 여겨
	새로 스물여덟 글자를 만드니
	사람마다 하여금 쉬이 익혀 날마다 씀에 편안케 하고자 할 따름이다.`

	var salt, encripttext, plaintext []byte
	if err := crypto.Encode([]byte(s), func(key, salt_, encript_text []byte) {
		t.Log("encode key:", string(key))
		t.Log("encode salt:", string(salt_))
		salt = salt_
		encripttext = encript_text
	}); err != nil {
		t.Error(err)
	}

	if err := crypto.Decode(encripttext, salt, func(key, salt_, plain_text []byte) {
		t.Log("decode key:", string(key))
		t.Log("decode salt:", string(salt_))

		plaintext = plain_text
	}); err != nil {
		t.Error(err)
	}

	t.Log(string(plaintext))

}
