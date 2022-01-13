package macro

import "time"

func TimeParse(s string) time.Time {
	const layout = "2006-01-02 15:04:05 MST"
	t, err := time.Parse(layout, s)
	if err != nil {
		panic(err)
	}

	return t
}
