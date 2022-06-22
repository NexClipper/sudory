package v2

import (
	"fmt"
	"strings"
	"time"

	noxorm "github.com/NexClipper/sudory/pkg/server/model/noxorm/v2"
)

type ServiceStep_essential struct {
	Name         string            `column:"name"          json:"name,omitempty"`
	Summary      noxorm.NullString `column:"summary"       json:"summary,omitempty"`
	Method       string            `column:"method"        json:"method,omitempty"`
	Args         noxorm.NullJson   `column:"args"          json:"args,omitempty"`
	ResultFilter noxorm.NullString `column:"result_filter" json:"result_filter,omitempty"`
}

func (ServiceStep_essential) TableName() string {
	return "service_step"
}

type ServiceStep struct {
	Uuid     string    `column:"uuid"     json:"uuid,omitempty"`     //pk
	Sequence int       `column:"sequence" json:"sequence,omitempty"` //pk
	Created  time.Time `column:"created"  json:"created,omitempty"`  //pk

	ServiceStep_essential `json:",inline"`
}

type ServiceStepStatus_essential struct {
	Status  StepStatus      `column:"status"  json:"status,omitempty"`
	Started noxorm.NullTime `column:"started" json:"started,omitempty"`
	Ended   noxorm.NullTime `column:"ended"   json:"ended,omitempty"`
}

func (ServiceStepStatus_essential) TableName() string {
	return "service_step_status"
}

type ServiceStepStatus struct {
	Uuid     string    `column:"uuid"     json:"uuid,omitempty"`     //pk
	Sequence int       `column:"sequence" json:"sequence,omitempty"` //pk
	Created  time.Time `column:"created"  json:"created,omitempty"`  //pk

	ServiceStepStatus_essential `json:",inline"`
}

type ServiceStep_tangled struct {
	// Uuid     string    `column:"uuid"     json:"uuid,omitempty"`     //pk
	// Sequence int       `column:"sequence" json:"sequence,omitempty"` //pk
	// Created  time.Time `column:"created"  json:"created,omitempty"`  //pk
	// Updated  time.Time `column:"updated"  json:"updated,omitempty"`  //pk

	// ServiceStep_essential       `json:",inline"` //step
	// ServiceStepStatus_essential `json:",inline"` //status

	ServiceStep                 `json:",inline"` //step
	ServiceStepStatus_essential `json:",inline"` //status

	Updated time.Time `column:"updated" json:"updated,omitempty"` //pk
}

/* *
`
SELECT A.uuid, A.sequence, A.created,
       name, summary, method, args, result_filter,
       B.created AS updated, status, started, ended
  FROM service_step A
  LEFT JOIN service_step_status B
         ON A.uuid = B.uuid
        AND A.sequence = B.sequence
        AND B.created = (
            SELECT MAX(B.created) AS MAX_created
              FROM service_step_status B
             WHERE A.uuid = B.uuid
               AND A.sequence = B.sequence
            )
`
*/
func (record ServiceStep_tangled) TableName() string {
	q := `(
    SELECT %v /**columns**/
      FROM %v A /**service_step A**/
      LEFT JOIN %v B /**service_step_status B**/
             ON A.uuid = B.uuid 
            AND A.sequence = B.sequence
            AND B.created = (
                SELECT MAX(B.created) AS MAX_created 
                  FROM %v B /**service_step_status B**/
                 WHERE A.uuid = B.uuid 
                   AND A.sequence = B.sequence
                )
    ) X`

	columns := []string{
		"A.uuid",
		"A.sequence",
		"A.created",
		"B.created AS updated",
	}
	columns = append(columns, ServiceStep_essential{}.ColumnNames()...)
	columns = append(columns, ServiceStepStatus_essential{}.ColumnNames()...)
	A := record.ServiceStep.TableName()
	B := record.ServiceStepStatus_essential.TableName()
	return fmt.Sprintf(q, strings.Join(columns, ", "), A, B, B)
}
