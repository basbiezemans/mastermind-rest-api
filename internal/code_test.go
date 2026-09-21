package internal

import (
	"reflect"
	"testing"
)

func TestNewSecret(t *testing.T) {
	var secret Code
	for i, n := 0, 10; i < n; i++ {
		secret = RandomCode()
		if !isValidCode(secret.String()) {
			t.Errorf("invalid secret code: %s", secret.String())
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
