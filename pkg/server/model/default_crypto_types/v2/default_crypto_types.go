package v2

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"

	"github.com/NexClipper/sudory/pkg/server/macro/enigma"
	"github.com/pkg/errors"
)

const DefaultCryptoName = "sudory.default.crypto"

func EnigmaEncode(bytes []byte) (out []byte, err error) {
	if out, err = enigma.CipherSet(CiperKeySudoryDefaultCrypto.String()).Encode(bytes); err != nil {
		return nil, errors.Wrapf(err, "enigma encode")
	}
	return
}

func EnigmaDecode(bytes []byte) (out []byte, err error) {
	if out, err = enigma.CipherSet(CiperKeySudoryDefaultCrypto.String()).Decode(bytes); err != nil {
		return out, errors.Wrapf(err, "enigma decode")
	}
	return
}

// CryptoString
type CryptoString string

func (cs CryptoString) String() string {
	return (string)(cs)
}

func (cs *CryptoString) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		return nil
	}

	var b []byte
	switch value := value.(type) {
	case string:
		var i sql.NullString
		if err := i.Scan(value); err != nil {
			return err
		}
		b = []byte(i.String)
	case []byte:
		b = value
	default:
		return errors.New("invalid type")
	}

	bytes, err := EnigmaDecode(b)
	if err != nil {
		return errors.Wrapf(err, "default crypto string: decode")
	}
	*cs = CryptoString(bytes)

	return nil
}
func (cs CryptoString) Value() (driver.Value, error) {
	out, err := EnigmaEncode([]byte(cs))
	if err != nil {
		return out, errors.Wrapf(err, "default crypto string: encode")
	}

	return out, nil
}

// CryptoObject
type CryptoObject map[string]interface{}

func (cj CryptoObject) Object() map[string]interface{} {
	return (map[string]interface{})(cj)
}

func (cj *CryptoObject) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		return nil
	}

	var b []byte
	switch value := value.(type) {
	case string:
		var i sql.NullString
		if err := i.Scan(value); err != nil {
			return err
		}
		b = []byte(i.String)
	case []byte:
		b = value
	default:
		return errors.New("invalid type")
	}

	bytes, err := EnigmaDecode(b)
	if err != nil {
		return errors.Wrapf(err, "default crypto hashset: decode")
	}

	if err := json.Unmarshal(bytes, cj); err != nil {
		return errors.Wrapf(err, "default crypto hashset: json unmarshal")
	}

	return nil
}
func (cj CryptoObject) Value() (driver.Value, error) {
	out, err := json.Marshal(cj)
	if err != nil {
		return string(out), errors.Wrapf(err, "default crypto hashset: json marshal")
	}
	out, err = EnigmaEncode(out)
	if err != nil {
		return string(out), errors.Wrapf(err, "default crypto hashset: encode")
	}

	return string(out), nil
}
