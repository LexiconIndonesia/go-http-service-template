package user

import (
	"github.com/go-chi/chi/v5"
)

// Router returns the router for user endpoints
func (u *User) Router() chi.Router {
	r := chi.NewRouter()
	r.Post("/", u.CreateUser)
	return r
}
