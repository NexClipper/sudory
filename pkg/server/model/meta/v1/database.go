package v1

import (
	"time"
)

//database meta info
type DbMeta struct {
	Id      uint64     `json:"id"                xorm:"'id'      bigint   notnull pk autoincr comment('id')"`
	Created *time.Time `json:"created,omitempty" xorm:"'created' datetime null    created     comment('created')"`
	Updated *time.Time `json:"updated,omitempty" xorm:"'updated' datetime null    updated     comment('updated')"`
	Deleted *time.Time `json:"deleted,omitempty" xorm:"'deleted' datetime null    deleted     comment('deleted')"`
}

//label meta info
//object extends
type LabelMeta struct {
	Name    string  `json:"name"              xorm:"'name'    varchar(255) notnull comment('name')"`
	Summary *string `json:"summary,omitempty" xorm:"'summary' varchar(255) null    comment('summary')"`
}

//uuid meta info
type UuidMeta struct {
	Uuid string `json:"uuid" xorm:"'uuid' char(32) notnull unique comment('uuid')"`
}
