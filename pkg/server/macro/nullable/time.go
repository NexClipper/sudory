package nullable

import (
	"time"
)

type nullTime struct {
	time.Time
	bool
}

func Time(v *time.Time) *nullTime {
	return new(nullTime).scan(v)
}

func (nullable nullTime) Value() time.Time {
	return nullable.Time
}
func (nullable nullTime) Ptr() *time.Time {
	return &nullable.Time
}
func (nullable nullTime) Has() bool {
	return nullable.bool
}

func (nullable *nullTime) scan(v interface{}) *nullTime {
	nullable.Time, nullable.bool = time.Time{}, false

	switch value := v.(type) {
	// case []byte:
	// 	nullable.Time, nullable.bool = string(value), true
	case time.Time:
		nullable.Time, nullable.bool = value, true
	case *time.Time:
		nullable.Time, nullable.bool = *value, true
		// case string:
		// 	if t, err := time.Parse(time.RFC3339, value); err != nil {
		// 		nullable.Time, nullable.bool = t, true
		// 		break
		// 	}
		// 	if t, err := time.Parse(time.RFC3339Nano, value); err != nil {
		// 		nullable.Time, nullable.bool = t, true
		// 		break
		// 	}
		// 	if t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", value); err != nil {
		// 		nullable.Time, nullable.bool = t, true
		// 		break
		// 	}
		// 	if t, err := time.Parse("2006-01-02 15:04:05", value); err != nil {
		// 		nullable.Time, nullable.bool = t, true
		// 		break
		// 	}
		// case *string:
		// 	if value != nil {
		// 		if t, err := time.Parse(time.RFC3339, *value); err != nil {
		// 			nullable.Time, nullable.bool = t, true
		// 			break
		// 		}
		// 		if t, err := time.Parse(time.RFC3339Nano, *value); err != nil {
		// 			nullable.Time, nullable.bool = t, true
		// 			break
		// 		}
		// 		if t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", *value); err != nil {
		// 			nullable.Time, nullable.bool = t, true
		// 			break
		// 		}
		// 		if t, err := time.Parse("2006-01-02 15:04:05", *value); err != nil {
		// 			nullable.Time, nullable.bool = t, true
		// 			break
		// 		}
		// 	}
		// case int:
		// 	nullable.Time, nullable.bool = strconv.FormatInt(int64(value), 10), true
		// case *int:
		// 	if value != nil {
		// 		nullable.Time, nullable.bool = strconv.FormatInt(int64(*value), 10), true
		// 	}
		// case int32:
		// 	nullable.Time, nullable.bool = strconv.FormatInt(int64(value), 10), true
		// case *int32:
		// 	if value != nil {
		// 		nullable.Time, nullable.bool = strconv.FormatInt(int64(*value), 10), true
		// 	}
		// case int64:
		// 	nullable.Time, nullable.bool = strconv.FormatInt(value, 10), true
		// case *int64:
		// 	if value != nil {
		// 		nullable.Time, nullable.bool = strconv.FormatInt(*value, 10), true
		// 	}
		// case uint:
		// 	nullable.Time, nullable.bool = strconv.FormatUint(uint64(value), 10), true
		// case *uint:
		// 	if value != nil {
		// 		nullable.Time, nullable.bool = strconv.FormatUint(uint64(*value), 10), true
		// 	}
		// case uint32:
		// 	nullable.Time, nullable.bool = strconv.FormatUint(uint64(value), 10), true
		// case *uint32:
		// 	if value != nil {
		// 		nullable.Time, nullable.bool = strconv.FormatUint(uint64(*value), 10), true
		// 	}
		// case uint64:
		// 	nullable.Time, nullable.bool = strconv.FormatUint(value, 10), true
		// case *uint64:
		// 	if value != nil {
		// 		nullable.Time, nullable.bool = strconv.FormatUint(*value, 10), true
		// 	}
	}

	return nullable
}
