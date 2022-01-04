package model

import "time"

type ReqService struct {
	Name      string     `json:"name"`
	ClusterID uint64     `json:"cluster_id"`
	StepCount uint       `json:"step_count"`
	Step      []*ReqStep `json:"step"`
}

type ReqStep struct {
	Name      string `json:"name"`
	Sequence  uint64 `json:"sequence"`
	Command   string `json:"command"`
	Parameter string `json:"parameter"`
}

type Service struct {
	ID           uint64    `xorm:"pk autoincr 'id'"`
	Name         string    `xorm:"name"`
	ClusterID    uint64    `xorm:"cluster_id"`
	StepCount    uint      `xorm:"step_count"`
	StepPosition uint      `xorm:"step_position"`
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
}

func (m *Service) GetType() string {
	return "SERVICE"
}

type Step struct {
	ID        uint64 `xorm:"pk autoincr 'id'"`
	Name      string `xorm:"name"`
	ServiceID uint64 `xorm:"service_id"`
	Sequence  uint64 `xorm:"sequence"`
	Command   string `xorm:"command"`
	Parameter string `xorm:"parameter"`
}

func (m *Step) GetType() string {
	return "STEP"
}

type ServiceStep struct {
	Service `xorm:"extends"`
	Step    `xorm:"extends"`
}

type ReqClientGetService struct {
	ClusterID uint64 `json:"cluster_id"`
}

type RespService struct {
	Name      string      `json:"name"`
	ClusterID uint64      `json:"cluster_id"`
	StepCount uint        `json:"step_count"`
	Step      []*RespStep `json:"step"`
}

func (m *RespService) GetType() string {
	return "RESPSERVICE"
}

type RespStep struct {
	Name      string `json:"name"`
	Sequence  uint64 `json:"sequence"`
	Command   string `json:"command"`
	Parameter string `json:"parameter"`
}
