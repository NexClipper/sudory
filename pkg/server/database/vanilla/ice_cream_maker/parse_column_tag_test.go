package ice_cream_maker_test

import (
	"reflect"
	"testing"
	"time"

	vanilla "github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
)

func TestColumnTags(t *testing.T) {
	objs := []interface{}{
		ServiceStep_essential{},
		ServiceStep{},
	}

	for _, obj := range objs {
		columninfo := vanilla.ParseColumnTag(reflect.TypeOf(obj), "")
		for _, columninfo := range columninfo {
			println(columninfo.String())
		}
	}
}

func TestFieldNames(t *testing.T) {
	objs := []interface{}{
		ServiceStep_essential{},
		ServiceStep{},
	}

	for _, obj := range objs {
		names := vanilla.FieldNames(reflect.TypeOf(obj))
		for _, name := range names {
			println(name)
		}
	}
}

func TypeFieldNames(obj interface{}) []string {
	tof := reflect.TypeOf(obj)
	tags := make([]string, tof.NumField())

	for i := 0; i < tof.NumField(); i++ {
		tags[i] = tof.Field(i).Name
	}

	return tags
}

// not working
func TypeFieldPtr(obj interface{}) []interface{} {
	vof := reflect.ValueOf(obj)
	ptr := make([]interface{}, vof.NumField())

	for i := 0; i < vof.NumField(); i++ {
		ptr[i] = vof.Field(i).Interface()
	}

	return ptr
}

type ServiceStep_essential struct {
	Name          string `column:"name"          json:"name,omitempty"`
	Summary       string `column:"summary"       json:"summary,omitempty"`
	Method        string `column:"method"        json:"method,omitempty"`
	Args          string `column:"args"          json:"args,omitempty"`
	Result_filter string `column:"result_filter" json:"result_filter,omitempty"`
}

func (ServiceStep_essential) TableName() string {
	return "service_step"
}

type ServiceStep struct {
	Uuid     string    `column:"uuid"     json:"uuid,omitempty"`     //pk
	Sequence string    `column:"sequence" json:"sequence,omitempty"` //pk
	Created  time.Time `column:"created"  json:"created,omitempty"`  //pk

	ServiceStep_essential `json:",inline"`
}

func (row *ServiceStep) Values() []interface{} {
	return []interface{}{
		row.Uuid,
		row.Sequence,
		row.Created,
		row.Name,
		row.Summary,
		row.Method,
		row.Args,
		row.Result_filter,
	}
}

func (row *ServiceStep) Dests() []interface{} {
	return []interface{}{
		&row.Uuid,
		&row.Sequence,
		&row.Created,
		&row.Name,
		&row.Summary,
		&row.Method,
		&row.Args,
		&row.Result_filter,
	}
}
