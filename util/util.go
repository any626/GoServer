package util

// Check panics if err is not nil
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
