package vanilla

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var json_null = []byte("null")

func IsJsonNull(data []byte) bool {
	return bytes.Equal(json_null, data)
}

type NullInt struct {
	sql.NullInt64
}

func NewNullInt(i int) *NullInt {
	return &NullInt{NullInt64: sql.NullInt64{Int64: int64(i), Valid: true}}
}

func NewNullInt64(i int64) *NullInt {
	return &NullInt{NullInt64: sql.NullInt64{Int64: i, Valid: true}}
}

func (null NullInt) Int() int {
	return int(null.Int64)
}

func (null NullInt) Ptr() (out *int) {
	if null.Valid {
		i := int(null.Int64)
		out = &i
	}

	return
}

func (null NullInt) MarshalJSON() ([]byte, error) {
	if !null.Valid {
		return json_null, nil
	}

	return json.Marshal(null.Int64)
}

func (null *NullInt) UnmarshalJSON(data []byte) error {
	if IsJsonNull(data) {
		return nil
	}

	var i int
	err := json.Unmarshal(data, &i)
	null.Int64 = int64(i)
	null.Valid = err == nil

	return err
}

func (null NullInt) Print() string {
	if null.Valid {
		return strconv.FormatInt(null.Int64, 10)
	}
	return string(json_null)
}

type NullUint8 struct {
	sql.NullByte
}

func NewNullUint8(b uint8) *NullUint8 {
	return &NullUint8{NullByte: sql.NullByte{Byte: b, Valid: true}}
}

func (null NullUint8) Uint8() uint8 {
	return null.Byte
}

func (null NullUint8) Ptr() (out *uint8) {
	if null.Valid {
		i := null.Byte
		out = &i
	}

	return
}

func (null NullUint8) MarshalJSON() ([]byte, error) {
	if !null.Valid {
		return json_null, nil
	}

	return json.Marshal(uint8(null.Byte))
}

func (null *NullUint8) UnmarshalJSON(data []byte) error {
	if IsJsonNull(data) {
		return nil
	}

	var i uint8
	err := json.Unmarshal(data, &i)
	null.Byte = i
	null.Valid = err == nil

	return err
}

func (null NullUint8) Print() string {
	if null.Valid {
		return fmt.Sprintf("%v", null.NullByte)
	}
	return string(json_null)
}

type NullBool struct {
	sql.NullBool
}

func NewNullBool(ok bool) *NullBool {
	return &NullBool{NullBool: sql.NullBool{Bool: ok, Valid: true}}
}

func (nb NullBool) Ptr() (out *bool) {
	if nb.Valid {
		out = &nb.Bool
	}

	return
}

func (null NullBool) MarshalJSON() ([]byte, error) {
	if !null.Valid {
		return json_null, nil
	}

	if null.Bool {
		return []byte("true"), nil
	} else {
		return []byte("false"), nil
	}
}

func (null *NullBool) UnmarshalJSON(data []byte) error {
	if IsJsonNull(data) {
		return nil
	}

	null.Bool = bytes.Equal(data, []byte("true"))
	null.Valid = true

	return nil
}

func (null NullBool) Print() string {
	if null.Valid {
		return strconv.FormatBool(null.Bool)
	}
	return string(json_null)
}

type NullString struct {
	sql.NullString
}

func NewNullString(s string) *NullString {
	return &NullString{NullString: sql.NullString{String: s, Valid: true}}
}

func (ns NullString) Ptr() (out *string) {
	if ns.Valid {
		out = &ns.NullString.String
	}

	return
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json_null, nil
	}

	return []byte(strconv.Quote(ns.NullString.String)), nil
}

func (ns *NullString) UnmarshalJSON(data []byte) (err error) {
	if IsJsonNull(data) {
		return nil
	}

	ns.NullString.String, err = strconv.Unquote(string(data))

	ns.Valid = err == nil

	return nil
}

func (null NullString) Print() string {
	if null.Valid {
		return null.NullString.String
	}
	return string(json_null)
}

type NullTime struct {
	sql.NullTime
}

func NewNullTime(t time.Time) *NullTime {
	return &NullTime{NullTime: sql.NullTime{Time: t, Valid: true}}
}

func (nt NullTime) Ptr() (out *time.Time) {
	if nt.Valid {
		out = &nt.Time
	}

	return
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return json_null, nil
	}

	return nt.Time.MarshalJSON()
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if IsJsonNull(data) {
		return nil
	}

	err := nt.Time.UnmarshalJSON(data)
	nt.Valid = err == nil

	return err
}

func (null NullTime) Print() string {
	if null.Valid {
		return null.Time.Format(time.RFC3339Nano)
	}
	return string(json_null)
}

type NullObject struct {
	Object map[string]interface{}
	Valid  bool
}

func NewNullObject(object map[string]interface{}) *NullObject {
	return &NullObject{Object: object, Valid: true}
}

func (null *NullObject) Scan(value interface{}) error {
	null.Object = map[string]interface{}{}
	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) != nil {
		if err := json.Unmarshal([]byte(i.String), &null.Object); err != nil {
			return err
		}
		null.Valid = true
	}

	return nil
}

func (null NullObject) Value() (driver.Value, error) {
	if !null.Valid {
		return nil, nil
	}

	b, err := json.Marshal(null.Object)
	if err != nil {
		return string(b), err
	}

	return string(b), err
}

func (null NullObject) MarshalJSON() ([]byte, error) {
	if !null.Valid {
		return json_null, nil
	}

	return json.Marshal(null.Object)
}

func (null *NullObject) UnmarshalJSON(data []byte) error {
	if IsJsonNull(data) {
		return nil
	}

	null.Object = map[string]interface{}{}
	err := json.Unmarshal(data, &null.Object)
	null.Valid = err == nil

	return err
}

func (null NullObject) Print() string {
	if null.Valid {
		b, _ := json.Marshal(null.Object)
		return string(b)
	}
	return string(json_null)
}

type NullKeyValue struct {
	KeyValue map[string]string
	Valid    bool
}

func NewNullMapStringString(object map[string]string) *NullKeyValue {
	return &NullKeyValue{KeyValue: object, Valid: true}
}

func (null *NullKeyValue) Scan(value interface{}) error {
	null.KeyValue = map[string]string{}
	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) != nil {
		if err := json.Unmarshal([]byte(i.String), &null.KeyValue); err != nil {
			return err
		}
		null.Valid = true
	}

	return nil
}

func (null NullKeyValue) Value() (driver.Value, error) {
	if !null.Valid {
		return nil, nil
	}

	b, err := json.Marshal(null.KeyValue)
	if err != nil {
		return string(b), err
	}

	return string(b), err
}

func (null NullKeyValue) MarshalJSON() ([]byte, error) {
	if !null.Valid {
		return json_null, nil
	}

	return json.Marshal(null.KeyValue)
}

func (null *NullKeyValue) UnmarshalJSON(data []byte) error {
	if IsJsonNull(data) {
		return nil
	}

	null.KeyValue = map[string]string{}
	err := json.Unmarshal(data, &null.KeyValue)
	null.Valid = err == nil

	return err
}

func (null NullKeyValue) Print() string {
	if null.Valid {
		b, _ := json.Marshal(null.KeyValue)
		return string(b)
	}
	return string(json_null)
}
