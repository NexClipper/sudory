//go:generate go run github.com/abice/go-enum --file=types.go --names --nocase=true
package service

import "github.com/pkg/errors"

/* ENUM(
	none
	remove
)
*/
type OnCompletion int8

/* ENUM(
	none
	database
	DigitalOcean:Spaces
)
*/
type ResultSaveType int

/* ENUM(
	regist 		= 0
	send 		= 1
	processing	= 2
	success		= 4
	fail		= 8
)
*/
type StepStatus int

func (status StepStatus) Valid() error {
	status_ := status.String()
	for _, iter := range StepStatusNames() {
		if status_ == iter {
			return nil
		}
	}
	return errors.Errorf("invalid %v", status)
}

/* ENUM(
	low
	middle
	high
)
*/
type Priority int
