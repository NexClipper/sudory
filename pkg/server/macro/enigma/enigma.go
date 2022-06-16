package enigma

import (
	"github.com/pkg/errors"
)

var (
	_CipherSet map[string]Cipher
)

func init() {
	if _CipherSet == nil {
		_CipherSet = make(map[string]Cipher)
	}
}

func LoadConfig(cfg Config) error {
	for k, v := range cfg.CryptoAlgorithmSet {
		machine, err := NewMachine(v.ToOption())
		if err != nil {
			return errors.Wrapf(err, "new machine")
		}
		_CipherSet[k] = machine
	}

	return nil
}

func CipherSet(k string) Cipher {
	return _CipherSet[k]
}
