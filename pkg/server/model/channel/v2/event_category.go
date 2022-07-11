package v2

import (
	"database/sql"
	"database/sql/driver"
)

//go:generate go run github.com/abice/go-enum --file=event_category.go --names --nocase=true

/* ENUM(
	NaV
	nonspecified
	client-auth-accept
	service-polling-out
	service-polling-in
)
*/
type EventCategory int

// func (enum EventCategory) MarshalJSON() ([]byte, error) {
// 	return []byte(strconv.Quote(enum.String())), nil
// }

// func (enum *EventCategory) UnmarshalJSON(data []byte) (err error) {
// 	s, err := strconv.Unquote(string(data))
// 	if err != nil {
// 		return err
// 	}
// 	*enum, err = ParseEventCategory(s)
// 	return
// }

// Scan implements the Scanner interface.
func (n *EventCategory) Scan(value interface{}) error {
	nn := sql.NullInt32{}

	if err := nn.Scan(value); err != nil {
		return err
	}

	*n = EventCategory(nn.Int32)

	return nil
}

// Value implements the driver Valuer interface.
func (n EventCategory) Value() (driver.Value, error) {
	if n == EventCategoryNaV {
		return nil, nil
	}
	return int64(n), nil
}
