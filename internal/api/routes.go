package api

import (
	"net/http"
	"strings"

	"github.com/makson2134/go-qa-service/internal/api/handlers"
)

func SetupRoutes(h *handlers.Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.HealthCheck)

	mux.HandleFunc("/questions/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/questions/")

		if path == "" {
			switch r.Method {
			case http.MethodGet:
				h.ListQuestions(w, r)
			case http.MethodPost:
				h.CreateQuestion(w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if strings.HasSuffix(path, "/answers/") {
			if r.Method == http.MethodPost {
				h.CreateAnswer(w, r)
			} else {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.GetQuestion(w, r)
		case http.MethodDelete:
			h.DeleteQuestion(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/answers/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetAnswer(w, r)
		case http.MethodDelete:
			h.DeleteAnswer(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
