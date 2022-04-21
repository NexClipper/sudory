package v1

import (
	"encoding/base64"
	"encoding/json"

	"github.com/NexClipper/sudory/pkg/server/macro/enigma"
	"github.com/pkg/errors"
)

const sudory_default_crypto = "sudory.default.crypto"

func toDB(bytes []byte) (out []byte, err error) {
	var buf []byte

	if buf, err = enigma.GetMachine(sudory_default_crypto).Encode(bytes); err != nil {
		return nil, errors.Wrapf(err, "encode")
	}

	out = []byte(base64.StdEncoding.EncodeToString(buf))

	return
}

func fromDB(bytes []byte) (out []byte, err error) {
	var buf []byte
	if buf, err = base64.StdEncoding.DecodeString(string(bytes)); err != nil {
		return out, errors.Wrapf(err, "base 64 decode")
	}

	if out, err = enigma.GetMachine(sudory_default_crypto).Decode(buf); err != nil {
		return out, errors.Wrapf(err, "decode ")
	}

	return
}

// String
type String string

func (field *String) FromDB(bytes []byte) (err error) {
	if bytes, err = fromDB(bytes); err != nil {
		return errors.Wrapf(err, "default enigma string: fromDB")
	}

	*field = String(bytes)

	return
}
func (field String) ToDB() (out []byte, err error) {
	if out, err = toDB([]byte(field)); err != nil {
		return out, errors.Wrapf(err, "default enigma string: toDB")
	}

	return
}

// Hashset
type Hashset map[string]interface{}

func (field *Hashset) FromDB(bytes []byte) (err error) {
	if bytes, err = fromDB(bytes); err != nil {
		return errors.Wrapf(err, "default enigma hashset: fromDB")
	}

	if err = json.Unmarshal(bytes, field); err != nil {
		return errors.Wrapf(err, "default enigma hashset: json unmarshal")
	}

	return
}
func (field Hashset) ToDB() (out []byte, err error) {
	if out, err = json.Marshal(field); err != nil {
		return nil, errors.Wrapf(err, "default enigma hashset: json marshal")
	}

	if out, err = toDB(out); err != nil {
		return out, errors.Wrapf(err, "default enigma hashset: toDB")
	}

	return
}
