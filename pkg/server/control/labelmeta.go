package control

import (
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

func NewLabelMeta(name, summary *string) metav1.LabelMeta {
	const api_version = "v1"
	return metav1.LabelMeta{
		ApiVersion: newist.String(api_version),
		Name:       name,
		Summary:    summary,
	}
}

func NewUuidMeta() metav1.UuidMeta {
	return metav1.UuidMeta{Uuid: macro.NewUuidString()}
}
