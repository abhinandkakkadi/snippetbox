package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {

	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock handler that we can pass to our secureHeaders middleware
	next := http.HandlerFunc(func(w http.ResponseWriter, t *http.Request) {
		w.Write([]byte("OK"))
	})

	// pass the mock http handler to secure headers middleware
	// Since secure headers returns http.Handler we can call its ServeHTTP() method
	secureHeaders(next).ServeHTTP(rr,r)

	rs := rr.Result()
	_ = rs

}