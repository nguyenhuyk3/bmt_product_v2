package convertors

import "fmt"

func ConvertInterfaceToSlice(input interface{}) ([]string, error) {
	if input == nil {
		return nil, nil
	}

	ifaceSlice, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("input is not []interface{}")
	}

	result := make([]string, 0, len(ifaceSlice))
	for _, v := range ifaceSlice {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("input element is not string")
		}

		result = append(result, s)
	}

	return result, nil
}
