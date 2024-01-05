package routes

import (
	"example.com/go/url-shortner/controllers"
	"example.com/go/url-shortner/middleware"
	"github.com/gorilla/mux"
)

func URLRoutes(r *mux.Router) {
	r.HandleFunc("/{short}", controllers.ResolveURL).Methods("GET")

	protectedR := r.NewRoute().Subrouter()
	protectedR.Use(middleware.Authentication)
	protectedR.HandleFunc("", controllers.ShortenURL).Methods("POST")
}
