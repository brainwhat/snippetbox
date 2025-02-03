package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"snippetbox.brainwhat/internal/models/mocks"
)

func NewTestApp(t *testing.T) *application {
	// The same as with the app
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	// Also the same
	formDecoder := form.NewDecoder()

	// We don't set .Store as there is no DB
	// Also default setting is using "transient in-memory store"
	// Which is perfect for testing
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		// Mock loggers that discard everything given to them
		// logRequest and recoverPanic middlewares use them
		infoLog:  log.New(io.Discard, "", 0),
		errorLog: log.New(io.Discard, "", 0),

		// Created mocks using interfaces so there's no need to spin up a DB instance
		snippets: &mocks.MockSnippetModel{},
		users:    &mocks.MockUserModel{},

		// This is the same as with normal app
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

// Define type that embeds Server instance
type testServer struct {
	*httptest.Server
}

// Return new instance of our custom server
func NewTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// We add this part to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	// We don't want to follow redirects but return the first response sent by the server
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// Make a GET req to url using test server client
// Return response status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	// Get resp after calling url
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	// Read the body
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
