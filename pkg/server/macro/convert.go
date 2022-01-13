package macro

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func ConvtKeyValuePairString(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func ConvtKeyValuePairJson(m map[string]string) string {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(b)
}
