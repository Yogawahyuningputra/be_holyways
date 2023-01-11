package handlers

import (
	"encoding/json"
	"net/http"
	donationdto "server/dto/donation"
	dto "server/dto/result"
	"server/models"
	"server/repositories"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerDonation struct {
	DonationRepository repositories.DonationRepository
}

func HandlerDonation(DonationRepository repositories.DonationRepository) *handlerDonation {
	return &handlerDonation{DonationRepository}

}
func (h *handlerDonation) FindDonations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	donations, err := h.DonationRepository.FindDonations()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: donations}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerDonation) GetDonation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	donation, err := h.DonationRepository.GetDonation(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseDonation(donation)}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerDonation) CreateDonation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	status := userInfo["role"]
	if status == "admin" {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "admin can't donations"}
		json.NewEncoder(w).Encode(response)
		return
	}
	funding_id, _ := strconv.Atoi(r.FormValue("funding_id"))
	money, _ := strconv.Atoi(r.FormValue("money"))
	request := donationdto.DonationRequest{
		FundingID: funding_id,
		Money:     money,
		Status:    r.FormValue("status"),
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	donation := models.Donation{
		Money:     request.Money,
		Status:    request.Status,
		FundingID: request.FundingID,
		UserID:    userID,
	}
	donation, err = h.DonationRepository.CreateDonation(donation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "Donation Failed"}
		json.NewEncoder(w).Encode(response)
		return
	}
	donation, _ = h.DonationRepository.GetDonation(donation.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseDonation(donation)}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerDonation) UpdateDonation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))
	status := userInfo["role"]
	if status == "admin" {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "admin can't donations"}
		json.NewEncoder(w).Encode(response)
		return
	}
	funding_id, _ := strconv.Atoi(r.FormValue("funding_id"))
	money, _ := strconv.Atoi(r.FormValue("money"))
	request := donationdto.DonationRequest{
		FundingID: funding_id,
		Money:     money,
		Status:    r.FormValue("status"),
		UserID:    userID,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	donation, err := h.DonationRepository.GetDonation(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	if request.Money != 0 {
		donation.Money = request.Money
	}
	if request.FundingID != 0 {
		donation.FundingID = request.FundingID
	}
	if request.Status != "" {
		donation.Status = request.Status
	}

	donation, err = h.DonationRepository.UpdateDonation(donation, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "Donation Failed"}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseDonation(donation)}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerDonation) DeleteDonation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	donation, err := h.DonationRepository.GetDonation(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	data, err := h.DonationRepository.DeleteDonation(donation, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}
func convertResponseDonation(u models.Donation) donationdto.DonationResponse {
	return donationdto.DonationResponse{
		Money:     u.Money,
		Status:    u.Status,
		UserID:    u.UserID,
		FundingID: u.FundingID,
	}
}
