package migrations

import (
	"bytes"
	_ "embed"
)

//go:embed sudory/latest
var sudoryLatest []byte

var SudoryLatest = new(Latest).SetReader(bytes.NewBuffer(sudoryLatest))
