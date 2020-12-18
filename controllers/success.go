package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/eltropycal/models/response"
)

func RespondWithSuccess(w http.ResponseWriter, code int, message string) {
	data, _ := json.Marshal(response.Response{Success: true, Data: response.Data{Message: message}})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithSuccessGeneric(w http.ResponseWriter, code int, body interface{}) {
	data, _ := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
