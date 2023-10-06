package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {

	// Define that this is a helper for testing
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want %v", actual, expected)
	}

}
