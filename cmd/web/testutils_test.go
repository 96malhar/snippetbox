package main

import (
	"bytes"
	"github.com/96malhar/snippetbox/internal/store/mocks"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// newTestApplication returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	// And a session manager instance. Note that we use the same settings as
	// production, except that we *don't* set a Store for the session manager.
	// If no store is set, the SCS package will default to using a transient
	// in-memory store, which is ideal for testing purposes.
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		errorLog:       log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		snippetStore:   mocks.NewMockSnippetStore(), // Use the mock.
		userStore:      mocks.NewMockUserStore(),    // Use the mock.
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

// Define a custom testServer type which embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// newTestServer initializes and returns a new instance of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	// Disable redirect-following for the test server client by setting a custom
	// CheckRedirect function. This function will be called whenever a 3xx
	// response is received by the client, and by always returning a
	// http.ErrUseLastResponse error it forces the client to immediately return
	// the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// get makes a GET request to a given url path using the test server client, and returns the
// response status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) *http.Response {
	resp, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	return resp
}

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) *http.Response {
	resp, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func getString(t *testing.T, r io.Reader) string {
	t.Helper()
	body, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	return string(body)
}
