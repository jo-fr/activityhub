package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

func CORSHandler(appURL string) func(h http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{appURL},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
