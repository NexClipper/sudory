package slicestrings

func Foreach(e []string, iterator func(s string) bool) bool {
	for i := range e {
		if !iterator(e[i]) {
			return false
		}
	}
	return true
}

func Map(e []string, mapper func(s string, i int) string) []string {
	// clone
	c := make([]string, len(e))
	copy(c, e)

	for i := range c {
		c[i] = mapper(c[i], i)
	}
	return c
}

func Aggregate(e []string, init string, iterator func(a string, b string, i int) string) string {
	for i := range e {
		init = iterator(init, e[i], i)
	}
	return init
}
