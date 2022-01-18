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
