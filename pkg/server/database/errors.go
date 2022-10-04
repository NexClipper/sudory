package database

import "fmt"

// func ErrorRecordWasNotFound() error {
// 	return fmt.Errorf("record was not found")
// }
// func ErrorNoAffected() error {
// 	return fmt.Errorf("no affected")
// }

var (
	ErrorFailedToCheckRecord = fmt.Errorf("failed to check record")
	ErrorRecordWasNotFound   = fmt.Errorf("record was not found")
	ErrorNoAffected          = fmt.Errorf("no affected")
	ErrorNoLastInsertId      = fmt.Errorf("no last insert id")
)
