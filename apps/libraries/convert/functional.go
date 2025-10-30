package convert

func MapSlice[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))

	for i, v := range slice {
		result[i] = fn(v)
	}

	return result
}

func UnwrapOr[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}

	return *value
}

func ValueToPointer[T any](value T) *T {
	return &value
}

func FromIntToInt64[T number](value *T) *int64 {
	if value == nil {
		return nil
	}

	int64Value := int64(*value)

	return &int64Value
}

func PointerToValue[T any](value *T) T {
	if value == nil {
		var defaultValue T
		return defaultValue
	}

	return *value
}
