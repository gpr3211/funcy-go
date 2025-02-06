package fp

// Set function that works with any slice (basic types or structs)
// returns nil on empty input.
func Set[T comparable](slice []T) []T {
	uniqueMap := make(map[T]bool)
	var uniqueSlice []T
	for _, item := range slice {
		if !uniqueMap[item] {
			uniqueMap[item] = true
			uniqueSlice = append(uniqueSlice, item)
		}
	}
	return uniqueSlice
}

// Union takes 2 sets of comparable types and returns their union.
// returns nil else
func Union[T comparable](set1, set2 []T) []T {
	combined := append(set1, set2...)
	return Set(combined)
}

// Intersection returns a slice containing the intersection of two slices.
// returns nil if there is no intersection between set1 and set2
func Intersection[T comparable](set1, set2 []T) []T {
	uniqueMap := make(map[T]bool)
	for _, item := range set1 {
		uniqueMap[item] = true
	}
	var intersectionSlice []T
	for _, item := range set2 {
		if uniqueMap[item] {
			intersectionSlice = append(intersectionSlice, item)
		}
	}
	return intersectionSlice
}

// Difference returns a slice containing elements in set1 but not in set2.
// returns nil if sets are equal
func Difference[T comparable](set1, set2 []T) []T {
	uniqueMap := make(map[T]bool)
	for _, item := range set2 {
		uniqueMap[item] = true
	}

	var differenceSlice []T
	for _, item := range set1 {
		if !uniqueMap[item] {
			differenceSlice = append(differenceSlice, item)
		}
	}
	//	fmt.Println(set1, set2)

	return differenceSlice
}
