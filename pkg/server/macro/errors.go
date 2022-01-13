package macro

func ErrorHandle(err error, handles ...func(err error)) bool {
	if err != nil {
		for _, h := range handles {
			h(err)
		}
		return true
	}
	return false
}

func HasError(err error) bool {
	return err != nil
}
