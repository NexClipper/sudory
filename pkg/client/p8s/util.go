package p8s

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/prometheus/common/model"
)

func mapToStruct(m map[string]interface{}, o interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, o)
}

type Duration time.Duration

func (d Duration) String() string {
	return time.Duration(d).String()
}

func parseDuration(s string) (Duration, error) {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		sec := f * float64(time.Second)
		if sec > float64(math.MaxInt64) || sec < float64(math.MinInt64) {
			return 0, fmt.Errorf("failed to parse %s to duration(out of int64 range)", s)
		}

		return Duration(sec), nil
	}

	if d, err := model.ParseDuration(s); err == nil {
		return Duration(d), nil
	}

	return 0, fmt.Errorf("failed to parse %s to duration", s)
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	var err error
	*d, err = parseDuration(s)
	return err
}

type apiResponse struct {
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data"`
	ErrorType string          `json:"errorType"`
	Error     string          `json:"error"`
	Warnings  []string        `json:"warnings,omitempty"`
}

func formatTime(t time.Time) string {
	return strconv.FormatFloat(float64(t.Unix())+float64(t.Nanosecond())/1e9, 'f', -1, 64)
}
