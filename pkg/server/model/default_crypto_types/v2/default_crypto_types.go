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
	var i sql.NullString
	var b []byte
	switch value := value.(type) {
	case string:
		if err := i.Scan(value); err != nil {
			return err
		}
		b = []byte(i.String)
	case []byte:
		b = value
	}

	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*cs = (CryptoString)(i.String)
	} else {
		bytes, err := EnigmaDecode(b)
		if err != nil {
			return errors.Wrapf(err, "default crypto string: decode")
		}
		*cs = CryptoString(bytes)
	}
	return nil
}
func (cs CryptoString) Value() (driver.Value, error) {
	out, err := EnigmaEncode([]byte(cs))
	if err != nil {
		return out, errors.Wrapf(err, "default crypto string: encode")
	}

	return out, nil
}

// CryptoJson
type CryptoJson map[string]interface{}

func (cj CryptoJson) Json() map[string]interface{} {
	return (map[string]interface{})(cj)
}

func (cj *CryptoJson) Scan(value interface{}) error {

	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*cj = map[string]interface{}{}
	} else {
		bytes, err := EnigmaDecode([]byte(i.String))
		if err != nil {
			return errors.Wrapf(err, "default crypto hashset: decode")
		}

		if err := json.Unmarshal(bytes, cj); err != nil {
			return errors.Wrapf(err, "default crypto hashset: json unmarshal")
		}
	}
	return nil
}
func (cj CryptoJson) Value() (driver.Value, error) {
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
