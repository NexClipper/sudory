package enigma

import (
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/reflected"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

/*
cryptos:
- crypto:
  name: aes-cbc-gen
  cipher: aes
  block-size: 128
  variable:
- crypto:
  name: aes-cbc-salt
  cipher: aes
  block-size: 192
  variable:
  salt:
- crypto:
  name: aes-gcm-gen
  cipher: aes
  block-size: 256
  variable:
- crypto:
  name: aes-gcm-salt
  cipher: aes
  block-size: 128
  variable:
*/

type Config struct {
	CryptoAlgorithms map[string]CryptoAlgorithm `yaml:"enigma"`
}
type CryptoAlgorithm struct {
	Block  `yaml:",inline"`
	Cipher `yaml:",inline"`
}

type Block struct {
	EncryptionMethod string `yaml:"method"`             // aes, des
	BlockSize        int    `yaml:"size" default:"128"` // [128|192|256], [64]
	BlockKey         string `yaml:"key"`                // (base64 string)
}

type Cipher struct {
	CipherMode    string  `yaml:"mode"`    // CBC, GCM
	CipherSalt    *string `yaml:"salt"`    // nil: auto-generate (base64 string)
	CipherPadding *string `yaml:"padding"` // PKCS5,
}

/*
CipherConfig {
	Cipher: "aes"

}
*/

type EncriptMethodAesCbcConfig struct {
	Cipher string `yaml:"cipher"`
}

// deepcopy
//  by yaml package
func deepcopy(dest, src interface{}) error {
	data, err := yaml.Marshal(src)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal yaml%s",
			logs.KVL(
				"src-type-name", reflected.TypeName(src),
				"src", src,
			))
	}

	if err := yaml.Unmarshal(data, dest); err != nil {
		return errors.Wrapf(err, "failed to unmarshal yaml%s",
			logs.KVL(
				"dest-type-name", reflected.TypeName(dest),
				"yaml", data,
			))
	}
	return nil
}
