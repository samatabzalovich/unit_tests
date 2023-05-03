package pointer

func ToString(s string) *string {
	return &s
}

func ToInt32(i int32) *int32 {
	return &i
}

func ToFloat64(f float64) *float64 {
	return &f
}

func ToBool(b bool) *bool {
	return &b
}
