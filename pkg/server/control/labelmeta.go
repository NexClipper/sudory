package control

import (
	"github.com/NexClipper/sudory/pkg/server/macro"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

func NewLabelMeta(name string, summary *string) metav1.LabelMeta {
	return metav1.LabelMeta{
		Name:    name,
		Summary: summary,
	}
}

func NewUuidMeta() metav1.UuidMeta {
	return metav1.UuidMeta{Uuid: macro.NewUuidString()}
}
