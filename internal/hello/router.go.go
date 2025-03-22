package hello

import (
	"net/http"

	"github.com/adryanev/go-http-service-template/common/db"
	"github.com/adryanev/go-http-service-template/common/messaging"
	"github.com/adryanev/go-http-service-template/common/utils"

	"github.com/go-chi/chi/v5"
)

// Hello contains all dependencies needed for the module
type Hello struct {
	DB         *db.DB
	NatsClient *messaging.NatsClient
}

// NewHello creates a new hello with dependencies
func NewHello(db *db.DB, natsClient *messaging.NatsClient) *Hello {
	return &Hello{
		DB:         db,
		NatsClient: natsClient,
	}
}

// Router returns a router for this module
func (m *Hello) Router() http.Handler {
	r := chi.NewMux()

	// Basic routes
	r.Get("/", m.testRoute)

	return r
}

// @Summary Health check
// @Description Simple health check endpoint to verify the service is running
// @Tags system
// @Produce json
// @Success 200 {object} utils.Response{message=string} "Service is healthy"
// @Router / [get]
func (m *Hello) testRoute(w http.ResponseWriter, r *http.Request) {
	utils.WriteMessage(w, 200, "Hello with dependency injection")
}
