//go:generate go run github.com/abice/go-enum --file=notifier_type.go --names --nocase=true
package v2

import (
	"database/sql"
	"database/sql/driver"
)

/* ENUM(
	NaV
	console
	webhook
	rabbitmq
)
*/
type NotifierType int

// func (enum NotifierType) MarshalJSON() ([]byte, error) {
// 	return []byte(strconv.Quote(enum.String())), nil
// }

// func (enum *NotifierType) UnmarshalJSON(data []byte) (err error) {
// 	s, err := strconv.Unquote(string(data))
// 	if err != nil {
// 		return err
// 	}
// 	*enum, err = ParseNotifierType(s)
// 	return
// }

// Scan implements the Scanner interface.
func (n *NotifierType) Scan(value interface{}) error {
	nn := sql.NullInt32{}

	if err := nn.Scan(value); err != nil {
		return err
	}

	*n = NotifierType(nn.Int32)

	return nil
}

// Value implements the driver Valuer interface.
func (n NotifierType) Value() (driver.Value, error) {
	if n == NotifierTypeNaV {
		return nil, nil
	}
	return int64(n), nil
}
