package service

type ServiceType int

const (
	ServiceTypeAtOnce = iota + 1
	ServiceTypePeriodic
)

type ExecStep struct {
	Name      string `json:"name"`
	Sequence  uint64 `json:"sequence"`
	Command   string `json:"command"`
	Parameter string `json:"parameter"`
}

type Service struct {
	Id   string      `json:"id"`
	Type ServiceType `json:"service_type"`

	Name      string      `json:"name"`
	ClusterID uint64      `json:"cluster_id"`
	StepCount uint        `json:"step_count"`
	Step      []*ExecStep `json:"step"`
}

func (s *Service) Execute(updateChan chan *Service) error {
	// Execute Service Task

	// Update Service status
	if updateChan != nil {
		updateChan <- s
	}

	return nil
}
