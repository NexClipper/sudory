package macro

import (
	"fmt"
)

type TestChapter struct {
	Subject string
	Action  func() error
}

type TestScenarios []TestChapter

func (scenarios TestScenarios) Foreach(fn func(string)) {
	for n, chapter := range scenarios {
		//if you wanna debugging in step? break here!!!
		err := chapter.Action()
		//^^^^^^^^^^
		ErrorHandle(err, func(err error) {
			fn(fmt.Sprintf("Chapter:'%d' Subject:'%s' Error: '%v'", n+1, chapter.Subject, err))
		})
	}
}
