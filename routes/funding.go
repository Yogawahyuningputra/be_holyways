package routes

import (
	"server/handlers"
	"server/pkg/middleware"
	"server/pkg/mysql"
	"server/repositories"

	"github.com/gorilla/mux"
)

func FundingRoutes(r *mux.Router) {
	fundingRepository := repositories.RepositoryFunding(mysql.DB)
	h := handlers.HandlerFunding(fundingRepository)
	r.HandleFunc("/fundings", h.FindFundings).Methods("GET")
	r.HandleFunc("/funding/{id}", h.GetFunding).Methods("GET")
	r.HandleFunc("/fundingByUser", middleware.Auth(h.GetFundingByUser)).Methods("GET")
	r.HandleFunc("/funding", middleware.Auth(middleware.UploadFile(h.CreateFunding))).Methods("POST")
	r.HandleFunc("/funding/{id}", middleware.Auth(middleware.UploadFile(h.UpdateFunding))).Methods("PATCH")
	r.HandleFunc("/funding/{id}", middleware.Auth(h.DeleteFunding)).Methods("DELETE")

}
