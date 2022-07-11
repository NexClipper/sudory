package repo

import "errors"

var errExistRepo = errors.New("repository already exists")

func ErrIsExistRepo(err error) bool {
	return err == errExistRepo
}
