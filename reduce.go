package fp

// type reduceFunc applies func to successive elements.
type reduceFunc[A any] func(a1, a2 A) A

func Reduce[A any](input []A, reducer reduceFunc[A]) A {
	if len(input) == 0 {
		// return default zero
		return *new(A)
	}
	result := input[0]
	for _, element := range input[1:] {
		result = reducer(result, element)
	}
	return result
}

type Number interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint |
		~int8 | ~int16 | ~int32 | ~int64 | ~int |
		~float32 | ~float64
}

func absolute(n int) int {
	if n < 0 {
		out := -1
		n = n * out
	}
	return n
}

func Sum[A Number](input []A) A {
	return Reduce(input, func(a1, a2 A) A { return a1 + a2 })
}

func Product[A Number](input []A) A {
	return Reduce(input, func(a1, a2 A) A { return a1 * a2 })
}
