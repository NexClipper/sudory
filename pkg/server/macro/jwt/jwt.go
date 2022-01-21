package macro

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
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

// JWT ew
//  make new JWT
func New(token map[string]interface{}, secret string) string {

	//define header
	var head = map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	json_mashal := json.Marshal
	json_mashal_ := func(v interface{}) []byte {
		value := func(a []byte, _ error) []byte { return a }
		return value(json_mashal(v))
	}

	encoder := base64.URLEncoding.EncodeToString
	encoder_ := func(src []byte) string {
		return strings.ReplaceAll(encoder(src), "=", "")
	}

	header := encoder_(json_mashal_(head))                         //make header
	payload := encoder_(json_mashal_(token))                       //make payload
	signature := encoder_(_HMACSHA256(header, payload, secret))    //make signature
	jwt := strings.Join([]string{header, payload, signature}, ".") //make jwt

	return jwt
}

// JWT Verify
//  valied signature
func Verify(jwt string, secret string) error {

	parts := strings.Split(jwt, ".") //split diffrent parts

	if len(parts) != 3 {
		return ErrorJwtInvalidLength()
	}

	signature_ := _HMACSHA256(parts[jwt_header], parts[jwt_payload], secret) //make signature

	if parts[jwt_signature] != string(signature_) { //compare
		return ErrorJwtInvalidSignature()
	}

	return nil
}

// GetToken
//  get a token part of jwt
func GetToken(jwt string, secret string) (map[string]interface{}, error) {
	var err error

	padd_recover := func(s string) string {
		mod := len(s) % 4
		if mod == 0 {
			return s
		}

		padd := strings.Repeat("=", 4-mod)
		return s + padd
	}
	decoder := base64.URLEncoding.DecodeString
	decoder_ := func(s string) []byte {
		value := func(a []byte, _ error) []byte { return a }
		return value(decoder(padd_recover(s)))
	}
	json_unmashal := json.Unmarshal

	jwt_ := strings.Split(jwt, ".") //split diffrent parts

	if len(jwt_) != 3 {
		return nil, ErrorJwtInvalidLength()
	}

	token := make(map[string]interface{})
	err = json_unmashal(decoder_(jwt_[jwt_payload]), &token) //token unmashal
	if err != nil {
		return nil, err
	}

	return token, nil
}

// HMACSHA256
//  HMACSHA256(
//  	base64UrlEncode(header) + "." +
//  	base64UrlEncode(payload),
//		secret
//    )
func _HMACSHA256(header, payload, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))                     //use hmac
	h.Write([]byte(strings.Join([]string{header, payload}, "."))) //hmac write
	signature := h.Sum(nil)                                       //hmac sum
	return signature
}
