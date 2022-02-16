package exceptions

import (
	"github.com/pkg/errors"
)

// Block {
// 	Try: func() {
// 		fmt.Println("Try..")
// 		Throw("stop it")
// 	},
// 	Catch: func(e Exception) {
// 		fmt.Printf("Caught %v\n", e)
// 	},
// 	Finally: func() {
// 		fmt.Println("Finally..")
// 	},
// }.Do()
// Block Try-Catch-Finally block struct
type Block struct {
	Try     func()
	Catch   func(error)
	Finally func()
}

// Do run block state
func (b Block) Do() {
	if b.Finally != nil {
		defer b.Finally()
	}

	if b.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				b.Catch(errors.Errorf("%+v", r))
			}
		}()
	}

	b.Try()
}

// Throw raise panic with exception
func Throw(ex Exception) {
	if ex != nil {
		panic(errors.Errorf("%+v", ex))
	}
}

// Exception pass exception to Catch
type Exception interface{}
