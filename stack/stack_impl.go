package Stack

type Stack[T interface{}] []T

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack[T]) Push(n T) {
	*s = append(*s, n)
}

func (s *Stack[T]) Top() (T, bool) {
	if s.IsEmpty() {
		var t T
		return t, false
	}
	return (*s)[len(*s)-1], true
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.IsEmpty() {
		var t T
		return t, false
	}
	el := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return el, true
}
