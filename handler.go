package main

import (
	"net/http"
)

func (s *server) handleCowsay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		token := r.Form.Get("token")
		if token != s.config.Mattermost.Token {
			http.Error(w, "UNAUTHORIZED", http.StatusUnauthorized)
			return
		}

		text := r.Form.Get("text")
		if _, err := w.Write([]byte(generateCow(text))); err != nil {
			http.Error(w, "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
			return
		}
	}
}