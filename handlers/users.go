package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	dto "server/dto/result"
	userdto "server/dto/user"
	"server/pkg/bcrypt"

	"net/http"
	"server/models"
	"server/repositories"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerUser struct {
	UserRepository repositories.UserRepository
}

func HandlerUser(UserRepository repositories.UserRepository) *handlerUser {
	return &handlerUser{UserRepository}
}

func (h *handlerUser) FindUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	users, err := h.UserRepository.FindUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}

		json.NewEncoder(w).Encode(response)

	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: users}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := h.UserRepository.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(user)}
	json.NewEncoder(w).Encode(response)
}
func convertResponse(u models.User) userdto.UserResponse {
	return userdto.UserResponse{
		ID:       u.ID,
		Fullname: u.Fullname,
		Email:    u.Email,
		Password: u.Password,
		Phone:    u.Phone,
		Gender:   u.Gender,
		Address:  u.Address,
		Image:    u.Image,
		Role:     u.Role,
	}
}

func (h *handlerUser) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	phone, _ := strconv.Atoi(r.FormValue("phone"))
	request := userdto.CreateUserRequest{
		Fullname: r.FormValue("fullname"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Gender:   r.FormValue("gender"),
		Phone:    phone,
		Address:  r.FormValue("address"),
		Role:     r.FormValue("role"),
	}
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	password, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}
	user := models.User{
		Fullname: request.Fullname,
		Email:    request.Email,
		Password: password,
		Gender:   request.Gender,
		Phone:    request.Phone,
		Address:  request.Address,
		Role:     "user",
		Image:    request.Image,
	}

	data, err := h.UserRepository.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	dataContex := r.Context().Value("dataFile")
	filepath := dataContex.(string)

	phone, _ := strconv.Atoi(r.FormValue("phone"))
	request := userdto.CreateUserRequest{
		Fullname: r.FormValue("fullname"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Gender:   r.FormValue("gender"),
		Phone:    phone,
		Address:  r.FormValue("address"),
		Image:    filepath,
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := h.UserRepository.GetUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add Cloudinary credentials
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "holyways/profile"})

	if err != nil {
		fmt.Println(err.Error())
	}
	password, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	if request.Fullname != "" {
		user.Fullname = request.Fullname
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Password != "" {
		user.Password = password
	}
	if request.Gender != "" {
		user.Gender = request.Gender
	}
	if request.Phone != 0 {
		user.Phone = request.Phone
	}
	if request.Address != "" {
		user.Address = request.Address
	}
	if request.Image != "" {
		user.Image = resp.SecureURL
	}
	data, err := h.UserRepository.UpdateUser(user, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := h.UserRepository.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.UserRepository.DeleteUser(user, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(data)}
	json.NewEncoder(w).Encode(response)

}
