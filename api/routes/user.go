package routes

import (
	"github.com/gorilla/mux"
	"github.com/ivinayakg/shorte.live/api/controllers"
	"github.com/ivinayakg/shorte.live/api/middleware"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/sign_in_with_google", controllers.SignInWithGoogle).Methods("GET")
	r.HandleFunc("/google/callback", controllers.CallbackSignInWithGoogle).Methods("GET")

	protectedR := r.NewRoute().Subrouter()
	protectedR.Use(middleware.Authentication)
	protectedR.HandleFunc("/self", controllers.SelfUser).Methods("GET")
}
