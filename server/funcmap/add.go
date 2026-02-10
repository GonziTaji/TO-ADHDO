package funcmap

// Simple add
func add(values ...int) int {
	r := 0
	for _, v := range values {
		r = r + v
	}
	return r
}
