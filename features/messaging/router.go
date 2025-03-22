package messaging

import (
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// Router returns the router for messaging endpoints
func (m *Messaging) Router() chi.Router {
	r := chi.NewRouter()

	// Only enable messaging routes in development environment
	env := os.Getenv("APP_ENV")
	if strings.ToLower(env) == "development" {
		r.Post("/publish", m.PublishMessage)
		r.Get("/subscribe/{subject}", m.SubscribeWebSocket)
		log.Info().Msg("Messaging endpoints enabled in development mode")
	} else {
		log.Info().Msg("Messaging endpoints disabled in production mode")
	}

	return r
}
