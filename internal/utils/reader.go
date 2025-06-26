package utils

type Preticate[T any] func(v T) bool

type Reader[T any] struct {
	S   []T
	idx int
}

func NewReader[T any](s []T) *Reader[T] {
	return &Reader[T]{S: s}
}

func (t *Reader[T]) Peek() (v T, ok bool) {
	if t.idx >= len(t.S) {
		return v, false
	}

	return t.S[t.idx], true
}

func (t *Reader[T]) Take() (v T, ok bool) {
	v, ok = t.Peek()
	if !ok {
		return v, false
	}

	t.idx++

	return v, true
}

func (t *Reader[T]) Back() {
	t.idx--
}

func (t *Reader[T]) SkipWhile(p Preticate[T]) bool {
	for {
		v, ok := t.Take()
		if !ok {
			return false
		}
		if !p(v) {
			t.Back()
			return true
		}
	}
}
