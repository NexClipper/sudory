package macro

func Do(err *error, fn func() (err error)) {
	if err == nil {
		panic(*err) // panic
	}
	if *err != nil {
		return
	}
	*err = fn()
}
