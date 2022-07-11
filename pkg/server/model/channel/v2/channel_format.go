//go:generate go run github.com/abice/go-enum --file=channel_format.go --names --nocase=true
package v2

import (
	"database/sql"
	"database/sql/driver"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

/* ENUM(
	disable
	fields
	jq
)
*/
type FormatType int

// func (enum FormatType) MarshalJSON() ([]byte, error) {
// 	return []byte(strconv.Quote(enum.String())), nil
// }

// func (enum *FormatType) UnmarshalJSON(data []byte) (err error) {
// 	s, err := strconv.Unquote(string(data))
// 	if err != nil {
// 		return err
// 	}
// 	*enum, err = ParseFormatType(s)
// 	return
// }

// Scan implements the Scanner interface.
func (n *FormatType) Scan(value interface{}) error {
	nn := sql.NullInt32{}

	if err := nn.Scan(value); err != nil {
		return err
	}

	*n = FormatType(nn.Int32)

	return nil
}

// Value implements the driver Valuer interface.
func (n FormatType) Value() (driver.Value, error) {
	return int64(n), nil
}

type Format_essential struct {
	// enums:"disable(0), fields(1), jq(2)"
	FormatType FormatType `column:"format_type,default(0)" json:"format_type,omitempty" enums:"0,1,2"`
	FormatData string     `column:"format_data,default('')" json:"format_data,omitempty"`
}
type Format_property struct {
	Format_essential `json:",inline"`

	Created vanilla.NullTime `column:"created" json:"created,omitempty" swaggertype:"string"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
}

func (Format_property) TableName() string {
	return "managed_channel_format"
}

type Format struct {
	Format_property `json:",inline"`

	Uuid string `column:"uuid" json:"uuid,omitempty"` // pk
}
