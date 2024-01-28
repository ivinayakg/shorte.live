package routes

import (
	"github.com/gorilla/mux"
	"github.com/ivinayakg/shorte.live/api/controllers"
	"github.com/ivinayakg/shorte.live/api/middleware"
)

func URLRoutes(r *mux.Router) {
	protectedR := r.NewRoute().Subrouter()
	protectedR.Use(middleware.Authentication)
	protectedR.HandleFunc("", controllers.ShortenURL).Methods("POST")
	protectedR.HandleFunc("/all", controllers.GetUserURL).Methods("GET")
	protectedR.HandleFunc("/{id}", controllers.UpdateUrl).Methods("PATCH")
	protectedR.HandleFunc("/{id}", controllers.DeleteUrl).Methods("DELETE")
}

func URLResolveRoutes(r *mux.Router) {
	r.HandleFunc("/{short}", controllers.ResolveURL).Methods("GET")
}
