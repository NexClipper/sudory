package v3

import (
	"time"
)

// HttpRsp_ClientServicePolling
//  http responce body; client service polling
type HttpRsp_ClientServicePolling struct {
	// Uuid    string                `json:"uuid"`              //pk
	// Created time.Time             `json:"created,omitempty"` //pk
	Service `json:",inline"`
	Steps   []ServiceStep `json:"steps,omitempty"`
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
	Service `json:",inline"`
	Steps   []ServiceStep `json:"steps,omitempty"`
}

type HttpRsp_Service_create struct {
	Service_create `json:",inline"`
	Steps          []ServiceStep_create `json:"steps,omitempty"`
}

type HttpReq_Service_create struct {
	Uuid              string `json:"uuid"` //pk
	Name              string `json:"name,omitempty"`
	Summary           string `json:"summary,omitempty"`
	ClusterUuid       string `json:"cluster_uuid,omitempty"`
	TemplateUuid      string `json:"template_uuid,omitempty"`
	SubscribedChannel string `json:"subscribed_channel,omitempty"`
	Steps             []struct {
		Args map[string]interface{} `json:"args,omitempty"`
	} `json:"steps,omitempty"`
}

type HttpRsp_ServiceStep = ServiceStep

type HttpRsp_ServiceResult = ServiceResult
