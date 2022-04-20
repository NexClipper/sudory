package enigma

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/NexClipper/sudory/pkg/version"
	"github.com/pkg/errors"
)

// Config
//  config-name:
//    method: aes   # NONE, AES, DES
//    size: 128     # NONE: default(1), AES: 128|192|256, DES: 64
//    key: secret   # (base64 string)
//    mode: gcm     # NONE: NONE|AES|DES , GCM: AES, CBC: NONE|AES|DES
//    salt: null    # NULL, (base64 string)
//    padding: PKCS # NONE: AES+NONE default(PKCS)|AES+GCM|DES+NONE default(PKCS)|DES+CBC, PKCS: ALL
type Config struct {
	CryptoAlgorithmSet map[string]ConfigCryptoAlgorithm `yaml:"enigma"`
}
type ConfigCryptoAlgorithm struct {
	ConfigBlock  `yaml:",inline"`
	ConfigCipher `yaml:",inline"`
}

type ConfigBlock struct {
	EncryptionMethod string `yaml:"method"`             // NONE, AES, DES
	BlockSize        int    `yaml:"size" default:"128"` // default(1), [128|192|256], [64]
	BlockKey         string `yaml:"key"`                // (base64 string)
}

type ConfigCipher struct {
	CipherMode    string  `yaml:"mode"`    // NONE, CBC, GCM
	CipherSalt    *string `yaml:"salt"`    // nil: auto-generate (base64 string)
	CipherPadding string  `yaml:"padding"` // [none|PKCS], [PKCS], [none|PKCS]
}

type EncriptMethodAesCbcConfig struct {
	Cipher string `yaml:"cipher"`
}

var (
	Machines map[string]Cipher
)

func init() {
	if Machines == nil {
		Machines = make(map[string]Cipher)
	}
}

func LoadConfig(cfg map[string]ConfigCryptoAlgorithm) error {
	for k, v := range cfg {
		machine, err := NewMachine(v)
		if err != nil {
			return errors.Wrapf(err, "new machine")
		}
		Machines[k] = machine
	}

	return nil
}

func GetMachine(k string) Cipher {
	return Machines[k]
}

func PrintConfig(w io.Writer, cfg map[string]ConfigCryptoAlgorithm) {
	fmt.Fprintln(w, "enigma configuration:")

	tabwrite := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	if strings.Compare(version.Version, "dev") == 0 {
		tabwrite.Write([]byte(strings.Join([]string{
			"",
			"name",
			"encryption-method",
			"block-size",
			"block-key",
			"cipher-mode",
			"cipher-salt",
			"cipher-padding",
		}, "\t") + "\n"))

		for name, cfg := range cfg {
			tabwrite.Write([]byte(strings.Join([]string{
				"-",
				name,
				cfg.EncryptionMethod,
				fmt.Sprintf("%v", cfg.BlockSize),
				cfg.BlockKey,
				cfg.CipherMode,
				fmt.Sprintf("%v", cfg.CipherSalt),
				cfg.CipherPadding,
			}, "\t") + "\n"))
		}
	} else {

		tabwrite.Write([]byte(strings.Join([]string{
			"",
			"name",
			"encryption-method",
			"block-size",
			// "block-key",
			"cipher-mode",
			// "cipher-salt",
			"cipher-padding",
		}, "\t") + "\n"))

		for name, cfg := range cfg {
			tabwrite.Write([]byte(strings.Join([]string{
				"-",
				name,
				cfg.EncryptionMethod,
				fmt.Sprintf("%v", cfg.BlockSize),
				// cfg.BlockKey,
				cfg.CipherMode,
				// fmt.Sprintf("%v", cfg.CipherSalt),
				cfg.CipherPadding,
			}, "\t") + "\n"))
		}
	}
	tabwrite.Flush()

	fmt.Fprintln(w, strings.Repeat("_", 40))
}
