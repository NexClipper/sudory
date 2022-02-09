package control

import (
	"github.com/NexClipper/sudory/pkg/server/macro"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

const api_version = "v1"

func NewLabelMeta(name, summary string) metav1.LabelMeta {
	return metav1.LabelMeta{
		Uuid:       macro.NewUuidString(),
		ApiVersion: api_version,
		Name:       name,
		Summary:    &summary,
	}
}
