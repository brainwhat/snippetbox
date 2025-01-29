package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {
	// This marks Equal() as test helper function
	// This means that Errorf will report filenamea and line of the func that called Equal()
	t.Helper()

	if actual != expected {
		t.Errorf("got %v, expected %v,", expected, actual)
	}
}
