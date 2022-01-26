package events

import (
	"bytes"
	"encoding/json"
)

type FileListenOption struct {
	Type string `yaml:"type,omitempty"`
	Path string `yaml:"path,omitempty"`
}
type FileListener struct {
	config EventConfig
	opt    FileListenOption
}

//check implementation
var _ ListenerContexter = (*FileListener)(nil)

func NewFileListener(config EventConfig, opt FileListenOption) (ListenerContexter, error) {

	err := Files.OpenFile(opt.Path)
	if err != nil {
		return nil, err
	}

	return &FileListener{config: config, opt: opt}, nil
}
func (me FileListener) Type() string {
	return ListenerTypeFile.String()
}

func (me FileListener) Name() string {
	return me.config.Name
}

func (me FileListener) Summary() string {
	return me.opt.Path
}

func (ctx FileListener) Raise(v interface{}) error {
	return file_raise(ctx.opt, v)
}

func file_raise(opt FileListenOption, v interface{}) error {

	buffer, err := serialize_json(v)
	if err != nil {
		return err
	}

	_, err = Files.WriteFile(opt.Path,
		buffer.Bytes(),
		[]byte{'\n'},
	)
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
