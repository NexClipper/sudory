package database

import "fmt"

func ErrorRecordWasNotFound() error {
	return fmt.Errorf("record was not found")
}
func ErrorNoAffecte() error {
	return fmt.Errorf("no affecte")
}
