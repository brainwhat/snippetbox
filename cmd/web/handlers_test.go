package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.brainwhat/internal/assert"
)

// We'll be using ResponseRecorder which is kinda ResponseWriter, but
// It records status code, headers and body, but doesn't write them to HTTP conn
func TestPing(t *testing.T) {
	rr := httptest.NewRecorder()

	// Dummy req
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call tested func, get the response, check if it's OK
	ping(rr, r)
	res := rr.Result()
	assert.Equal(t, res.StatusCode, http.StatusOK)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
