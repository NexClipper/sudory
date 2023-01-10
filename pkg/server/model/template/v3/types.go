//go:generate go run github.com/abice/go-enum --file=types.go --names --nocase=true
package template

/* ENUM(
	none
	predefined
	system
)
*/
type Origin int
