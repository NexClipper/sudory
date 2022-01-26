package events

import (
	"fmt"
)

func ErrorInvaildListenerType(v interface{}) error {
	return fmt.Errorf("invalid event listener type value=%v", v)
}
func ErrorNotFoundListenerType() error {
	return fmt.Errorf("not found event listener type")
}
func ErrorInvaildListenerTypeZeroLength() error {
	return fmt.Errorf("invalid event listener type (zero-length)")
}
func ErrorUndifinedEventListener(t string) error {
	return fmt.Errorf("undefined event listener type='%s'", t)
}

func ErrorInvaildListenerConfig() error {
	return fmt.Errorf("invalid listener config")
}

func ErrorInvalidEventPatternZeroLength() error {
	return fmt.Errorf("invalid event pattern (zero-length)")
}
func ErrorRegexComplileEventPattern(pattern string, err error) error {
	return fmt.Errorf("failed to regexp compile event pattern='%s' error='%w'", pattern, err)
}

func ErrorInvalidEventNameZeroLength() error {
	return fmt.Errorf("invalid event name (zero-length)")
}
