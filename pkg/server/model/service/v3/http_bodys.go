package service

import (
	"encoding/json"
	"strconv"
	"strings"
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
	Name              string   `json:"name,omitempty"`
	Summary           string   `json:"summary,omitempty"`
	ClusterUuid       []string `json:"cluster_uuid,omitempty"`
	TemplateUuid      string   `json:"template_uuid,omitempty"`
	SubscribedChannel string   `json:"subscribed_channel,omitempty"`
	Steps             []struct {
		Args map[string]interface{} `json:"args,omitempty"`
	} `json:"steps,omitempty"`

	IsMultiCluster bool `json:"-"`
}

func (obj HttpReq_Service_create) MarshalJSON() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *HttpReq_Service_create) UnmarshalJSON(bytes []byte) error {

	type T struct {
		Name              string          `json:"name,omitempty"`
		Summary           string          `json:"summary,omitempty"`
		ClusterUuid       json.RawMessage `json:"cluster_uuid,omitempty"`
		TemplateUuid      string          `json:"template_uuid,omitempty"`
		SubscribedChannel string          `json:"subscribed_channel,omitempty"`
		Steps             []struct {
			Args map[string]interface{} `json:"args,omitempty"`
		} `json:"steps,omitempty"`
	}

	var v T
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	var cluster_uuid = []string{}
	if 0 < len(v.ClusterUuid) &&
		string(v.ClusterUuid)[0] == '[' &&
		string(v.ClusterUuid)[len(string(v.ClusterUuid))-1] == ']' {

		if err := json.Unmarshal(v.ClusterUuid, &cluster_uuid); err != nil {
			return err
		}

		obj.IsMultiCluster = true
	}

	if 0 < len(v.ClusterUuid) &&
		string(v.ClusterUuid)[0] == '"' &&
		string(v.ClusterUuid)[len(string(v.ClusterUuid))-1] == '"' {

		s := string(v.ClusterUuid)
		s, _ = strconv.Unquote(s)
		s = strings.TrimSpace(s)

		cluster_uuid = append(cluster_uuid, s)

		obj.IsMultiCluster = false
	}

	// obj.Uuid = v.Uuid
	obj.Name = v.Name
	obj.Summary = v.Summary
	obj.ClusterUuid = cluster_uuid
	obj.TemplateUuid = v.TemplateUuid
	obj.SubscribedChannel = v.SubscribedChannel
	obj.Steps = v.Steps
	obj.Steps = v.Steps

	return nil
}

type HttpRsp_ServiceStep = ServiceStep

type HttpRsp_ServiceResult = ServiceResult
