package v1

import (
	"strings"

	"github.com/google/uuid"
)

func NewUuidString() string {
	u := uuid.NewString()
	return strings.Replace(u, "-", "", -1)

}
