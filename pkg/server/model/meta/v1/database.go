package v1

import (
	"time"
)

//database meta info
type DbMeta struct {
	//아이디 PK
	Id uint64 `json:"id" xorm:"bigint notnull pk autoincr 'id' comment('db meta's id')"`
	//생성시간
	Created *time.Time `json:"created,omitempty" xorm:"datetime null created 'created' comment('db meta's created')"`
	//수정시간
	Updated *time.Time `json:"updated,omitempty" xorm:"datetime null updated 'updated' comment('db meta's updated')"`
	//삭제시간, 삭제 플래그
	Deleted *time.Time `json:"deleted,omitempty" xorm:"datetime null deleted 'deleted' comment('db meta's deleted')"`
}

//label meta info
//object extends
type LabelMeta struct {
	//label name
	Name *string `json:"name" xorm:"varchar(255) notnull 'name' comment('label meta's name')"`
	//label summary
	Summary *string `json:"summary,omitempty" xorm:"varchar(255) null 'summary' comment('label meta's summary')"`
	//api version
	ApiVersion *string `json:"api_version" xorm:"varchar(255) notnull 'api_version' comment('label meta's api version')"`
}

//uuid meta info
type UuidMeta struct {
	//UUID
	Uuid string `json:"uuid" xorm:"char(32) notnull unique 'uuid' comment('uuid meta's uuid')"`
}
