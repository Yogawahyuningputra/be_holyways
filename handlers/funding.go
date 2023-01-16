package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	fundingdto "server/dto/funding"
	dto "server/dto/result"
	"server/models"
	"server/repositories"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerFunding struct {
	FundingRepository repositories.FundingRepository
}

func HandlerFunding(FundingRepository repositories.FundingRepository) *handlerFunding {
	return &handlerFunding{FundingRepository}
}

func (h *handlerFunding) FindFundings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fundings, err := h.FundingRepository.FindFundings()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: fundings}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerFunding) GetFunding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	funding, err := h.FundingRepository.GetFunding(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: funding}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerFunding) GetFundingByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	fundings, err := h.FundingRepository.GetFundingByUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: fundings}
	json.NewEncoder(w).Encode(response)

}
func (h *handlerFunding) CreateFunding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	dataContext := r.Context().Value("dataFile")
	filepath := dataContext.(string)
	goals, _ := strconv.Atoi(r.FormValue("goals"))

	request := fundingdto.FundingRequest{
		Title:       r.FormValue("title"),
		Goals:       goals,
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
	}
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	// Get cloudinary from .env
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add Cloudinary credentials
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "holyways/funding"})

	if err != nil {
		fmt.Println(err.Error())
	}
	funding := models.Funding{
		Title:       request.Title,
		Goals:       request.Goals,
		Description: request.Description,
		Image:       resp.SecureURL,
		UserID:      userID,
	}
	funding, err = h.FundingRepository.CreateFunding(funding)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	funding, _ = h.FundingRepository.GetFunding(funding.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: funding}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerFunding) UpdateFunding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	dataContext := r.Context().Value("dataFile")
	filepath := dataContext.(string)

	goals, _ := strconv.Atoi(r.FormValue("goals"))

	request := fundingdto.FundingUpdate{
		Title:       r.FormValue("title"),
		Goals:       goals,
		Description: r.FormValue("description"),
		Image:       filepath,
		UserID:      userID,
	}
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	// Get cloudinary from .env
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add Cloudinary credentials
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "holyways/funding"})

	if err != nil {
		fmt.Println(err.Error())
	}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	funding, err := h.FundingRepository.GetFunding(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "get data failed"}
		json.NewEncoder(w).Encode(response)
		return
	}
	if request.Title != "" {
		funding.Title = request.Title
	}
	if request.Goals != 0 {
		funding.Goals = request.Goals
	}
	if request.Description != "" {
		funding.Description = request.Description
	}
	if request.Image != "" {
		funding.Image = resp.SecureURL
	}

	dataFuding, err := h.FundingRepository.UpdateFunding(funding, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: dataFuding}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerFunding) DeleteFunding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	funding, err := h.FundingRepository.GetFunding(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	dataFunding, err := h.FundingRepository.DeleteFunding(funding, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: dataFunding}
	json.NewEncoder(w).Encode(response)
}

func convertResponseFunding(u models.Funding) fundingdto.FundingResponse {
	return fundingdto.FundingResponse{
		ID:          u.ID,
		Title:       u.Title,
		Goals:       u.Goals,
		Description: u.Description,
		Image:       u.Image,
		UserID:      u.UserID,
		User:        u.User,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
