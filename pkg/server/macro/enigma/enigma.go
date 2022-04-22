package enigma

import (
	"github.com/pkg/errors"
)

var (
	Machines map[string]Cipher
)

func init() {
	if Machines == nil {
		Machines = make(map[string]Cipher)
	}
}

func LoadConfig(cfg Config) error {
	for k, v := range cfg.CryptoAlgorithmSet {
		machine, err := NewMachine(v.ToOption())
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
