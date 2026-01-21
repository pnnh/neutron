package convert

import "github.com/pnnh/neutron/models"

func ToFloat64(value any) (float64, error) {
	if value == nil {
		return 0, models.ErrNilValue
	}
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, nil
	}
}

func ToFloat32(value any) (float32, error) {
	if value == nil {
		return 0, models.ErrNilValue
	}
	switch v := value.(type) {
	case int:
		return float32(v), nil
	case int8:
		return float32(v), nil
	case uint8:
		return float32(v), nil
	case int16:
		return float32(v), nil
	case uint16:
		return float32(v), nil
	case int32:
		return float32(v), nil
	case uint32:
		return float32(v), nil
	case int64:
		return float32(v), nil
	case uint64:
		return float32(v), nil
	case float32:
		return v, nil
	case float64:
		return float32(v), nil
	default:
		return 0, nil
	}
}
