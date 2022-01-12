package v1

import (
	"time"
)

//database meta info
type DbMeta struct {
	//아이디 PK
	Id int64 `json:"id" orm:"id,primary"`
	//UUID
	Uuid string `json:"uuid" orm:"uuid" xorm:"unique","not null"`
	//생성자
	CreatedBy string `json:"created_by,omitempty" orm:"created_by"`
	//생성시간
	// CreatedAt serializev1.JSONTime `json:"created_at,omitempty" orm:"created_at" xorm:"extends","created"`
	CreatedAt time.Time `json:"created_at,omitempty" orm:"created_at" xorm:"created"`
	//수정자
	UpdatedBy string `json:"updated_by,omitempty" orm:"updated_by"`
	//수정시간
	// UpdatedAt serializev1.JSONTime `json:"updated_at,omitempty" orm:"updated_at" xorm:"extends","updated"`
	UpdatedAt time.Time `json:"updated_at,omitempty" orm:"updated_at" xorm:"updated"`
	//삭제자
	DeletedBy string `json:"deleted_by,omitempty" orm:"deleted_by"`
	//삭제시간, 삭제 플래그
	// DeletedAt serializev1.JSONTime `json:"deleted_at,omitempty" orm:"deleted_at" xorm:"extends","deleted"`
	DeletedAt time.Time `json:"deleted_at,omitempty" orm:"deleted_at" xorm:"deleted"`
}

//label meta info
//object extends
type LabelMeta struct {
	//label name
	Name string `json:"name" orm:"name"`
	//label summary
	Summary string `json:"summary" orm:"summary"`
	//api version
	ApiVersion string `json:"api_version,omitempty" orm:"api_version"`
}
