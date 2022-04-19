package enigma_test

import (
	"fmt"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/macro/enigma"
	"gopkg.in/yaml.v2"
)

func TestLoadCryptoConfig(t *testing.T) {

	const s = `
cryptos:
- crypto 1:
  encryption: AES
  size:       128
  key:        'YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n' 
  mode:       CBC
  salt:       null
- crypto 2:
  encryption: AES
  size:       256
  key:        'YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n' 
  mode:       GCM
  salt:       null
- crypto 3:
  encryption: DES
  size:       64
  key:        'YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n' 
  mode:       CBC
  salt:       null
- crypto 4:
  encryption: DES
  size:       64
  key:        'YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n' 
  mode:       CBC
  salt:       null
`

	var cfg map[string]enigma.CryptoAlgorithm

	if err := yaml.Unmarshal([]byte(s), &cfg); err != nil {
		t.Error(err)
	}

}

func TestMashalCryptoConfig(t *testing.T) {

	const s = `
cryptos:
- crypto 1:
  encryption: AES
  size:       128
  key:        'YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n' 
  mode:       CBC
  salt:       null
- crypto 2:
  encryption: AES
  size:       256
  key:        'YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n' 
  mode:       GCM
  salt:       null
- crypto 3:
  encryption: DES
  size:       64
  key:        'YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n' 
  mode:       CBC
  salt:       null
- crypto 4:
  encryption: DES
  size:       64
  key:        'YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n' 
  mode:       CBC
  salt:       null
`

	var crypto_1 enigma.CryptoAlgorithm
	crypto_1.EncryptionMethod = "aes"
	crypto_1.BlockSize = 128
	crypto_1.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_1.CipherMode = "cbc"
	crypto_1.CipherSalt = nil

	var crypto_2 enigma.CryptoAlgorithm
	crypto_2.EncryptionMethod = "aes"
	crypto_2.BlockSize = 128
	crypto_2.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_2.CipherMode = "gcm"
	crypto_2.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")

	var crypto_3 enigma.CryptoAlgorithm
	crypto_3.EncryptionMethod = "des"
	crypto_3.BlockSize = 64
	crypto_3.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_3.CipherMode = "cbc"
	crypto_3.CipherSalt = nil

	var crypto_4 enigma.CryptoAlgorithm
	crypto_4.EncryptionMethod = "des"
	crypto_4.BlockSize = 64
	crypto_4.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
	crypto_4.CipherMode = "gcm"
	crypto_4.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")

	var cfg = map[string]enigma.CryptoAlgorithm{}

	cfg["crypto_1"] = crypto_1
	cfg["crypto_2"] = crypto_2
	cfg["crypto_3"] = crypto_3
	cfg["crypto_4"] = crypto_4

	b, err := yaml.Marshal(cfg)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(b))
	fmt.Println(string(b))
}

func NewString(s string) *string { return &s }
