package service

//go:generate go run -mod=mod github.com/abice/go-enum --file=types.go --names --nocase=true

import "github.com/pkg/errors"

/* ENUM(
	none
	database
	DigitalOcean:Spaces
)
*/
type ResultSaveType int

/* ENUM(
	regist 		= 0
	sent 		= 1
	processing	= 2
	succeeded   = 4
	failed		= 8
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
