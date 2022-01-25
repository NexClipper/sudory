package macro

func ErrorWithHandler(err error, handler ...func(err error)) bool {
	if err != nil {
		for _, fn_h := range handler {
			fn_h(err)
		}
		return true
	}
	return false
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
