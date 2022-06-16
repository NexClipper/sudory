package v1

import (
	"encoding/json"

	"github.com/NexClipper/sudory/pkg/server/macro/enigma"
	"github.com/pkg/errors"
)

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

// String
type String string

func (field *String) FromDB(bytes []byte) (err error) {
	if bytes, err = EnigmaDecode(bytes); err != nil {
		return errors.Wrapf(err, "default crypto string: decode")
	}

	*field = String(bytes)

	return
}
func (field String) ToDB() (out []byte, err error) {
	if out, err = EnigmaEncode([]byte(field)); err != nil {
		return out, errors.Wrapf(err, "default crypto string: encode")
	}

	return
}

// Hashset
type Hashset map[string]interface{}

func (field *Hashset) FromDB(bytes []byte) (err error) {
	if bytes, err = EnigmaDecode(bytes); err != nil {
		return errors.Wrapf(err, "default crypto hashset: decode")
	}

	if err = json.Unmarshal(bytes, field); err != nil {
		return errors.Wrapf(err, "default crypto hashset: json unmarshal")
	}

	return
}
func (field Hashset) ToDB() (out []byte, err error) {
	if out, err = json.Marshal(field); err != nil {
		return nil, errors.Wrapf(err, "default crypto hashset: json marshal")
	}

	if out, err = EnigmaEncode(out); err != nil {
		return out, errors.Wrapf(err, "default crypto hashset: encode")
	}

	return
}
