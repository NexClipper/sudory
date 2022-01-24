package events

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"sync"
)

type FileListenOption struct {
	Type    string `yaml:"type,omitempty"`
	Name    string `yaml:"name,omitempty"`
	Pattern string `yaml:"pattern,omitempty"`
	Option  struct {
		Path string `yaml:"path,omitempty"`
	} `yaml:"option,omitempty"`
}
type FileListener struct {
	opt FileListenOption
	mux sync.Mutex
}

//check implementation
var _ ListenerContext = (*FileListener)(nil)

func NewFileListener(opt FileListenOption) ListenerContext {
	dir := path.Dir(opt.Option.Path)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			panic(err) //didn't make mkdir
		}
	}
	return &FileListener{opt: opt}
}
func (me *FileListener) Type() string {
	return "file"
}
func (me *FileListener) Name() string {
	return me.opt.Name
}
func (me *FileListener) Pattern() string {
	return me.opt.Pattern
}
func (me *FileListener) Dest() string {
	return me.opt.Option.Path
}

func (me *FileListener) Raise(v interface{}) error {
	me.mux.Lock()
	defer me.mux.Unlock()
	return me.raise(v)
}

func (me *FileListener) raise(v interface{}) error {

	buffer, err := serialize_json(v)
	if err != nil {
		return err
	}

	fd, err := os.OpenFile(me.opt.Option.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	_, err = fd.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	return nil
}

func serialize_json(v interface{}) (*bytes.Buffer, error) {

	if v == nil {
		return nil, nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}
