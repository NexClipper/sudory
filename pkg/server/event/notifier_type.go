//go:generate go run github.com/abice/go-enum --file=notifier_type.go --names --nocase=true
package event

/* ENUM(
console
webhook
file
rabbitMQ
)
*/
type NotifierType int32
