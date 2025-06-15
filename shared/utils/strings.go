package utils

import (
	"fmt"
	"strconv"
)

func ToString(val any) (string, error) {

	switch v := val.(type) {
	case []byte:
		return string(v), nil
	case int:
		return strconv.Itoa(v), nil
	case int8:
		return strconv.Itoa(int(v)), nil
	case int16:
		return strconv.Itoa(int(v)), nil
	case int32:
		return strconv.Itoa(int(v)), nil
	case int64:
		return strconv.Itoa(int(v)), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	default:
		res, ok := v.(string)

		if !ok {
			return "", fmt.Errorf("failed to cast string")
		}

		return res, nil
	}
}
