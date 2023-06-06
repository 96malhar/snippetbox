package main

import (
	"github.com/96malhar/snippetbox/internal/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	resp := ts.get(t, "/ping")
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, getString(t, resp.Body), "OK")
}

func TestHome(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	resp := ts.get(t, "/")
	defer resp.Body.Close()
	body := getString(t, resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	titleTag := "<title>Home - Snippetbox</title>"
	assert.StringContains(t, body, titleTag)
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	app.snippetStore.Insert("Snippet 1", "Content for snippet 1...", 10)
	app.snippetStore.Insert("Snippet 2", "Content for snippet 2...", 5)

	testcases := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID 1",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "Content for snippet 1...",
		},
		{
			name:     "Valid ID 2",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusOK,
			wantBody: "Content for snippet 2...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/3",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			resp := ts.get(t, tc.urlPath)
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, tc.wantCode)

			if tc.wantBody != "" {
				assert.StringContains(t, getString(t, resp.Body), tc.wantBody)
			}
		})
	}
}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		resp := ts.get(t, "/snippet/create")
		assert.Equal(t, resp.StatusCode, http.StatusSeeOther)
		assert.Equal(t, resp.Header.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		// Insert a dummy user in the userStore
		app.userStore.Insert("alice", "alice@example.com", "pa$$word")

		// Make a POST /user/login request using the dummy user inserted above
		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "pa$$word")
		ts.postForm(t, "/user/login", form)

		// Then check that the authenticated user is shown the create snippet form.
		resp := ts.get(t, "/snippet/create")
		defer resp.Body.Close()
		body := getString(t, resp.Body)
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.StringContains(t, body, "<form action='/snippet/create' method='POST'>")

		// logout the user
		resp = ts.postForm(t, "/user/logout", nil)
		assert.Equal(t, resp.StatusCode, http.StatusSeeOther)
		assert.Equal(t, resp.Header.Get("Location"), "/")

		// Try getting /snippet/create again - this should fail
		resp = ts.get(t, "/snippet/create")
		assert.Equal(t, resp.StatusCode, http.StatusSeeOther)
		assert.Equal(t, resp.Header.Get("Location"), "/user/login")
	})
}

func TestSnippetCreatePost(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	validTitle := "This is a snippet title"
	validContent := "This is snippet content"
	validExpires := "7"

	t.Run("Unauthenticated", func(t *testing.T) {
		form := url.Values{}
		form.Add("title", validTitle)
		form.Add("content", validContent)
		form.Add("expires", validExpires)
		resp := ts.postForm(t, "/snippet/create", form)

		// The post request fails and redirects the user to the login page
		assert.Equal(t, resp.StatusCode, http.StatusSeeOther)
		assert.Equal(t, resp.Header.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		// Insert a dummy user in the userStore
		app.userStore.Insert("alice", "alice@example.com", "pa$$word")

		// Make a POST /user/login request using the dummy user inserted above
		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "pa$$word")
		ts.postForm(t, "/user/login", form)

		tests := []struct {
			name           string
			snippetTitle   string
			snippetContent string
			snippetExpires string
			wantStatusCode int
			wantHeaders    map[string]string
		}{
			{
				name:           "Valid form",
				snippetTitle:   validTitle,
				snippetContent: validContent,
				snippetExpires: validExpires,
				wantStatusCode: http.StatusSeeOther,
				wantHeaders:    map[string]string{"Location": "/snippet/view/1"},
			},
			{
				name:           "Empty title",
				snippetTitle:   "",
				snippetContent: validContent,
				snippetExpires: validExpires,
				wantStatusCode: http.StatusUnprocessableEntity,
			},
			{
				name:           "Empty content",
				snippetTitle:   validTitle,
				snippetContent: "",
				snippetExpires: validExpires,
				wantStatusCode: http.StatusUnprocessableEntity,
			},
			{
				name:           "Invalid expiration",
				snippetTitle:   validTitle,
				snippetContent: validContent,
				snippetExpires: "10",
				wantStatusCode: http.StatusUnprocessableEntity,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				form := url.Values{}
				form.Add("title", tc.snippetTitle)
				form.Add("content", tc.snippetContent)
				form.Add("expires", tc.snippetExpires)
				resp := ts.postForm(t, "/snippet/create", form)

				assert.Equal(t, resp.StatusCode, tc.wantStatusCode)

				for key, val := range tc.wantHeaders {
					assert.Equal(t, resp.Header.Get(key), val)
				}
			})
		}
	})
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	const (
		wantTitle         = "<title>Signup - Snippetbox</title>"
		wantNameInput     = "<input type='text' name='name' value=''>"
		wantEmailInput    = "<input type='email' name='email' value=''>"
		wantPasswordInput = "<input type='password' name='password'>"
	)

	resp := ts.get(t, "/user/signup")
	defer resp.Body.Close()

	body := getString(t, resp.Body)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.StringContains(t, body, wantTitle)
	assert.StringContains(t, body, wantNameInput)
	assert.StringContains(t, body, wantEmailInput)
	assert.StringContains(t, body, wantPasswordInput)
}

