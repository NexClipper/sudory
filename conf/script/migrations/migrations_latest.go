package migrations

import (
	"bufio"
	"io"
	"strings"
	"sync"
)

type Latest struct {
	r    io.Reader
	m    map[string]string
	once sync.Once
	err  error
}

func (la *Latest) SetReader(r io.Reader) *Latest {
	la.r = r
	return la
}

func (la *Latest) Version() string {
	const key = "version"
	if err := la.scan(); err != nil {
		return ""
	}

	return la.m[key]
}

func (la *Latest) scan() error {
	la.once.Do(func() {
		fs := bufio.NewScanner(la.r)
		fs.Split(bufio.ScanLines)

		la.m = map[string]string{}

		for fs.Scan() {
			l := fs.Text()
			l = strings.TrimSpace(l)

			l = strings.Split(l, "#")[0]
			if len(l) == 0 {
				continue
			}

			ll := strings.Split(l, "=")
			if len(ll) == 1 {
				continue
			}

			k, v := ll[0], ll[1]
			la.m[k] = v
		}

	})
	return la.err
}

func (l *Latest) Err() error {
	l.scan()
	return l.err
}
