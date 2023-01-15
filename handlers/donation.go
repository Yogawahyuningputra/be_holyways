package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	donationdto "server/dto/donation"
	dto "server/dto/result"
	"server/models"
	"server/repositories"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
)

var c = coreapi.Client{
	ServerKey: os.Getenv("SERVER_KEY"),
	ClientKey: os.Getenv("CLIENT_KEY"),
}

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
	//create unique idtransaction
	var transIDisMatch = false
	var donationID int

	for !transIDisMatch {
		donationID = int(time.Now().Unix())
		request, _ := h.DonationRepository.GetDonation(donationID)
		if request.ID == 0 {
			transIDisMatch = true
		}

	}
	funding_id, _ := strconv.Atoi(r.FormValue("funding_id"))
	money, _ := strconv.Atoi(r.FormValue("money"))

	request := donationdto.DonationRequest{
		FundingID: funding_id,
		Money:     money,
		Status:    r.FormValue("status"),
	}
	// log.Print("this is request =>", request)

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	donation := models.Donation{
		ID:        donationID, // unixID
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
	// fmt.Println("this get donationbyID", donation)

	// request token from midtranss
	// initial snap client
	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(donation.ID),
			GrossAmt: int64(donation.Money),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: donation.User.Fullname,
			Email: donation.User.Email,
		},
	}
	snapResp, _ := s.CreateTransaction(req)
	// fmt.Println("this snap req =>", req)
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: snapResp}
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
		User:      u.User,
		FundingID: u.FundingID,
		Funding:   u.Funding,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (h *handlerDonation) GetDonationByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	donations, err := h.DonationRepository.GetDonationByUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: donations}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerDonation) GetDonationByFunding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	donations, err := h.DonationRepository.GetDonationByFunding(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: donations}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerDonation) GetDonationPending(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	donations, err := h.DonationRepository.GetDonationPending(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: donations}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerDonation) Notification(w http.ResponseWriter, r *http.Request) {
	var notificationPayload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "error di sini"}
		json.NewEncoder(w).Encode(response)
		return
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderId := notificationPayload["order_id"].(string)

	// fmt.Println("this order id", orderId)

	IDtrans, _ := strconv.Atoi(orderId)

	donation, _ := h.DonationRepository.GetDonation(IDtrans)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			// TODO set transaction status on your database to 'challenge'
			// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			SendMail("challenge", donation)
			h.DonationRepository.UpdateStatus("pending", donation.ID)

		} else if fraudStatus == "accept" {
			// TODO set transaction status on your database to 'success'
			SendMail("approve", donation)
			h.DonationRepository.UpdateStatus("success", donation.ID)
		}
	} else if transactionStatus == "settlement" {
		// TODO set transaction status on your database to 'success'
		SendMail("success", donation)
		h.DonationRepository.UpdateStatus("success", donation.ID)

	} else if transactionStatus == "deny" {
		// TODO you can ignore 'deny', because most of the time it allows payment retries
		// and later can become success
		SendMail("failed", donation)
		h.DonationRepository.UpdateStatus("failed", donation.ID)

	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		// TODO set transaction status on your database to 'failure'
		SendMail("cancel", donation)
		h.DonationRepository.UpdateStatus("failed", donation.ID)

	} else if transactionStatus == "pending" {
		// TODO set transaction status on your database to 'pending' / waiting payment
		SendMail("pending", donation)
		h.DonationRepository.UpdateStatus("pending", donation.ID)
	}

	w.WriteHeader(http.StatusOK)
}

func SendMail(status string, donation models.Donation) {

	var CONFIG_SMTP_HOST = "smtp.gmail.com"
	var CONFIG_SMTP_PORT = 587
	var CONFIG_SENDER_NAME = "Dewetour <demo.dumbways@gmail.com>"
	var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
	var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")

	var Funding = donation.Funding.Title
	var money = strconv.Itoa(donation.Money)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", donation.User.Email)
	mailer.SetHeader("Subject", "Transaction Status")
	mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
	  <html lang="en">
		<head>
		<meta charset="UTF-8" />
		<meta http-equiv="X-UA-Compatible" content="IE=edge" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Document</title>
		<style>
		  h1 {
		  color: brown;
		  }
		</style>
		</head>
		<body>
		<h2>Transaction :</h2>
		<ul style="list-style-type:none;">
		  <li>Name : %s</li>
		  <li>Total payment: Rp.%s</li>
		  <li>Status : <b>%s</b></li>
		</ul>
		</body>
	  </html>`, Funding, money, status))

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Mail sent! to " + donation.User.Email)
}
