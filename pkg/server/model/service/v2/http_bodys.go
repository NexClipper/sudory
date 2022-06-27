package v2

import (
	"time"
)

// HttpRsp_ClientServicePolling
//  http responce body; client service polling
type HttpRsp_ClientServicePolling struct {
	// Uuid    string                `json:"uuid"`              //pk
	// Created time.Time             `json:"created,omitempty"` //pk
	Service_status `json:",inline"`
	Steps          []ServiceStep_tangled `json:"steps,omitempty"`
}

// HttpReq_ClientServiceUpdate
//  http request body; client service update
type HttpReq_ClientServiceUpdate struct {
	Uuid     string     `json:"uuid"`     //pk
	Sequence int        `json:"sequence"` //pk
	Status   StepStatus `json:"status"`
	Result   string     `json:"result"` //StepStatus 값에 따라; 결과 혹은 에러 메시지
	Started  time.Time  `json:"started"`
	Ended    time.Time  `json:"ended"`
}

type HttpRsp_Service struct {
	Service_tangled `json:",inline"`
	Steps           []ServiceStep_tangled `json:"steps,omitempty"`
}

type HttpRsp_Service_status struct {
	Service_status `json:",inline"`
	Steps          []ServiceStep_tangled `json:"steps,omitempty"`
}

type HttpRsp_Service_create struct {
	Service `json:",inline"`
	Steps   []ServiceStep `json:"steps,omitempty"`
}

type HttpReq_Service_Create struct {
	Name              string `json:"name,omitempty"`
	Summary           string `json:"summary,omitempty"`
	ClusterUuid       string `json:"cluster_uuid,omitempty"`
	TemplateUuid      string `json:"template_uuid,omitempty"`
	SubscribedChannel string `json:"subscribed_channel,omitempty"`
}

type HttpReq_ServiceStep_Create struct {
	Args map[string]interface{} `json:"args,omitempty"`
}

type HttpReq_Service_create struct {
	Uuid                   string `json:"uuid"` //pk
	HttpReq_Service_Create `json:",inline"`
	Steps                  []HttpReq_ServiceStep_Create `json:"steps,omitempty"`
}

type HttpRsp_ServiceStep struct {
	ServiceStep_tangled `json:",inline"`
}

// type HttpReq_ServiceUpdate struct {
// 	Service Service_essential       `json:",inline"`
// 	Steps   []ServiceStep_essential `json:",inline"`
// }
