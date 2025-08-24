package convert

func ConvertBool(value any) (bool, error) {
	if boolValue, ok := value.(bool); ok {
		return boolValue, nil
	}
	return false, nil

}
