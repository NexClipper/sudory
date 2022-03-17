package filepool

import (
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/pkg/errors"
)

type FilePool struct {
	pool map[string]*os.File
	sync.Mutex
}

//OpenFile
func (handler *FilePool) OpenFile(filename string) (*os.File, error) {
	//lock guard
	handler.Lock()
	defer handler.Unlock()

	if handler.pool == nil {
		handler.pool = make(map[string]*os.File)
	}

	if fd, ok := handler.pool[filename]; ok {
		return fd, nil
	}

	//directory check
	dir := path.Dir(filename)
	if _, err := os.Stat(dir); err != nil {
		//if path is not exists
		if os.IsNotExist(err) {
			//then make directory
			if err := os.MkdirAll(dir, 0700); err != nil {
				return nil, errors.Wrapf(err, "os.MkdirAll path=%s perm=0700", dir) //didn't make mkdir
			}
		} else {
			//handle for other
			return nil, errors.Wrapf(err, "os.Stat path=%s", dir) //didn't make mkdir
		}
	}

	//open file
	fd, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, errors.Wrapf(err, "os.OpenFile name=%s flag=os.O_APPEND|os.O_WRONLY|os.O_CREATE perm=0600", filename)
	}

	handler.pool[filename] = fd

	return fd, nil
}

//File
func (handler *FilePool) Write(filename string, b []byte) error {
	//lock guard
	handler.Lock()
	defer handler.Unlock()

	fd, ok := handler.pool[filename]
	if !ok {
		return fmt.Errorf("not found file descriptor from table filename=%s", filename)
	}

	if _, err := fd.Write(b); err != nil {
		hs := hex.EncodeToString(b)
		return errors.Wrapf(err, "fd.Write b=%s", hs)
	}

	return nil
}

//File
func (handler *FilePool) File(filename string) *os.File {
	//lock guard
	handler.Lock()
	defer handler.Unlock()

	return handler.pool[filename]
}

//CloseAll
func (handler *FilePool) CloseAll() {
	handler.Lock()
	defer handler.Unlock()

	for filename := range handler.pool {
		//sync
		handler.pool[filename].Sync()
		//file close
		handler.pool[filename].Close()
		//remove item
		delete(handler.pool, filename)
	}

}

//Close
func (handler *FilePool) Close(filename string) (err error) {
	handler.Lock()
	defer handler.Unlock()

	if _, ok := handler.pool[filename]; ok {
		//sync
		handler.pool[filename].Sync()
		//file close
		err = handler.pool[filename].Close()
		//remove item
		delete(handler.pool, filename)
	}

	return
}

var inst = FilePool{}

func OpenFile(filename string) (*os.File, error) {
	return inst.OpenFile(filename)
}

func Write(filename string, b []byte) error {
	return inst.Write(filename, b)
}

func Close(filename string) (err error) {
	return inst.Close(filename)
}

func CloseAll() {
	inst.CloseAll()
}
