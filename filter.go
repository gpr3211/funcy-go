package fp

// Predicate type takes any type of slice  and a function which returns a bool
type Predicate[A any] func(A) bool

// Filter takes as input a slice of Any and a Predicate type function to be applied to each item in the slice.
// returns a new slice containing only the satisfied elements
func Filter[A any](input []A, pred Predicate[A]) []A {
	output := []A{}
	for _, element := range input {
		if pred(element) {
			output = append(output, element)
		}
	}
	return output
}
