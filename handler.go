package main

import (
	"encoding/json"
	"net/http"
)

type MattermostResponse struct {
	ResponseType string `json:"response_type"`
	Text string `json:"text"`
}

func prepareResponseMatter(text string, responseType string) ([]byte, error) {
	res := &MattermostResponse{
		ResponseType: responseType,
		Text:         text,
	}

	ret, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *server) handleCowsay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		token := r.Form.Get("token")
		userId := r.Form.Get("user_id")
		text := r.Form.Get("text")

		activity := &ActivityLog{
			UserID: userId,
			Message: text,
		}

		if token != s.config.Mattermost.Token {
			http.Error(w, "UNAUTHORIZED", http.StatusUnauthorized)
			activity.Success = false
			s.db.Create(activity)
			return
		}

		res, err := prepareResponseMatter(generateCow(text), "in_channel")
		if err != nil {
			panic(err)
		}

		_, err = w.Write(res);
		if err != nil {
			http.Error(w, "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
			activity.Success = false
			s.db.Create(activity)
			return
		}

		activity.Success = true
		s.db.Create(activity)
	}
}

func (s *server) handleRetrieveLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res, err := prepareResponseMatter("something", "in_channel")
		if err != nil {
			panic(err)
		}
		w.Write(res)
	}
}