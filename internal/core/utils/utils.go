package utils

// Pointer pointer
func Pointer[Value any](v Value) *Value {
	return &v
}
