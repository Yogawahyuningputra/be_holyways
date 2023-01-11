package routes

import (
	"server/handlers"
	"server/pkg/middleware"
	"server/pkg/mysql"
	"server/repositories"

	"github.com/gorilla/mux"
)

func DonationRoutes(r *mux.Router) {
	donationRepository := repositories.RepositoryDonation(mysql.DB)
	h := handlers.HandlerDonation(donationRepository)
	r.HandleFunc("/donations", middleware.Auth(h.FindDonations)).Methods("GET")
	r.HandleFunc("/donation/{id}", middleware.Auth(h.GetDonation)).Methods("GET")
	r.HandleFunc("/donation", middleware.Auth(h.CreateDonation)).Methods("POST")
	r.HandleFunc("/donation/{id}", middleware.Auth(h.UpdateDonation)).Methods("PATCH")
	r.HandleFunc("/donation/{id}", middleware.Auth(h.DeleteDonation)).Methods("DELETE")

}
