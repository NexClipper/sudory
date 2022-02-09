package macro

import (
	"strings"

	"github.com/google/uuid"
)

func NewUuidString() string {
	u := uuid.NewString()
	return strings.Replace(u, "-", "", -1)

}

func EmptyUuid(s *string) {
	if len(*s) == 0 {
		*s = NewUuidString()
	}
}
