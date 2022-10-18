package handlers

import (
	"encoding/json"
	"net/http"
)

type responseBody struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func respondWithJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var status string
	if statusCode >= 100 && statusCode < 400 {
		status = "success"
	} else {
		status = "error"
	}
	p := responseBody{
		Status:  status,
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(p)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, message, nil)
}
