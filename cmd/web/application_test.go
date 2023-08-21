package main

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplication_ServerError(t *testing.T) {
	app := newTestApplication(t)
	rr := httptest.NewRecorder()
	app.serverError(rr, errors.New("some error"))

	res := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestApplication_ClientError(t *testing.T) {
	statusCodes := []int{http.StatusBadRequest, http.StatusUnprocessableEntity}
	app := newTestApplication(t)

	for _, code := range statusCodes {
		t.Run(http.StatusText(code), func(t *testing.T) {
			rr := httptest.NewRecorder()
			app.clientError(rr, code)

			assert.Equal(t, code, rr.Result().StatusCode)
		})
	}
}

func TestApplication_NotFound(t *testing.T) {
	app := newTestApplication(t)
	rr := httptest.NewRecorder()
	app.notFound(rr)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestApplication_Render(t *testing.T) {
	testcases := []struct {
		name        string
		pageName    string
		data        *templateData
		wantContent string
		wantStatus  int
	}{
		{
			name:        "Valid template",
			pageName:    "create.tmpl",
			data:        &templateData{Form: snippetCreateForm{}},
			wantContent: "Create a New Snippet",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "Invalid template",
			pageName:    "does-not-exist.tmpl",
			data:        &templateData{},
			wantContent: "Internal Server Error",
			wantStatus:  http.StatusInternalServerError,
		},
	}

	app := newTestApplication(t)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			app.render(rr, http.StatusOK, tc.pageName, tc.data)

			statusCode := rr.Code
			body, err := io.ReadAll(rr.Body)

			assert.NoError(t, err)
			assert.Contains(t, string(body), tc.wantContent)
			assert.Equal(t, tc.wantStatus, statusCode)
		})
	}
}

func TestApplication_IsAuthenticated(t *testing.T) {
	app := newTestApplication(t)
	testcases := []struct {
		name                string
		requestContext      context.Context
		wantIsAuthenticated bool
	}{
		{
			name:                "Is authenticated",
			requestContext:      context.WithValue(context.Background(), isAuthenticatedContextKey, true),
			wantIsAuthenticated: true,
		},
		{
			name:                "Not authenticated",
			requestContext:      context.Background(),
			wantIsAuthenticated: false,
		},
		{
			name:                "AuthenticatedContextKey is not bool",
			requestContext:      context.WithValue(context.Background(), isAuthenticatedContextKey, 100),
			wantIsAuthenticated: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(tc.requestContext, http.MethodGet, "/", nil)
			assert.NoError(t, err)
			isAuthenticated := app.isAuthenticated(req)
			assert.Equal(t, tc.wantIsAuthenticated, isAuthenticated)
		})
	}
}

func TestApplication_DecodePostForm(t *testing.T) {
	app := newTestApplication(t)

	type personForm struct {
		FirstName string `form:"firstName"`
		LastName  string `form:"lastName"`
		Age       int    `form:"age"`
	}

	validFormData := map[string]string{"firstName": "John", "lastName": "Smith", "age": "100"}
	invalidFormData := map[string]string{"firstName": "John", "lastName": "Smith", "age": "Hello, world!"}

	t.Run("Valid decoding", func(t *testing.T) {
		postReq := createPostRequest(t, "/", validFormData)

		dst := personForm{}
		err := app.decodePostForm(postReq, &dst)
		if err != nil {
			t.Errorf("Unexpected error occurred while decoding post form. Error = %v", err)
		}

		wantDst := personForm{
			FirstName: "John", LastName: "Smith", Age: 100,
		}

		assert.Equal(t, wantDst, dst)
	})

	t.Run("Should panic on non-pointer destination", func(t *testing.T) {
		postReq := createPostRequest(t, "/", validFormData)

		assert.Panics(t, func() {
			app.decodePostForm(postReq, personForm{})
		})
	})

	t.Run("Should error on invalid form data", func(t *testing.T) {
		postReq := createPostRequest(t, "/", invalidFormData)

		err := app.decodePostForm(postReq, &personForm{})
		assert.Error(t, err)
	})
}
