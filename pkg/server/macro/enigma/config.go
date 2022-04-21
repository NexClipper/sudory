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
	ConfigBlock   `yaml:",inline"`
	ConfigCipher  `yaml:",inline"`
	ConfigPadding `yaml:",inline"`
}

type ConfigBlock struct {
	EncryptionMethod string `yaml:"method"`                   // NONE,       AES,           DES
	BlockSize        int    `yaml:"block-size" default:"128"` // default(1), [128|192|256], [64]
	BlockKey         string `yaml:"block-key"`                // (base64 string)
}

type ConfigCipher struct {
	CipherMode string  `yaml:"cipher-mode"` // NONE, CBC, GCM
	CipherSalt *string `yaml:"cipher-salt"` // nil: auto-generate (base64 string)
}

type ConfigPadding struct {
	Padding string `yaml:"padding"` // none, PKCS
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
			"padding",
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
				cfg.Padding,
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
			"padding",
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
				cfg.Padding,
			}, "\t") + "\n"))
		}
	}
	tabwrite.Flush()

	fmt.Fprintln(w, strings.Repeat("_", 40))
}
