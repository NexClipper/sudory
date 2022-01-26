package events

import (
	"os"
	"path"
	"sync"
)

var Files FileHandler

type FileHandler struct {
	table map[string]*os.File
	mux   sync.Mutex
}

func (fh *FileHandler) OpenFile(filename string) error {
	fh.mux.Lock()
	defer fh.mux.Unlock()

	if fh.table == nil {
		fh.table = make(map[string]*os.File)
	}

	if _, ok := fh.table[filename]; ok {
		return nil
	}

	dir := path.Dir(filename)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err //didn't make mkdir
		}
	}

	fd, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	fh.table[filename] = fd

	return nil
}

func (fh *FileHandler) GetFile(filename string) *os.File {
	fh.mux.Lock()
	defer fh.mux.Unlock()

	return fh.table[filename]
}

func (fh *FileHandler) WriteFile(filename string, b ...[]byte) (int, error) {
	fh.mux.Lock()
	defer fh.mux.Unlock()

	var cnt int
	for _, it := range b {
		n, err := fh.table[filename].Write(it)
		if err != nil {
			return cnt, err
		}
		cnt += n
	}

	return cnt, nil
}

func (fh *FileHandler) CloseFileAll() {
	fh.mux.Lock()
	defer fh.mux.Unlock()

	//close
	for filename := range fh.table {
		fh.table[filename].Close()
	}
	//remove
	for filename := range fh.table {
		delete(fh.table, filename)
	}
}

func (fh *FileHandler) CloseFile(filename string) {
	fh.mux.Lock()
	defer fh.mux.Unlock()

	//close
	fh.table[filename].Close()

	//remove
	delete(fh.table, filename)

}
