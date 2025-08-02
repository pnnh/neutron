package convert

import "fmt"

func ConvertString(value any) (string, error) {

	if value == nil {
		return "", ErrNilValue
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case int:
	case int8:
	case int16:
	case int32:
	case int64:
	case uint:
	case uint8:
	case uint16:
	case uint32:
	case uint64:
		return fmt.Sprintf("%d", v), nil
	case []byte:
		return string(v), nil

	}
	return "", fmt.Errorf("unsupported type %T for conversion to string", value)
}
