package monad

import "strconv"

// Safe map access
func getFromMap[K comparable, V any](m map[K]V, key K) Maybe[V] {
	if value, ok := m[key]; ok {
		return Just(value)
	}
	return Nothing[V]()
}

func parseNumber(s string) Maybe[int] {
	if num, err := strconv.Atoi(s); err == nil {
		return Just(num)
	}
	return Nothing[int]()
}

// Example using fromNullable
func fromNullable[A any](ptr *A) Maybe[A] {
	if ptr == nil {
		return Nothing[A]()
	}
	return Just(*ptr)
}

/*
// Example functions
func getUserByID(id int) Maybe[User] {
	users := map[int]User{
		1: {"Alice", 30},
		2: {"Bob", 25},
	}
	if user, exists := users[id]; exists {
		return Just(user)
	}
	return Nothing[User]()
}
*/
