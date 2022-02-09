package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

const (
	jwt_header    = 0
	jwt_payload   = 1
	jwt_signature = 2
)

func ErrorJwtInvalidLength() error {
	return errors.New("jwt have 3 parts")
}
func ErrorJwtInvalidSignature() error {
	return errors.New("invalid signature")
}

func ErrorJwtBindPayloadNonePointer() error {
	return errors.New("jwt bind payload none pointer")
}
func ErrorJwtBindPayloadNil() error {
	return errors.New("jwt bind payload nil")
}

// JWT new
//  make new JWT
func New(payload interface{}, secret []byte) (jwt string, err error) {
	defer func() {
		var ok bool
		if r := recover(); r != nil {
			err, ok = r.(error)
			if !ok {
				panic(r)
			}
		}
	}()

	//define header
	var head = map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	json_mashal := func(v interface{}) []byte {
		// json_mashal := func(v interface{}) ([]byte, error) { return json.MarshalIndent(v, "", " ") }
		json_mashal := json.Marshal
		right := func(b []byte, err error) []byte {
			if err != nil {
				panic(err)
			}
			return b
		}
		return right(json_mashal(v))
	}

	header := byte_encoder(json_mashal(head))                               //make header
	payload_ := byte_encoder(json_mashal(payload))                          //make payload
	signature := byte_encoder(HMACSHA256(header, payload_, []byte(secret))) //make signature
	jwt = strings.Join([]string{header, payload_, signature}, ".")          //make jwt

	return jwt, err
}

// JWT Verify
//  valied signature
func Verify(jwt string, secret []byte) error {

	parts := strings.Split(jwt, ".") //split diffrent parts

	if len(parts) != 3 {
		return ErrorJwtInvalidLength()
	}

	signature := byte_encoder(HMACSHA256(parts[jwt_header], parts[jwt_payload], []byte(secret))) //make signature

	if parts[jwt_signature] != signature { //compare
		return ErrorJwtInvalidSignature()
	}

	return nil
}

// GetPayload
//  get a payload part of jwt
func GetPayload(jwt string) (payload map[string]interface{}, err error) {
	defer func() {
		var ok bool
		if r := recover(); r != nil {
			err, ok = r.(error)
			if !ok {
				panic(r)
			}
		}
	}()

	jwt_ := strings.Split(jwt, ".") //split diffrent parts

	if len(jwt_) < 2 {
		return nil, ErrorJwtInvalidLength()
	}

	payload = make(map[string]interface{})
	err = json_unmashal(byte_decoder(jwt_[jwt_payload]), &payload) //unmashal
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func BindPayload(jwt string, payload interface{}) (err error) {
	defer func() {
		var ok bool
		if r := recover(); r != nil {
			err, ok = r.(error)
			if !ok {
				panic(r)
			}
		}
	}()

	if payload == nil {
		return ErrorJwtBindPayloadNil()
	}

	if reflect.ValueOf(payload).Kind() != reflect.Ptr {
		return ErrorJwtBindPayloadNonePointer()
	}

	jwt_ := strings.Split(jwt, ".") //split diffrent parts

	if len(jwt_) != 3 {
		return ErrorJwtInvalidLength()
	}

	err = json_unmashal(byte_decoder(jwt_[jwt_payload]), payload) //unmashal
	if err != nil {
		return err
	}

	return nil
}

// HMACSHA256
//  HMACSHA256(
//  	base64UrlEncode(header) + "." +
//  	base64UrlEncode(payload),
//		secret
//    )
func HMACSHA256(header, payload string, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)                             //use hmac
	h.Write([]byte(strings.Join([]string{header, payload}, "."))) //hmac write
	signature := h.Sum(nil)                                       //hmac sum
	return signature
}

func byte_encoder(src []byte) string {
	encoder := base64.URLEncoding.EncodeToString
	return strings.ReplaceAll(encoder(src), "=", "") //remove padd
}

func byte_decoder(s string) []byte {
	decoder := base64.URLEncoding.DecodeString
	padd_recover := func(s string) string {
		mod := len(s) % 4
		if mod == 0 {
			return s
		}
		return s + strings.Repeat("=", 4-mod)
	}
	right := func(b []byte, err error) []byte {
		if err != nil {
			panic(err)
		}
		return b
	}
	return right(decoder(padd_recover(s))) //recover padd
}

var json_unmashal = json.Unmarshal
