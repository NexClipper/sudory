package v1

import (
	"encoding/json"
	"sort"

	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
)

//Template Property
type TemplateProperty struct {
	Origin string `json:"origin" xorm:"'origin' varchar(255) null comment('origin')"`
}

//DATABASE SCHEMA: TEMPLATE
type Template struct {
	metav1.DbMeta    `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta  `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	TemplateProperty `json:",inline" xorm:"extends"` //inline property
}

func (Template) TableName() string {
	return "template"
}

type HttpReqTemplate_Create struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"`                      //inline labelmeta
	Origin           string                                               `json:"origin,omitempty"`
	Commands         []commandv1.HttpReqTemplateCommand_Create_ByTemplate `json:"commands,omitempty"`
}

type HttpReqTemplate_Update struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	Origin           string                          `json:"origin,omitempty"`
}

type HttpRspTemplate struct {
	Template `json:",inline"`
	Commands []commandv1.TemplateCommand `json:"commands,omitempty"`
}

func (object HttpRspTemplate) MarshalJSON() ([]byte, error) {
	//sort commands by sequence
	sort.Slice(object.Commands, func(i, j int) bool {
		var a, b int32 = 0, 0
		if object.Commands[i].Sequence != nil {
			a = *object.Commands[i].Sequence
		}
		if object.Commands[j].Sequence != nil {
			b = *object.Commands[j].Sequence
		}
		return a < b
	})

	v := struct {
		Template `json:",inline"`
		Commands []commandv1.TemplateCommand `json:",inline"`
	}{
		object.Template,
		object.Commands,
	}

	return json.Marshal(v)
}
