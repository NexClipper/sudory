//go:generate go run github.com/abice/go-enum --file=types.go --names --nocase=true
package v2

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
type ResultType int

/* ENUM(
	regist 		= 0
	send 		= 1
	processing	= 2
	success		= 4
	fail		= 8
)
*/
type StepStatus int
