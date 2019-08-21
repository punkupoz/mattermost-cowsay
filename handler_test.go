package main

import (
	"github.com/jinzhu/gorm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/go-chi/chi"
	"strings"
)

func TestHandleCow(t *testing.T) {
	assert := assert.New(t)
	db, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		panic(err)
	}
	srv := server{
		config: &Config{
			Mattermost: Mattermost{
				Token: "test",
			},
		},
		router: chi.NewRouter(),
		db: db,
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