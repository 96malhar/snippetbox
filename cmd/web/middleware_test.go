package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/96malhar/snippetbox/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to our secureHeaders
	// middleware, which writes a 200 status code and an "OK" response body.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)
	rs := rr.Result()

	expectedHeaderValues := map[string]string{
		"Content-Security-Policy": "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		"Referrer-Policy":         "origin-when-cross-origin",
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "deny",
		"X-XSS-Protection":        "0",
	}

	for header, value := range expectedHeaderValues {
		assert.Equal(t, rs.Header.Get(header), value)
	}

	// Check that the middleware has correctly called the next handler in line
	// and the response status code and body are as expected.
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}

func TestRequireAuthentication(t *testing.T) {
	tests := []struct {
		name        string
		requestCtx  context.Context
		wantStatus  int
		wantHeaders map[string]string
	}{
		{
			name:       "Authenticated context",
			requestCtx: context.WithValue(context.Background(), isAuthenticatedContextKey, true),
			wantStatus: http.StatusOK,
			wantHeaders: map[string]string{
				"Cache-Control": "no-store",
			},
		},
		{
			name:       "Unauthenticated context",
			requestCtx: context.Background(),
			wantStatus: http.StatusSeeOther,
			wantHeaders: map[string]string{
				"Location": "/user/login",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			r = r.WithContext(tc.requestCtx)

			// Create a mock HTTP handler that we can pass to our requireAuthentication middleware
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})

			app := newTestApplication(t)
			app.requireAuthentication(next).ServeHTTP(rr, r)

			result := rr.Result()
			assert.Equal(t, result.StatusCode, tc.wantStatus)

			for key, value := range tc.wantHeaders {
				assert.Equal(t, result.Header.Get(key), value)
			}
		})
	}
}
