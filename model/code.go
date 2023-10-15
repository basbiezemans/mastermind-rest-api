package model

import (
	"errors"
	"math/rand"
	"unicode"

	slice "github.com/basbiezemans/gofunctools/functools"
)

type Secret struct {
	Code Code
}

type Code struct {
	Digits []rune
}

func (c Code) String() string {
	return string(c.Digits)
}

func NewSecret() Secret {
	var digits = make([]rune, 4)
	var valid = []rune("123456")
	// The below rand.Perm returns, as a slice of 6 ints, a pseudo-random
	// permutation of the integers in the half-open interval [0,6).
	// Even capped at 4 ints, it's sufficiently random for the purpose
	// of generating a secret code.
	var rs = rand.Perm(6)[:4]
	for i, r := range rs {
		digits[i] = valid[r]
	}
	return Secret{Code{digits}}
}

func newCode(code string) Code {
	return Code{
		Digits: []rune(code),
	}
}

func isValidDigit(r rune) bool {
	var valid = []rune("123456")
	return unicode.IsDigit(r) && slice.Any(equals(r), valid)
}

func equals(r rune) func(rune) bool {
	return func(e rune) bool {
		return r == e
	}
}

func CodeFromString(guess string) (Code, error) {
	var chars = []rune(guess)
	var isValidLen = len(chars) == 4
	var isValidGuess = isValidLen && slice.All(isValidDigit, chars)
	if isValidGuess {
		return newCode(guess), nil
	}
	return Code{}, errors.New("invalid guess")
}