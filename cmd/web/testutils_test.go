package main

import (
	"bytes"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
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

	// We don"t set .Store as there is no DB
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

		// Created mocks using interfaces so there"s no need to spin up a DB instance
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

	// We don"t want to follow redirects but return the first response sent by the server
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

// Code is identical to get(), we just call Client().PostForm() instead of Get()
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

// Я не хочу учить regexp
var csrfTokenRx = regexp.MustCompile(`<input\s+type="hidden"\s+name="csrf_token"\s+value="([^"]+)"`)

func extactCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRx.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no CSRF tokens found in body")
	}

	// html/template escapes(encodes) some characters
	// like `+` will be returned as `&#43;`
	// We use UnescapeString to get original string
	return html.UnescapeString(string(matches[1]))
}
