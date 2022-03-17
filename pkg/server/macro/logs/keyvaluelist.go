package logs

import (
	"bytes"

	"github.com/NexClipper/sudory/pkg/server/macro/logs/internal/serialize"
)

// KVL
//  KeyValueList
func KVL(keysAndValues ...interface{}) string {
	return parseKVList(keysAndValues...)
}

func parseKVList(keysAndValues ...interface{}) string {

	buf := bytes.Buffer{}

	serialize.KVListFormat(&buf, keysAndValues...)

	return buf.String()
}
