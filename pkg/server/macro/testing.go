package macro

import (
	"fmt"
)

type TestChapter struct {
	Subject string
	Action  func() error
}

type TestScenarios []TestChapter

func (scenarios TestScenarios) Foreach(fn func(error)) {
	for n, chapter := range scenarios {
		//if you wanna debugging in step? break here!!!
		err := chapter.Action()
		//^^^^^^^^^^
		ErrorHandle(err, func(err error) {
			fn(fmt.Errorf("chapter:'%d' subject:'%s' error: '%w'", n+1, chapter.Subject, err))
		})
	}
}
