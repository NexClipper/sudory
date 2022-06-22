package v2

import (
	"fmt"
	"strings"
	"time"

	crypto "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
	noxorm "github.com/NexClipper/sudory/pkg/server/model/noxorm/v2"
)

type Service_essential struct {
	Name              string            `column:"name"               json:"name,omitempty"`
	Summary           noxorm.NullString `column:"summary"            json:"summary,omitempty"`
	ClusterUuid       string            `column:"cluster_uuid"       json:"cluster_uuid,omitempty"`
	TemplateUuid      string            `column:"template_uuid"      json:"template_uuid,omitempty"`
	StepCount         int               `column:"step_count"         json:"step_count,omitempty"`
	SubscribedChannel noxorm.NullString `column:"subscribed_channel" json:"subscribed_channel,omitempty"`
	OnCompletion      OnCompletion      `column:"on_completion"      json:"on_completion,omitempty"`
}

func (Service_essential) TableName() string {
	return "service"
}

type Service struct {
	Uuid              string    `column:"uuid"    json:"uuid,omitempty"`    //pk
	Created           time.Time `column:"created" json:"created,omitempty"` //pk
	Service_essential `json:",inline"`
}

type ServiceStatus_essential struct {
	AssignedClientUuid string            `column:"assigned_client_uuid" json:"assigned_client_uuid,omitempty"`
	StepPosition       int               `column:"step_position"        json:"step_position,omitempty"`
	Status             StepStatus        `column:"status"               json:"status,omitempty"`
	Message            noxorm.NullString `column:"message"              json:"message,omitempty"`
}

func (ServiceStatus_essential) TableName() string {
	return "service_status"
}

type ServiceStatus struct {
	Uuid    string    `column:"uuid"    json:"uuid,omitempty"`    //pk
	Created time.Time `column:"created" json:"created,omitempty"` //pk

	ServiceStatus_essential `json:",inline"`
}

type ServiceResults_essential struct {
	ResultType ResultType          `column:"result_type" json:"result_type,omitempty"`
	Result     crypto.CryptoString `column:"result"      json:"result,omitempty"`
}

func (ServiceResults_essential) TableName() string {
	return "service_result"
}

type ServiceResult struct {
	Uuid    string    `column:"uuid"    json:"uuid,omitempty"`    //pk
	Created time.Time `column:"created" json:"created,omitempty"` //pk

	ServiceResults_essential `json:",inline"`
}

type Service_tangled struct {
	// Uuid    string    `column:"uuid"    json:"uuid"`    //pk
	// Created time.Time `column:"created" json:"created"` //pk
	// Updated time.Time `column:"updated" json:"updated"` //pk

	// Service_essential        `json:",inline"` //service
	// ServiceStatus_essential  `json:",inline"` //status
	// ServiceResults_essential `json:",inline"` //result

	Service                  `json:",inline"` //service
	ServiceStatus_essential  `json:",inline"` //status
	ServiceResults_essential `json:",inline"` //result

	Updated time.Time `column:"updated" json:"updated"` //pk
}

/*
`
SELECT A.uuid, A.created,
       name, summary, cluster_uuid, template_uuid, step_count, subscribed_channel, on_completion,
       B.created AS updated, assigned_client_uuid, step_position, status, message, result_type, result
  FROM service A
  LEFT JOIN service_status B
         ON A.uuid = B.uuid
        AND B.created = (
            SELECT MAX(B.created) AS MAX_created
              FROM service_status B
             WHERE A.uuid = B.uuid
			)
  LEFT JOIN service_result C
         ON A.uuid = C.uuid
        AND C.created = (
            SELECT MAX(C.created) AS MAX_created
              FROM service_result C
             WHERE A.uuid = C.uuid
			)
`
*/
func (record Service_tangled) TableName() string {
	q := `(
    SELECT %v /**columns**/
      FROM %v A /**service A**/
      LEFT JOIN %v B /**service_status B**/
             ON A.uuid = B.uuid 
            AND B.created = (
                SELECT MAX(B.created) AS MAX_created
                  FROM %v B /**service_status B**/
                 WHERE A.uuid = B.uuid 
                )
      LEFT JOIN %v C /**service_result C**/
             ON A.uuid = C.uuid
            AND C.created = (
                SELECT MAX(C.created) AS MAX_created
                  FROM %v C /**service_result C**/
                 WHERE A.uuid = C.uuid
                )
    ) X`

	columns := []string{
		"A.uuid",
		"A.created",
		"B.created AS updated",
	}
	columns = append(columns, Service_essential{}.ColumnNames()...)
	columns = append(columns, ServiceStatus_essential{}.ColumnNames()...)
	columns = append(columns, ServiceResults_essential{}.ColumnNames()...)
	A := record.Service.TableName()
	B := record.ServiceStatus_essential.TableName()
	C := record.ServiceResults_essential.TableName()
	return fmt.Sprintf(q, strings.Join(columns, ", "), A, B, B, C, C)
}
