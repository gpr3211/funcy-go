package fp

// Exists takes in any array and returns true if predicate is true for any element.
//
//	-type Predicate func(A any) bool
func Exists[A any](input []A, pred Predicate[A]) bool {
	for _, element := range input {
		if pred(element) {
			return true
		}
	}
	return false
}

// ForAll returns true if all elements satisfy predicate.
func ForAll[A any](input []A, pred Predicate[A]) bool {
	for _, element := range input {
		if !pred(element) {
			return false
		}
	}
	return true
}
