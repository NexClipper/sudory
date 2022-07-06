package config

import (
	"bufio"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type Latest struct {
	Source string
	m      map[string]string
	once   sync.Once
	err    error
}

func (l Latest) Filename() string {
	return path.Join(l.Source, "latest")
}

func (l *Latest) Version() string {
	const key = "version"
	if err := l.scan(); err != nil {
		return ""
	}

	return l.m[key]
}

func (la *Latest) scan() error {
	la.once.Do(func() {
		fd, err := os.Open(la.Filename())
		err = errors.Wrapf(err, "can not open file in migrations")
		if err != nil {
			la.err = err
			return
		}

		fs := bufio.NewScanner(fd)
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
