package managed_channel

import (
	"context"
	"encoding/json"
	"sync"

	channelv2 "github.com/NexClipper/sudory/pkg/server/model/channel/v2"
	"github.com/itchyny/gojq"
)

type Formatter interface {
	Format(map[string]interface{}) (interface{}, error)
	Type() channelv2.FormatType
}

type Formatter_fields struct {
	FormatData string

	err error
	sync.Once
	field_names []string
}

func (format Formatter_fields) Type() channelv2.FormatType {
	return channelv2.FormatTypeFields
}

func (format *Formatter_fields) Format(a map[string]interface{}) (interface{}, error) {
	rst := map[string]interface{}{}

	format.Do(func() {
		format.err = json.Unmarshal([]byte(format.FormatData), &format.field_names)
	})
	if format.err != nil {
		return nil, format.err
	}

	// make table
	field_names := map[string]struct{}{}
	for _, field_name := range format.field_names {
		field_names[field_name] = struct{}{}
	}

	for k, v := range a {
		ok := func(key string) bool {
			for field_name := range field_names {
				if key == field_name {
					return true
				}
			}
			return false
		}(k)
		if ok {
			rst[k] = v
			// reduce table
			delete(field_names, k)
		}
	}

	return rst, nil
}

type Formatter_jq struct {
	FormatData string

	err error
	sync.Once
	query *gojq.Query
}

func (format Formatter_jq) Type() channelv2.FormatType {
	return channelv2.FormatTypeJq
}

func (format *Formatter_jq) Format(a map[string]interface{}) (interface{}, error) {
	rst := make([]interface{}, 0, len(a))

	format.Do(func() {
		format.query, format.err = gojq.Parse(format.FormatData)
	})

	if format.err != nil {
		return nil, format.err
	}

	var err error
	iter := format.query.RunWithContext(context.TODO(), a)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		switch v := v.(type) {
		case error:
			err = v
		case map[string]interface{}:
			rst = append(rst, v)
		default:
			rst = append(rst, v)
		}
	}
	if err != nil {
		return nil, err
	}

	if len(rst) == 1 {
		return rst[0], err
	}

	return rst, err
}
