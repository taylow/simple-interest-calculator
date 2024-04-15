package main

import "testing"

func TestRandomString(t *testing.T) {
	for i := range 10 {
		str := randomString(i)
		if len(str) != i {
			t.Errorf("Unexpected length of random string. got %d, want %d", len(str), i)
		}
	}
}
