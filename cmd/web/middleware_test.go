package main

import (
	"bytes"
	"context"
	"github.com/96malhar/snippetbox/internal/store"
	"github.com/96malhar/snippetbox/internal/store/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
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
		assert.Equal(t, value, rs.Header.Get(header))
	}

	// Check that the middleware has correctly called the next handler in line
	// and the response status code and body are as expected.
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	require.NoError(t, err)
	body = bytes.TrimSpace(body)
	assert.Equal(t, "OK", string(body))
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
			app := newTestApplication(t)
			reqCtx, err := app.sessionManager.Load(tc.requestCtx, "")
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r, err := http.NewRequestWithContext(reqCtx, http.MethodGet, "/", nil)
			require.NoError(t, err)

			// Create a mock HTTP handler that we can pass to our requireAuthentication middleware
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})

			app.requireAuthentication(next).ServeHTTP(rr, r)

			result := rr.Result()
			assert.Equal(t, tc.wantStatus, result.StatusCode)

			for key, value := range tc.wantHeaders {
				assert.Equal(t, value, result.Header.Get(key))
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	// Create a mock HTTP handler that we can pass to our authenticate middleware
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
		if !ok {
			isAuthenticated = false
		}
		w.Header().Add("IsAuthenticated", strconv.FormatBool(isAuthenticated))
		w.Write([]byte("OK"))
	})

	// Initialize app
	app := newTestApplication(t)
	reqCtx, err := app.sessionManager.Load(context.Background(), "")
	app.sessionManager.Put(reqCtx, "authenticatedUserID", 1)
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		name                string
		userStore           userStoreInterface
		wantIsAuthenticated string
	}{
		{
			name:                "User exists",
			userStore:           mocks.NewMockUserStore(&store.User{ID: 1, Name: "John"}),
			wantIsAuthenticated: "true",
		},
		{
			name:                "User does not exists",
			userStore:           mocks.NewMockUserStore(),
			wantIsAuthenticated: "false",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			r, err := http.NewRequestWithContext(reqCtx, http.MethodGet, "/", nil)
			require.NoError(t, err)

			app.userStore = tc.userStore
			app.authenticate(next).ServeHTTP(rr, r)
			isAuthenticated := rr.Result().Header.Get("IsAuthenticated")
			assert.Equal(t, tc.wantIsAuthenticated, isAuthenticated)
		})
	}
}
