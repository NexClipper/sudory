package macro

import (
	"fmt"
	"testing"
)

type TestChapter struct {
	Subject string
	Method  func() error
}

type TestScenarios []TestChapter

func (scenarios TestScenarios) Foreach(t *testing.T) {
	for n, chapter := range scenarios {
		//if you wanna debugging in step? break here!!!
		err := chapter.Method()
		//^^^^^^^^^^
		ErrorHandle(err, func(err error) {
			t.Error(fmt.Sprintf("Chapter:'%d' Subject:'%s' Error: '%v'", n+1, chapter.Subject, err))
		})
	}
}
