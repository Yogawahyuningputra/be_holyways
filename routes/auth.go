package routes

import (
	"server/handlers"
	"server/pkg/middleware"
	"server/pkg/mysql"
	"server/repositories"

	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router) {
	authRepository := repositories.RepositoryAuth(mysql.DB)
	h := handlers.HandlerAuth(authRepository)
	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/checkAuth", middleware.Auth(h.CheckAuth)).Methods("GET")

}