func TestUserSignupPost(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		wantCode     int
		wantFormTag  string
		wantHeaders  map[string]string
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			wantCode:     http.StatusSeeOther,
			wantHeaders: map[string]string{
				"Location": "/user/login",
			},
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "bob@example.",
			userPassword: validPassword,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pa$$",
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate email",
			userName:     validName,
			userEmail:    "dupe@example.com",
			userPassword: validPassword,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)

			resp := ts.postForm(t, "/user/signup", form)
			defer resp.Body.Close()
			body := getString(t, resp.Body)
			assert.Equal(t, resp.StatusCode, tt.wantCode)
			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}

			for key, value := range tt.wantHeaders {
				assert.Equal(t, resp.Header.Get(key), value)
			}
		})
	}
}

func TestUserLogin(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	const (
		wantFormTag       = "<form action='/user/login' method='POST' novalidate>"
		wantTitle         = "<title>Login - Snippetbox</title>"
		wantEmailInput    = "<input type='email' name='email' value=''>"
		wantPasswordInput = "<input type='password' name='password'>"
	)

	resp := ts.get(t, "/user/login")
	defer resp.Body.Close()

	body := getString(t, resp.Body)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.StringContains(t, body, wantFormTag)
	assert.StringContains(t, body, wantTitle)
	assert.StringContains(t, body, wantEmailInput)
	assert.StringContains(t, body, wantPasswordInput)
}

func TestUserLoginPost(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	const (
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/user/login' method='POST' novalidate>"
	)

	app.userStore.Insert("validName", validEmail, validPassword)

	tests := []struct {
		name         string
		userEmail    string
		userPassword string
		wantCode     int
		wantFormTag  string
		wantHeaders  map[string]string
	}{
		{
			name:         "User found",
			userEmail:    validEmail,
			userPassword: validPassword,
			wantCode:     http.StatusSeeOther,
			wantHeaders: map[string]string{
				"Location": "/snippet/create",
			},
		},
		{
			name:         "user not found",
			userEmail:    "notfound@example.com",
			userPassword: validPassword,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "invalid email",
			userEmail:    "Invalid.Email",
			userPassword: validPassword,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "empty email",
			userEmail:    "",
			userPassword: validPassword,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "empty password",
			userEmail:    validEmail,
			userPassword: "",
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tc.userEmail)
			form.Add("password", tc.userPassword)

			resp := ts.postForm(t, "/user/login", form)
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, tc.wantCode)
			if tc.wantFormTag != "" {
				assert.StringContains(t, getString(t, resp.Body), tc.wantFormTag)
			}

			for key, value := range tc.wantHeaders {
				assert.Equal(t, resp.Header.Get(key), value)
			}
		})
	}
}

func TestUserLogoutPost(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		resp := ts.postForm(t, "/user/logout", nil)
		assert.Equal(t, resp.StatusCode, http.StatusSeeOther)
		assert.Equal(t, resp.Header.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		// Insert a dummy user in the userStore
		app.userStore.Insert("alice", "alice@example.com", "pa$$word")

		// Make a POST /user/login request using the dummy user inserted above
		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "pa$$word")
		ts.postForm(t, "/user/login", form)

		// Then check that the authenticated user is logged out successfully
		resp := ts.postForm(t, "/user/logout", nil)
		assert.Equal(t, resp.StatusCode, http.StatusSeeOther)
		assert.Equal(t, resp.Header.Get("Location"), "/")
	})

}
