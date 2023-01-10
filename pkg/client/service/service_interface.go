package service

import "time"

type Version string

const (
	// client service version
	SERVICE_VERSION_V1 Version = "v1"
	SERVICE_VERSION_V2 Version = "v2"
)

func (v Version) GetMatchServerServiceVersion() string {
	switch v {
	case SERVICE_VERSION_V1:
		return "v3"
	case SERVICE_VERSION_V2:
		return "v4"
	}

	return ""
}

type StepStatus int32

const (
	StepStatusPreparing = iota + 1
	StepStatusProcessing
	StepStatusSuccess
	StepStatusFail
)

type ServiceInterface interface {
	Version() Version
	GetId() string
	GetPriority() int
	GetCreatedTime() time.Time
}

type ServiceUpdateInterface interface {
	Version() Version
	GetId() string
	GetStepCount() int
	GetSequence() int
	GetStatus() StepStatus
}
