package batch

func check(err interface{}, f func(string, ...interface{}), msg string, args ...interface{}) {
	if err != nil {
		f(msg, args...)
	}
}
