//go:generate go run github.com/abice/go-enum --file=channel_filter.go --names --nocase=true
package v2

import (
	"fmt"
	"strings"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

/* ENUM(
	NaV
	EQ
	NOT
	GT
	GTE
	LT
	LTE
)
*/
type FilterOp int

// func (enum FilterOp) MarshalJSON() ([]byte, error) {
// 	return []byte(strconv.Quote(enum.String())), nil
// }

// func (enum *FilterOp) UnmarshalJSON(data []byte) (err error) {
// 	s, err := strconv.Unquote(string(data))
// 	if err != nil {
// 		return err
// 	}
// 	*enum, err = ParseFilterOp(s)
// 	return
// }

type Filter struct {
	Uuid string `column:"uuid"         json:"uuid,omitempty"` // pk
	// enums:"NaV(0), EQ(1), NOT(2), GT(3), GTE(4), LT(5), LTE(6)"
	FilterOp    FilterOp         `column:"filter_op"    json:"filter_op,omitempty"   enums:"0,1,2,3,4,5,6"`
	FilterKey   string           `column:"filter_key"   json:"filter_key,omitempty"`
	FilterValue string           `column:"filter_value" json:"filter_value,omitempty"`
	Created     time.Time        `column:"created"      json:"created,omitempty"`
	Updated     vanilla.NullTime `column:"updated"      json:"updated,omitempty"     swaggertype:"string"`
	Deleted     vanilla.NullTime `column:"deleted"      json:"deleted,omitempty"     swaggertype:"string"`
}

func (Filter) TableName() string {
	return "managed_channel_filter"
}

type Channel_EventCategory_Filter struct {
	Uuid          string        `column:"uuid"`
	EventCategory EventCategory `column:"event_category"`
	FilterOp      FilterOp      `column:"filter_op"`
	FilterKey     string        `column:"filter_key"`
	FilterValue   string        `column:"filter_value"`
}

func (Channel_EventCategory_Filter) TableName() string {
	q := `(
		SELECT %v /**columns**/
		  FROM %v A /**managed_channel A**/
		  INNER JOIN %v B /**managed_channel_filter B**/
				ON A.uuid = B.uuid 
		) X`

	columns := []string{
		"A.uuid",
		fmt.Sprintf("IFNULL(A.event_category, '%v') AS event_category", int(EventCategoryNaV)),
		"B.filter_key",
		"B.filter_value",
	}
	A := "managed_channel"
	B := "managed_channel_filter"
	return fmt.Sprintf(q, strings.Join(columns, ", "), A, B)
}
