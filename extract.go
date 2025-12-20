package botc

func convertStringSlice(x any) ([]string, bool) {
	xs := x.([]any)
	r := make([]string, len(xs))
	for i, s := range xs {
		S, ok := s.(string)
		if !ok {
			return r, false
		}
		r[i] = S
	}
	return r, true
}

func extractSlice(k string, m map[string]any) ([]any, bool, error) {
	v, ok := m[k]
	if !ok {
		return []any{}, false, nil
	}
	values, ok := v.([]any)
	if !ok {
		return values, true, NewConversionError(k, v)
	}
	return values, true, nil
}

func extractString(k string, m map[string]any) (string, bool, error) {
	v, ok := m[k]
	if !ok {
		return "", false, nil
	}
	value, ok := v.(string)
	if !ok {
		return "", true, NewConversionError(k, v)
	}
	return value, true, nil
}

func extractRequiredString(k string, m map[string]any) (string, error) {
	value, ok, err := extractString(k, m)
	if !ok {
		return value, NewRequiredFieldMissingError(k)
	}
	if err != nil {
		return value, err
	}
	return value, nil
}

func extractStringSlice(k string, m map[string]any) ([]string, bool, error) {
	v, ok := m[k]
	if !ok {
		return []string{}, false, nil
	}
	if v == nil {
		return []string{}, true, nil
	}
	values, ok := convertStringSlice(v)
	if !ok {
		return []string{}, true, NewConversionError(k, v)
	}
	return values, true, nil
}

func extractInt(k string, m map[string]any) (int, bool, error) {
	v, ok := m[k]
	if !ok {
		return 0, false, nil
	}
	value64, ok := v.(float64)
	if !ok {
		return 0, true, NewConversionError(k, v)
	}
	value := int(value64)
	return value, true, nil
}

func extractRoleRequiredInt(k string, m map[string]any) (int, error) {
	value, ok, err := extractInt(k, m)
	if !ok {
		return 0, NewRequiredFieldMissingError(k)
	}
	if err != nil {
		return 0, err
	}
	return value, nil
}

func extractBool(k string, m map[string]any) (bool, bool, error) {
	v, ok := m[k]
	if !ok {
		return false, false, nil
	}
	value, ok := v.(bool)
	if !ok {
		return false, true, NewConversionError(k, v)
	}
	return value, true, nil
}

func extractRequiredBool(k string, m map[string]any) (bool, error) {
	value, ok, err := extractBool(k, m)
	if !ok {
		return false, NewRequiredFieldMissingError(k)
	}
	if err != nil {
		return false, err
	}
	return value, nil
}
