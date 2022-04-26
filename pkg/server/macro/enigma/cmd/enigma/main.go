package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/macro/enigma"
	"github.com/jinzhu/configor"
	"gopkg.in/yaml.v2"
)

const default_config_filename = "enigma.yml"

func main() {
	flag.Usage = flagUsageBuilder(func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s [encode|decode|config]\n", procName())
	})

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	for i, arg := range os.Args {
		switch strings.ToLower(arg) {
		case "encode":
			fnCipher(strings.ToLower(arg))
			return
		case "decode":
			fnCipher(strings.ToLower(arg))
			return
		case "config":
			fnConfig(os.Args[i:])
			return
		}
	}

	flag.Usage()
	os.Exit(1)
}

func fnCipher(f string) {
	flag.Usage = flagUsageBuilder(func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s [options] %s [string|stdin]\n", procName(), f)
	})

	var (
		yml string
		cfg enigma.ConfigCryptoAlgorithm
	)
	flag.StringVar(&yml, "yaml", default_config_filename, "*.yml")
	flagEnigmaConfig(&cfg)

	flag.Parse()

	//load config; custom config file name\
	panicking(func(s string) error {
		if exists(s) {
			return configor.Load(&cfg, s)
		}
		return nil
	}(yml))

	do := func(input []byte) {
		var output []byte
		switch strings.ToLower(flag.Arg(0)) {
		case "encode":
			output = right(right(enigma.NewMachine(cfg.ToOption())).(*enigma.Machine).Encode(input)).([]byte)
		case "decode":
			output = right(right(enigma.NewMachine(cfg.ToOption())).(*enigma.Machine).Decode(input)).([]byte)
		default:
			fmt.Fprintln(os.Stderr, "invalid function")
			return
		}

		fmt.Fprintln(os.Stdout, string(output))
	}
	if 0 < len(flag.Arg(1)) {
		do([]byte(flag.Arg(1)))
	} else {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				panicking(err)
			}

			do([]byte(line))
		}
	}
}

func fnConfig(args []string) {
	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "write":
			fnYamlFileWrite()
			return
		case "read":
			fnYamlFileRead()
			return
		}
	}

	fmt.Fprintln(os.Stderr, "invalid function")
}

func fnYamlFileRead() {
	flag.Usage = flagUsageBuilder(func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s [options] config read\n", procName())
	})

	var (
		cfg enigma.ConfigCryptoAlgorithm
		yml string
	)

	flagEnigmaConfig(&cfg)
	flag.StringVar(&yml, "yaml", default_config_filename, "*.yml")

	flag.Parse()

	if exists(yml) {
		output := right(ioutil.ReadFile(yml)).([]byte)
		// fmt.Fprintln(os.Stdout, string(output))

		yaml.Unmarshal(output, &cfg)

		cfg_ := enigma.Config{CryptoAlgorithmSet: map[string]enigma.ConfigCryptoAlgorithm{
			"enigma": cfg,
		}}
		enigma.PrintConfig(os.Stdout, cfg_)

	}
}
func fnYamlFileWrite() {
	flag.Usage = flagUsageBuilder(func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s [options] config write\n", procName())
	})

	var (
		cfg enigma.ConfigCryptoAlgorithm
		yml string
	)

	flagEnigmaConfig(&cfg)
	flag.StringVar(&yml, "yaml", default_config_filename, "*.yml")

	flag.Parse()

	//print config
	cfg_ := enigma.Config{CryptoAlgorithmSet: map[string]enigma.ConfigCryptoAlgorithm{
		"enigma": cfg,
	}}
	enigma.PrintConfig(os.Stdout, cfg_)

	//write to file
	output := right(yaml.Marshal(cfg)).([]byte)
	panicking(ioutil.WriteFile(yml, output, os.ModePerm))

}

func flagEnigmaConfig(cfg *enigma.ConfigCryptoAlgorithm) {

	flag.StringVar(&cfg.EncryptionMethod, "block-method", selectString(cfg.EncryptionMethod, "AES", func(a string) bool { return 0 < len(a) }), `encryption method
	NONE, AES, DES
`)
	flag.IntVar(&cfg.BlockSize, "block-size", selectInt(cfg.BlockSize, 128, func(a int) bool { return 0 < a }), `block size
	NONE: default(1)
	AES: 128|192|256
	DES: 64
`)
	flag.StringVar(&cfg.BlockKey, "block-key", cfg.BlockKey, `block key
	base64 string`)
	flag.StringVar(&cfg.CipherMode, "cipher-mode", selectString(cfg.CipherMode, "GCM", func(a string) bool { return 0 < len(a) }), `cipher mode
	NONE: 
		NONE|AES|DES
	GCM: 
		AES
	CBC: 
		NONE|AES|DES
`)
	flag.Func("cipher-salt", `cipher salt
	base64 string`, func(s string) error {
		if 0 < len(s) {
			cfg.CipherSalt = &s
		}
		return nil
	})
	flag.StringVar(&cfg.Padding, "padding", selectString(cfg.Padding, "NONE", func(a string) bool { return 0 < len(a) }), `padding
	NONE: 
		AES+NONE default(PKCS)
		AES+GCM
		AES+CBC
		DES+NONE default(PKCS)
		DES+CBC 
	PKCS: 
		ALL
`)
	flag.StringVar(&cfg.StrConv, "strconv", selectString(cfg.StrConv, "base64", func(a string) bool { return 0 < len(a) }), `strconv
	none|base64|hex
`)
}

func flagUsageBuilder(fn ...func()) func() {
	return func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", procName())
		for _, fn := range fn {
			fn()
		}
		flag.PrintDefaults()
	}
}

func procName() string {
	return path.Base(strings.ReplaceAll(os.Args[0], "\\", "/"))
}

func panicking(err ...error) {
	for _, err := range err {
		if err != nil {
			panic(err)
		}
	}
}

func right(i interface{}, err error) interface{} {
	panicking(err)

	return i
}

func left(i interface{}, err error) {
	panicking(err)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func selectString(a, b string, filter func(a string) bool) string {
	if filter(a) {
		return a
	}
	return b
}

func selectInt(a, b int, filter func(a int) bool) int {
	if filter(a) {
		return a
	}
	return b
}
