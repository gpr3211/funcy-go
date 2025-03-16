package monad

import ()

// Maybe interface defines the common operations
type Maybe[A any] interface {
	Get() A
	GetOrElse(def A) A
}

// JustMaybe represents a present value
type JustMaybe[A any] struct {
	value A
}

func (j JustMaybe[A]) Get() A {
	return j.value
}

func (j JustMaybe[A]) GetOrElse(def A) A {
	return j.value
}

// NothingMaybe represents an absent value
type NothingMaybe[A any] struct{}

func (n NothingMaybe[A]) Get() A {
	return *new(A)
}

func (n NothingMaybe[A]) GetOrElse(def A) A {
	return def
}

// Constructor functions
func Just[A any](a A) JustMaybe[A] {
	return JustMaybe[A]{value: a}
}

func Nothing[A any]() Maybe[A] {
	return NothingMaybe[A]{}
}

// Functor implementation
func fmap[A, B any](m Maybe[A], mapFunc func(A) B) Maybe[B] {
	switch m.(type) {
	case JustMaybe[A]:
		j := m.(JustMaybe[A])
		return JustMaybe[B]{
			value: mapFunc(j.value),
		}
	case NothingMaybe[A]:
		return NothingMaybe[B]{}
	default:
		panic("unknown type")
	}
}
