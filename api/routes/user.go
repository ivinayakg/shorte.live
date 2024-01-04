package routes

import (
	"example.com/go/url-shortner/controllers"
	"example.com/go/url-shortner/middleware"
	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/sign_in_with_google", controllers.SignInWithGoogle).Methods("GET")
	r.HandleFunc("/google/callback", controllers.CallbackSignInWithGoogle).Methods("POST")

	protectedR := r.NewRoute().Subrouter()
	protectedR.Use(middleware.Authentication)
	protectedR.HandleFunc("/self", controllers.SelfUser).Methods("GET")
}
