package vanilla_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

func TestNullInt(t *testing.T) {
	// obj := map[string]interface{}{}
	{
		obj := vanilla.NewNullInt(987654321)
		expected := `987654321`
		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullInt)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if v.Int() != 987654321 {
				return fmt.Errorf("diff value")
			}

			if v.Valid != true {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}
	}
	{
		obj := &vanilla.NullInt{}
		expected := `null`
		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullInt)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if v.Int() != 0 {
				return fmt.Errorf("diff value")
			}

			if v.Valid != false {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}
	}

}

func TestNullUint8(t *testing.T) {
	{
		obj := vanilla.NewNullUint8(128)
		expected := `128`

		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullUint8)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if v.Uint8() != 128 {
				return fmt.Errorf("diff value")
			}

			if v.Valid != true {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}
	}
	{
		obj := &vanilla.NullUint8{}
		expected := `null`

		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullUint8)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if v.Uint8() != 0 {
				return fmt.Errorf("diff value")
			}

			if v.Valid != false {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}
	}
}

func TestNullBool(t *testing.T) {
	{
		obj := vanilla.NewNullBool(true)
		expected := `true`
		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullBool)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if v.Bool != true {
				return fmt.Errorf("diff value")
			}

			if v.Valid != true {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}
	}
	{
		obj := vanilla.NewNullBool(false)
		expected := `false`
		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullBool)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if v.Bool != false {
				return fmt.Errorf("diff value")
			}

			if v.Valid != true {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}
	}
	{
		obj := &vanilla.NullBool{}
		expected := `null`
		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullBool)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if v.Bool != false {
				return fmt.Errorf("diff value")
			}

			if v.Valid != false {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}
	}
}

func TestNullString(t *testing.T) {
	{
		obj := vanilla.NewNullString("foo")
		expected := `"foo"`

		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullString)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if v.String != `"foo"` {
				return fmt.Errorf("diff value")
			}

			if v.Valid != true {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}
	}
	{
		obj := &vanilla.NullString{}
		expected := `null`
		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullString)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if v.String != `` {
				return fmt.Errorf("diff value")
			}

			if v.Valid != false {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}
	}
}

func TestNullTime(t *testing.T) {

	{
		obj := vanilla.NewNullTime(time.Unix(1656567675, 0))
		expected := `"2022-06-30T05:41:15Z"`

		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullTime)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if tt, _ := time.Parse(time.RFC3339, `2022-06-30T05:41:15Z`); tt.Unix() != v.Time.Unix() {
				return fmt.Errorf("diff value")
			}

			if v.Valid != true {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}
	}
	{
		obj := &vanilla.NullTime{}
		expected := `null`

		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullTime)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if tt, _ := time.Parse(time.RFC3339, `0001-01-01T00:00:00Z`); tt.Unix() != v.Time.Unix() {
				return fmt.Errorf("diff value")
			}

			if v.Valid != false {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}

	}
}

func TestNullObject(t *testing.T) {
	{
		obj := vanilla.NewNullObject(map[string]interface{}{
			"string": "foo",
			"number": 987654321,
			"bool":   true,
			"nil":    nil,
		})
		expected := `{"bool":true,"nil":null,"number":987654321,"string":"foo"}`

		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullObject)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if v.Object["string"] != `foo` {
				return fmt.Errorf("diff value")
			}
			if v.Object["number"] != float64(987654321) {
				return fmt.Errorf("diff value")
			}
			if v.Object["bool"] != true {
				return fmt.Errorf("diff value")
			}
			if v.Object["nil"] != nil {
				return fmt.Errorf("diff value")
			}

			if v.Valid != true {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}
	}
	{
		obj := &vanilla.NullObject{}
		expected := `null`
		if err := JsonMarshal(
			obj,
			expected,
		); err != nil {
			t.Error(err)
		}

		if err := JsonUnmarshal(expected, obj, func(a interface{}) error {
			v, ok := a.(*vanilla.NullObject)
			if !ok {
				return fmt.Errorf("diff type")
			}

			if v.Object != nil {
				return fmt.Errorf("diff value")
			}

			if v.Valid != false {
				return fmt.Errorf("diff valid")
			}

			return nil
		}); err != nil {
			t.Error(err)
		}
	}
}

func JsonMarshal(obj interface{}, expected string) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	if !bytes.Equal([]byte(expected), b) {
		return fmt.Errorf("expected=\"%s\" actual=\"%s\"", expected, b)
	}

	return nil
}

func JsonUnmarshal(j string, obj interface{}, diff func(a interface{}) error) error {

	if err := json.Unmarshal([]byte(j), obj); err != nil {
		return err
	}

	if err := diff(obj); err != nil {
		return fmt.Errorf("is diffrent: %w", err)
	}

	return nil
}

func init() {
	time.Local = time.UTC
}
