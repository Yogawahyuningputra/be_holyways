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
	r.HandleFunc("/donationByUser", middleware.Auth(h.GetDonationByUser)).Methods("GET")
	r.HandleFunc("/donationByFunding/{id}", middleware.Auth(h.GetDonationByFunding)).Methods("GET")
	r.HandleFunc("/donationPending/{id}", middleware.Auth(h.GetDonationPending)).Methods("GET")
	r.HandleFunc("/notification", h.Notification).Methods("POST")

}
