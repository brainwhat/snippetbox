package main

import (
	"net/http"
	"testing"

	"snippetbox.brainwhat/internal/assert"
)

func TestPing(t *testing.T) {
	app := NewTestApp(t)

	ts := NewTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
