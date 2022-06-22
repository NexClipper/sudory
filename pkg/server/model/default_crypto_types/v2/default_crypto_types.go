package v2

import (
	"encoding/json"

	"github.com/NexClipper/sudory/pkg/server/macro/enigma"
	v2 "github.com/NexClipper/sudory/pkg/server/model/noxorm/v2"
	"github.com/pkg/errors"
)

const DefaultCryptoName = "sudory.default.crypto"

func EnigmaEncode(bytes []byte) (out []byte, err error) {
	if out, err = enigma.GetMachine(DefaultCryptoName).Encode(bytes); err != nil {
		return nil, errors.Wrapf(err, "enigma encode")
	}
	return
}

func EnigmaDecode(bytes []byte) (out []byte, err error) {
	if out, err = enigma.GetMachine(DefaultCryptoName).Decode(bytes); err != nil {
		return out, errors.Wrapf(err, "enigma decode")
	}
	return
}

// CryptoString
type CryptoString v2.NullString

func (cs CryptoString) String() string {
	return (string)(cs)
}

func (field *CryptoString) FromDB(bytes []byte) (err error) {
	if bytes, err = EnigmaDecode(bytes); err != nil {
		return errors.Wrapf(err, "default crypto string: decode")
	}

	*field = CryptoString(bytes)

	return
}
func (field CryptoString) ToDB() (out []byte, err error) {
	if out, err = EnigmaEncode([]byte(field)); err != nil {
		return out, errors.Wrapf(err, "default crypto string: encode")
	}

	return
}

// CryptoJson
type CryptoJson v2.NullJson

func (cj CryptoJson) Json() map[string]interface{} {
	return (map[string]interface{})(cj)
}

func (field *CryptoJson) FromDB(bytes []byte) (err error) {
	if bytes, err = EnigmaDecode(bytes); err != nil {
		return errors.Wrapf(err, "default crypto hashset: decode")
	}

	if err = json.Unmarshal(bytes, field); err != nil {
		return errors.Wrapf(err, "default crypto hashset: json unmarshal")
	}

	return
}
func (field CryptoJson) ToDB() (out []byte, err error) {
	if out, err = json.Marshal(field); err != nil {
		return nil, errors.Wrapf(err, "default crypto hashset: json marshal")
	}

	if out, err = EnigmaEncode(out); err != nil {
		return out, errors.Wrapf(err, "default crypto hashset: encode")
	}

	return
}
