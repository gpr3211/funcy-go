package fp

// MapFunc type of Any func(A any) A
type MapFunc[A any] func(A) A

// Map applies MapFunc to each element of the input slice
// returns a new modified slice with the same number of elements
func Map[A any](input []A, m MapFunc[A]) ([]A, error) {
	output := make([]A, len(input))
	for i, element := range input {
		output[i] = m(element)
	}
	return output, nil
}

// FlatMap is a one-to-many map. Example below
// - input slice of any type
// - func(A) []A to be applied to each element
func FlatMap[A any](input []A, m func(A) []A) []A {
	output := []A{}
	for _, element := range input {
		newlemenets := m(element)
		output = append(output, newlemenets...)
	}
	return output
}
