package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/cors"
)

func CORSHandler(appHost string) func(h http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{fmt.Sprintf("https://%s", appHost)},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
