//go:generate go run github.com/abice/go-enum --file=polling_option.go --names --nocase=true
package v1

import (
	"time"
)

/* ENUM(
regular
smart
)
*/
type PollingType int32

type PollingHandler interface {
	Interval(_default time.Duration, service_count int) time.Duration
	ToMap() map[string]interface{}
}

func (opt ClusterProperty) GetPollingOption() PollingHandler {
	pollingType := PollingTypeRegular

	if opt_type, ok := opt.PollingOption["type"]; ok {
		opt_type, _ := opt_type.(string)
		if opt_type, err := ParsePollingType(opt_type); err == nil {
			pollingType = opt_type
		}
	}

	switch pollingType {
	case PollingTypeSmart:
		idle := 0
		busy := 0
		if opt, ok := opt.PollingOption["idle"]; ok {
			opt, _ := opt.(float64) //json 숫자타입은 부동소수
			idle = int(opt)
		}
		if opt, ok := opt.PollingOption["busy"]; ok {
			opt, _ := opt.(float64) //json 숫자타입은 부동소수
			busy = int(opt)
		}
		return &SmartPollingOption{IdleInterval: idle, BusyInterval: busy}
	case PollingTypeRegular:
		fallthrough
	default:
		return &RagulerPollingOption{}
	}
}

func (opt *ClusterProperty) SetPollingOption(handle PollingHandler) {
	switch handle := handle.(type) {
	case *SmartPollingOption:
		opt.PollingOption = handle.ToMap()
	case *RagulerPollingOption:
		opt.PollingOption = handle.ToMap()
	}
}

type RagulerPollingOption struct{}

func (opt RagulerPollingOption) Interval(_default time.Duration, service_count int) time.Duration {
	return _default
}
func (opt RagulerPollingOption) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"type": PollingTypeRegular.String(),
	}
	return m
}

type SmartPollingOption struct {
	IdleInterval int `json:"idle,omitempty"` //(초)
	BusyInterval int `json:"busy,omitempty"` //(초)
}

func (opt SmartPollingOption) Interval(_default time.Duration, service_count int) time.Duration {
	busy := func() time.Duration {
		if opt.BusyInterval == 0 {
			return _default
		}
		return time.Duration(opt.BusyInterval) * time.Second
	}
	idle := func() time.Duration {
		if opt.IdleInterval == 0 {
			return _default
		}
		return time.Duration(opt.IdleInterval) * time.Second
	}

	if 0 < service_count {
		return busy()
	} else {
		return idle()
	}
}

func (opt SmartPollingOption) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"type": PollingTypeSmart.String(),
		"idle": opt.IdleInterval,
		"busy": opt.BusyInterval,
	}
	return m
}
