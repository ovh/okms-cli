package utils

func PtrTo[T any](v T) *T {
	return &v
}

func DerefOrDefault[T any](v *T) T {
	if v == nil {
		var def T
		return def
	}
	return *v
}
