package v2

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"time"
)

// NullInt is an alias for sql.NullInt data type
type NullInt int

func (ni NullInt) Int() int {
	return (int)(ni)
}

// Scan implements the Scanner interface for NullInt64
func (ni *NullInt) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*ni = (NullInt)(int(i.Int64))
	} else {
		*ni = (NullInt)(int(i.Int64))
	}
	return nil
}

// // MarshalJSON for NullInt64
// func (ni *NullInt) MarshalJSON() ([]byte, error) {
// 	if !ni.Valid {
// 		return []byte("null"), nil
// 	}
// 	return json.Marshal(ni.Int64)
// }

// func (b *NullInt) UnmarshalJSON(data []byte) error {

// 	var v []interface{}
// 	if err := json.Unmarshal(data, &v); err != nil {
// 		return err
// 	}

// 	return nil
// }

type NullString string

func (ns NullString) String() string {
	return (string)(ns)
}

// Scan implements the Scanner interface for NullString
func (ns *NullString) Scan(value interface{}) error {
	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = (NullString)(i.String)
	} else {
		*ns = (NullString)(i.String)
	}
	return nil
}

// // MarshalJSON for NullInt64
// func (ni *NullString) MarshalJSON() ([]byte, error) {
// 	if !ni.Valid {
// 		return []byte("null"), nil
// 	}
// 	return json.Marshal(ni.String)
// }

// func (b *NullString) UnmarshalJSON(data []byte) error {
// 	b.String = string(data)
// 	return nil
// }

type NullTime time.Time

func (nt NullTime) Time() time.Time {
	return (time.Time)(nt)
}

// Scan implements the Scanner interface for NullTime
func (nt *NullTime) Scan(value interface{}) error {
	var i sql.NullTime
	if err := i.Scan(value); err != nil {
		return err
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*nt = (NullTime)(i.Time)
	} else {
		*nt = (NullTime)(i.Time)
	}
	return nil
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	t := (time.Time)(nt)
	// if t.IsZero() {
	// 	return []byte{}, nil
	// }
	return t.MarshalJSON()
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	t := (*time.Time)(nt)
	err := t.UnmarshalJSON(data)
	if err == nil {
		*nt = *(*NullTime)(t)
	}
	return err
}

type NullJson map[string]interface{}

func (nj NullJson) Json() map[string]interface{} {
	return (map[string]interface{})(nj)
}

// Scan implements the Scanner interface for NullTime
func (nj *NullJson) Scan(value interface{}) error {
	m := map[string]interface{}{}
	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*nj = (NullJson)(m)
	} else {

		if err := json.Unmarshal([]byte(i.String), &m); err != nil {
			return err
		}

		*nj = (NullJson)(m)
	}
	return nil
}

func (nj NullJson) Value() (driver.Value, error) {
	b, err := json.Marshal(nj)
	if err != nil {
		return string(b), err
	}

	return string(b), err
}
