package cast

func Ptr[T any](arg T) *T {
	return &arg
}
