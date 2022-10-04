//go:generate go run github.com/abice/go-enum --file=channel_format.go --names --nocase=true
package v3

import (
	"database/sql"
	"database/sql/driver"
)

/* ENUM(
	disable
	fields
	jq
)
*/
type FormatType int

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

type Format_update = Format_property

type Format_property struct {
	FormatType FormatType `column:"format_type,default(0)"  json:"format_type,omitempty" enums:"0,1,2"` // enums:"disable(0), fields(1), jq(2)"
	FormatData string     `column:"format_data,default('')" json:"format_data,omitempty"`
}

type Format struct {
	Uuid string `column:"uuid"                    json:"uuid,omitempty"` // pk

	Format_property `json:",inline"`

	// Created vanilla.NullTime `column:"created"                 json:"created,omitempty"     swaggertype:"string"`
	// Updated vanilla.NullTime `column:"updated"                 json:"updated,omitempty"     swaggertype:"string"`
}

func (Format) TableName() string {
	return "managed_channel_format"
}
