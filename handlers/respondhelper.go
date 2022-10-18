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

type responseError struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

func respondWithJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var p interface{}
	if statusCode >= 100 && statusCode < 400 {
		p = responseBody{
			Status:  "success",
			Message: message,
			Data:    data,
		}
	} else {
		p = responseError{
			Status:  "error",
			Message: message,
			Error:   data,
		}
	}

	json.NewEncoder(w).Encode(p)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string, err error) {
	respondWithJSON(w, statusCode, message, err.Error())
}
