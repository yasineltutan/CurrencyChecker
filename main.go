package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type normalResponse struct {
	Data     float64 `json:"data"`
	Currency string  `json:"currency"`
}

type errorResponse struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

func responseHandler(w http.ResponseWriter, status int, currency string, data float64) {
	var nRes normalResponse
	var eRes errorResponse
	if status == http.StatusOK {
		nRes.Data = data
		nRes.Currency = currency
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(nRes)
	} else if status == http.StatusBadRequest {
		eRes.Error = 400
		eRes.Message = "Bad Request"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(eRes)
	} else if status == http.StatusNotFound {
		eRes.Error = 404
		eRes.Message = "Not Found"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(eRes)
	} else {
		eRes.Error = 500
		eRes.Message = "Internal server error"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(eRes)
	}
}

func currencyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		switch r.URL.Path {
		case "/USD":
			data, err := getCurrency("USD")
			if err != nil {
				responseHandler(w, 500, "", 0.0)
			} else {
				responseHandler(w, 200, "USD", data)
			}
		case "/EUR":
			data, err := getCurrency("EUR")
			if err != nil {
				responseHandler(w, 500, "", 0.0)
			} else {
				responseHandler(w, 200, "EUR", data)
			}
		case "/GBP":
			data, err := getCurrency("GBP")
			if err != nil {
				responseHandler(w, 500, "", 0.0)
			} else {
				responseHandler(w, 200, "GBP", data)
			}
		default:
			responseHandler(w, 404, "", 0.0)
		}
	} else {
		responseHandler(w, 400, "", 0.0)
	}
}

func handleRequests() {
	http.HandleFunc("/", currencyHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func main() {
	handleRequests()
}
