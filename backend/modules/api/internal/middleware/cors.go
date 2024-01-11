package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/cors"
	"github.com/jo-fr/activityhub/backend/pkg/config"
)

func CORSHandler(config config.Config) func(h http.Handler) http.Handler {

	var origins []string
	origins = append(origins, fmt.Sprintf("https://%s", config.AppHost))
	if config.Environment.IsLocal() {
		origins = append(origins, "http://localhost:5173")
	}
	return cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
