package enigma

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/NexClipper/sudory/pkg/version"
)

// Config
//  config-name:
//    method: aes     # NONE, AES, DES
//    size: 128       # NONE: default(1), AES: 128|192|256, DES: 64
//    key: secret     # (base64 string)
//    mode: gcm       # NONE: NONE|AES|DES , GCM: AES, CBC: NONE|AES|DES
//    salt: null      # NULL, (base64 string)
//    padding: PKCS   # NONE: AES+NONE default(PKCS)|AES+GCM|DES+NONE default(PKCS)|DES+CBC, PKCS: ALL
//    strconv: base64 # none, base64, hex
type Config struct {
	CryptoAlgorithmSet map[string]ConfigCryptoAlgorithm `yaml:"enigma"`
}
type ConfigCryptoAlgorithm struct {
	ConfigBlock   `yaml:",inline"`
	ConfigCipher  `yaml:",inline"`
	ConfigPadding `yaml:",inline"`
	ConfigStrConv `yaml:",inline"`
}

func (cfg ConfigCryptoAlgorithm) ToOption() MachineOption {
	return configToOption(cfg)
}

func configToOption(cfg ConfigCryptoAlgorithm) (opt MachineOption) {
	opt.Block.Method = cfg.ConfigBlock.EncryptionMethod
	opt.Block.Size = cfg.ConfigBlock.BlockSize
	opt.Block.Key = cfg.ConfigBlock.BlockKey
	opt.Cipher.Mode = cfg.ConfigCipher.CipherMode
	opt.Cipher.Salt = cfg.ConfigCipher.CipherSalt
	opt.Padding = cfg.ConfigPadding.Padding
	opt.StrConv = cfg.ConfigStrConv.StrConv

	return
}

type ConfigBlock struct {
	EncryptionMethod string `yaml:"method"`     // NONE,       AES,           DES
	BlockSize        int    `yaml:"block-size"` // default(1), [128|192|256], [64]
	BlockKey         string `yaml:"block-key"`  // (base64 string)
}

type ConfigCipher struct {
	CipherMode string  `yaml:"cipher-mode"` // NONE, CBC, GCM
	CipherSalt *string `yaml:"cipher-salt"` // nil: auto-generate (base64 string)
}

type ConfigPadding struct {
	Padding string `yaml:"padding"` // none, PKCS
}

type ConfigStrConv struct {
	StrConv string `yaml:"strconv"` // none, base64, hex
}

func PrintConfig(w io.Writer, cfg Config) {
	fmt.Fprintln(w, "enigma configuration:")

	tabwrite := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	insecure := strings.EqualFold(version.Version, "dev")

	col := []string{}
	col = append(col, "")
	col = append(col, "name")
	col = append(col, "encryption-method")
	col = append(col, "block-size")
	if insecure {
		col = append(col, "block-key")
	}
	col = append(col, "cipher-mode")
	if insecure {
		col = append(col, "cipher-salt")
	}
	col = append(col, "padding")
	col = append(col, "strconv")

	tabwrite.Write([]byte(strings.Join(col, "\t") + "\n"))

	for name, cfg := range cfg.CryptoAlgorithmSet {

		row := []string{}
		row = append(row, "-")
		row = append(row, name)
		row = append(row, cfg.EncryptionMethod)
		row = append(row, fmt.Sprintf("%v", cfg.BlockSize))
		if insecure {
			row = append(row, cfg.BlockKey)
		}
		row = append(row, cfg.CipherMode)
		if insecure {
			row = append(row, fmt.Sprintf("%v", cfg.CipherSalt))
		}
		row = append(row, cfg.Padding)
		row = append(row, cfg.StrConv)

		tabwrite.Write([]byte(strings.Join(row, "\t") + "\n"))
	}

	tabwrite.Flush()

	fmt.Fprintln(w, strings.Repeat("_", 40))
}
