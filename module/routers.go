package module

import (
	"net/http"

	"github.com/adryanev/go-http-service-template/common/db"
	"github.com/adryanev/go-http-service-template/common/utils"

	"github.com/go-chi/chi/v5"
)

// Module contains all dependencies needed for the module
type Module struct {
	DB *db.DB
}

// NewModule creates a new module with dependencies
func NewModule(db *db.DB) *Module {
	return &Module{
		DB: db,
	}
}

// Router returns a router for this module
func (m *Module) Router() http.Handler {
	r := chi.NewMux()

	// Basic routes
	r.Get("/", m.testRoute)

	// User routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/", m.userHandler)
	})

	return r
}

// Router is a legacy function (will be deprecated)
func Router() *chi.Mux {
	r := chi.NewMux()
	r.Get("/", testRoute)
	return r
}

// testRoute is the handler for the test route
func (m *Module) testRoute(w http.ResponseWriter, r *http.Request) {
	utils.WriteMessage(w, 200, "Hello with dependency injection")
}

// testRoute is a legacy function (will be deprecated)
func testRoute(w http.ResponseWriter, req *http.Request) {
	utils.WriteMessage(w, 200, "Hello")
}
