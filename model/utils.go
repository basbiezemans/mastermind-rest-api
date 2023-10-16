package model

type Pair[T any] struct {
	fst T
	snd T
}

func NewPair[T any](fst, snd T) Pair[T] {
	return Pair[T]{fst, snd}
}

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

func IsEqual[C comparable](a C) func(C) bool {
	return func(b C) bool {
		return a == b
	}
}

// FindIndex, takes a predicate and a slice and returns the index of the first
// element in the slice satisfying the predicate (index,true), or (-1,false) if
// there is no such element.
func FindIndex[C comparable](fn func(C) bool, xs []C) (int, bool) {
	for i, x := range xs {
		if fn(x) {
			return i, true
		}
	}
	return -1, false
}

func Delete[C comparable](x C, xs []C) []C {
	if i, ok := FindIndex(IsEqual(x), xs); ok {
		return append(xs[:i], xs[i+1:]...)
	}
	return xs
}
