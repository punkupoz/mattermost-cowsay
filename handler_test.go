package main

import (
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandleCow(t *testing.T) {
	assert := assert.New(t)
	repo, _ := newRepo("sqlite3", "test.db")
	srv := server{
		config: &Config{
			Mattermost: Mattermost{
				Token: "test",
			},
		},
		router: chi.NewRouter(),
		repo: repo,
	}
	srv.routes()

	t.Run("cow does say", func (t *testing.T) {
		data := url.Values{}
		data.Set("token", "test")
		data.Set("text", "cow")

		body := strings.NewReader(data.Encode())
		req, err := http.NewRequest("POST", "/cowsay", body)
		assert.NoError(err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		srv.router.ServeHTTP(w, req)
		assert.Equal(w.Code, http.StatusOK)
	})

	t.Run("cow not authorized", func (t *testing.T) {
		data := url.Values{}
		data.Set("token", "testtest")
		data.Set("text", "cow")

		body := strings.NewReader(data.Encode())
		req, err := http.NewRequest("POST", "/cowsay", body)
		assert.NoError(err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		srv.router.ServeHTTP(w, req)
		assert.Equal(w.Code, http.StatusUnauthorized)
	})
}