package prepare

import "fmt"

func ErrorInvalidArgumentEmptyString() error {
	return fmt.Errorf("invalid argument: empty string")
}

func ErrorNotFoundHandler(key string) error {
	return fmt.Errorf("not found handler key=%s", key)
}

func ErrorUnsupportedType(value interface{}) error {
	return fmt.Errorf("unsupported type value=%#v", value)
}
