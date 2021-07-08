package helpers

func Int64(v int64) *int64 {
	return &v
}

func Int(v int) *int {
	return &v
}

func Int32(v int32) *int32 {
	return &v
}

func String(v string) *string {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

func True() *bool {
	return Bool(true)
}

func False() *bool {
	return Bool(false)
}
