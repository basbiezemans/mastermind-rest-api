package internal

import (
	"errors"
	"math/rand"
	"regexp"
)

type Code struct {
	Digits []rune `json:"digits"`
}

func (c Code) String() string {
	return string(c.Digits)
}

func RandomCode() Code {
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
	return Code{digits}
}

func newCode(code string) Code {
	return Code{
		Digits: []rune(code),
	}
}

func isValidCode(code string) bool {
	re := regexp.MustCompile(`^[1-6]{4}$`)
	return re.MatchString(code)
}

func CodeFromString(code string) (Code, error) {
	if isValidCode(code) {
		return newCode(code), nil
	}
	return Code{}, errors.New("invalid guess")
}
