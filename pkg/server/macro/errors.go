package macro

func ErrorWithHandler(err error, handler ...func(err error)) bool {

	if err == nil {
		return false
	}

	for n := range handler {
		func(fn func(err error)) {
			defer func() {
				_ = recover()
			}()
			fn(err)
		}(handler[n])
	}
	return true
}

func HasError(err error) bool {
	return err != nil
}

func Eqaul(a, b error) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Error() == b.Error()
}

func NotEqaul(a, b error) bool {
	return !Eqaul(a, b)
}
