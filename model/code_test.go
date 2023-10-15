package model

import (
	"reflect"
	"testing"
)

func TestNewSecret(t *testing.T) {
	var secret Secret
	var length int
	for i, n := 0, 10; i < n; i++ {
		secret = NewSecret()
		length = len(secret.Code.Digits)
		// fmt.Printf("Secret: %s\n", secret.Code.String())
		if length != 4 {
			t.Errorf("invalid length, expected 4, got %d", length)
		}
		for _, d := range secret.Code.Digits {
			if !isValidDigit(d) {
				t.Errorf("invalid digit, expected 1..6, got %s", string(d))
			}
		}
	}
}

func TestCodeFromString(t *testing.T) {
	var tests = []struct {
		descr  string
		guess  string
		expect Code
	}{
		{"valid guess", "1234", newCode("1234")},
		{"guess too short", "123", Code{}},
		{"guess too long", "12345", Code{}},
		{"invalid digits", "1237", Code{}},
		{"invalid digits", "0234", Code{}},
		{"invalid digits", "1e34", Code{}},
	}
	for _, test := range tests {
		code, _ := CodeFromString(test.guess)
		if !reflect.DeepEqual(code, test.expect) {
			ans := test.expect.String()
			got := code.String()
			t.Errorf("%s, expected %q, got %q", test.descr, ans, got)
		}
	}
}
