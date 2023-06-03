package main

import (
	"context"
	"github.com/96malhar/snippetbox/internal/assert"
	"net/http"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2022 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}

func TestNewTemplateCache(t *testing.T) {
	cache, err := newTemplateCache()
	if err != nil {
		t.Fatalf("Unexpected error = %v", err)
	}

	expectedCacheEntries := []string{
		"create.tmpl", "home.tmpl", "login.tmpl", "signup.tmpl", "view.tmpl",
	}

	if len(cache) != len(expectedCacheEntries) {
		t.Errorf("len(cache) = %d; want = %d", len(cache), len(expectedCacheEntries))
	}

	for _, key := range expectedCacheEntries {
		if cache[key] == nil {
			t.Errorf("Could not find template for key = %s", key)
		}
	}
}

func TestNewTemplateData(t *testing.T) {
	app := newTestApplication(t)
	ctx, err := app.sessionManager.Load(context.Background(), "session-token")
	if err != nil {
		t.Fatalf("An error occurred loading data via sessionManager; err = %v", err)
	}

	flashMessage := "This is a flash message"
	app.sessionManager.Put(ctx, "flash", flashMessage)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", nil)
	if err != nil {
		t.Fatalf("Unexpected error while creating a new request; err = %v", err)
	}

	td := app.newTemplateData(req)
	if td.Flash != flashMessage {
		t.Fatalf("td.Flash = %s; want = %s", td.Flash, flashMessage)
	}
}
