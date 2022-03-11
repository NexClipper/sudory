package database

import "fmt"

func ErrorRecordWasNotFound() error {
	return fmt.Errorf("record was not found")
}
func ErrorNoAffected() error {
	return fmt.Errorf("no affected")
}
