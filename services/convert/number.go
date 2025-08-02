package convert

import "fmt"

func ConvertInt(value any) (int, error) {
	if value == nil {
		return 0, ErrNilValue
	}
	switch v := value.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case uint8:
		return int(v), nil
	case int16:
		return int(v), nil
	case uint16:
		return int(v), nil
	case int32:
		return int(v), nil
	case uint32:
		return int(v), nil
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("ConvertInt error, negative value: %d", v)
		}
		if v > int64(^uint(0)>>1) {
			return 0, fmt.Errorf("ConvertInt error, value too large: %d", v)
		}
		return int(v), nil
	case uint64:
		if v > uint64(^uint(0)>>1) {
			return 0, fmt.Errorf("ConvertInt error, value too large: %d", v)
		}
		return int(v), nil
	case float32:
		if v < 0 {
			return 0, fmt.Errorf("ConvertInt error, negative value: %f", v)
		}
		if v > float32(int(^uint(0)>>1)) {
			return 0, fmt.Errorf("ConvertInt error, value too large: %f", v)
		}
		return int(v), nil
	case float64:
		if v < 0 {
			return 0, fmt.Errorf("ConvertInt error, negative value: %f", v)
		}
		if v > float64(int(^uint(0)>>1)) {
			return 0, fmt.Errorf("ConvertInt error, value too large: %f", v)
		}
		return int(v), nil
	case string:
		var intValue int
		_, err := fmt.Sscanf(v, "%d", &intValue)
		if err != nil {
			return 0, fmt.Errorf("ConvertInt error, cannot convert string to int: %s", v)
		}
		return intValue, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("ConvertInt error, unsupported type: %T", value)
}
