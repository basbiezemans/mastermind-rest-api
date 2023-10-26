package model

type Pair[T any] struct {
	fst T
	snd T
}

// Create a new Pair.
func NewPair[T any](fst, snd T) Pair[T] {
	return Pair[T]{fst, snd}
}

// Separate a Pair into its two elements.
func Unpair[T any](p Pair[T]) (T, T) {
	return p.fst, p.snd
}

// Unequal, given a slice of pairs, returns the pairs with unequal elements.
func Unequal[C comparable](pairs []Pair[C]) []Pair[C] {
	var unequal = make([]Pair[C], 0)
	for _, p := range pairs {
		if p.fst != p.snd {
			unequal = append(unequal, p)
		}
	}
	return unequal
}

// IsEqual, given a "comparable" value, returns a function that will compare
// for equality with another value of the same type.
func IsEqual[C comparable](a C) func(C) bool {
	return func(b C) bool {
		return a == b
	}
}

// Removes the first occurrence of x from its slice argument, if it's present.
func Remove[C comparable](x C, xs []C) ([]C, bool) {
	for i, y := range xs {
		if x == y {
			return append(xs[:i], xs[i+1:]...), true
		}
	}
	return xs, false
}
