package newist

import (
	cryptov1 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v1"
)

func Cryptov1String(s string) *cryptov1.String {
	return func(s cryptov1.String) *cryptov1.String { return &s }(cryptov1.String(s))
}
