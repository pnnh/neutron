package convert

import (
	"fmt"
	"time"

	"neutron/models"
)

func ConvertTime(value any) (time.Time, error) {

	if value == nil {
		return time.Time{}, models.ErrNilValue
	}

	switch v := value.(type) {
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to parse time from string: %w", err)
		}
		return t, nil
	case time.Time:
		return v, nil
	case int64:
		return time.Unix(v, 0), nil
	case int:
		return time.Unix(int64(v), 0), nil
	case uint64:
		return time.Unix(int64(v), 0), nil
	case uint:
		return time.Unix(int64(v), 0), nil
	case float64:
		return time.Unix(int64(v), 0), nil
	case float32:
		return time.Unix(int64(v), 0), nil
	case []byte:
		t, err := time.Parse(time.RFC3339, string(v))
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to parse time from byte slice: %w", err)
		}
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unsupported type %T for conversion to string", value)
}
