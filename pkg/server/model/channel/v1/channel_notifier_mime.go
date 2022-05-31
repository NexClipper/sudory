package v1

import "github.com/pkg/errors"

//go:generate go run github.com/abice/go-enum --file=channel_notifier_mime.go --names --nocase=true

/* ENUM(
application/json
application/xml
)
*/
type ContentType int

//database meta info
type MIME struct {
	ContentType string `json:"Content-Type"              xorm:"'content_type'              notnull"`
}

func (mime MIME) Valid() error {
	if _, err := ParseContentType(mime.ContentType); err != nil {
		return errors.Wrapf(err, "valid Content-Type")
	}

	return nil
}
