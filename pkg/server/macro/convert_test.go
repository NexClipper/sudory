package macro_test

import (
	"fmt"
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
)

func TestConvtKeyValuePairs(t *testing.T) {

	m := make(map[string]string)

	for i := 0; i < 10; i++ {
		k := fmt.Sprintf("%c", rune('A'+i))
		v := fmt.Sprintf("%d", i)
		m[k] = v
	}

	s := ConvtKeyValuePairString(m)
	t.Log(s)
}

func TestConvtKeyValuePairJson(t *testing.T) {

	m := make(map[string]string)

	for i := 0; i < 10; i++ {
		k := fmt.Sprintf("%c", rune('A'+i))
		v := fmt.Sprintf("%d", i)
		m[k] = v
	}

	s := ConvtKeyValuePairJson(m)
	t.Log(s)
}
