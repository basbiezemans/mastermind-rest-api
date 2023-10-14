package model

import (
	"errors"
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

// Stub
func NewSecret() Secret {
	var digits = []rune{'1', '2', '3', '4'}
	return Secret{Code{digits}}
}

// Stub
func CodeFromString(guess string) (Code, error) {
	var digits = []rune(guess)
	if len(digits) != 4 {
		return Code{}, errors.New("invalid guess")
	}
	// isDigit check for each rune
	// inRange check for each digit [1..6]
	return Code{digits}, nil
}
