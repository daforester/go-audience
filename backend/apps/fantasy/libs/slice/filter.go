package slice

func Filter[T any](s *[]T, f func(any) bool) {
	n := 0
	for _, v := range *s {
		if f(v) {
			(*s)[n] = v
			n = n + 1
		}
	}
	*s = (*s)[:n]
	return
}
