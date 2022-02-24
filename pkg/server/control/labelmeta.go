package control

import (
	"github.com/NexClipper/sudory/pkg/server/macro"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

func NewLabelMeta(name string, summary *string) metav1.LabelMeta {
	const api_version = "v1"
	return metav1.LabelMeta{
		Uuid:       macro.NewUuidString(),
		ApiVersion: api_version,
		Name:       name,
		Summary:    summary,
	}
}

func LabelMeta(uuid, name string, summary *string) metav1.LabelMeta {
	const api_version = "v1"
	return metav1.LabelMeta{
		Uuid:       uuid,
		ApiVersion: api_version,
		Name:       name,
		Summary:    summary,
	}
}
