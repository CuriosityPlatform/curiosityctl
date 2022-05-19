package types

// EmptyToPtr Check if passed value is zero for current type and returns pointer otherwise returns nil ptr
// !Not efficient use since allocates in heap type T to check of zero value
func EmptyToPtr[T comparable](t T) *T {
	t2 := new(T)
	if t == *t2 {
		return nil
	}
	return &t
}
