package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/eltropycal/models/request"
	log "github.com/sirupsen/logrus"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var apiReq request.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&apiReq)
	if err != nil {
		log.WithFields(log.Fields{"api": "Login", "error": "Invalid request"}).Error(err.Error())
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err = apiReq.Validate(); err != nil {
		log.WithFields(log.Fields{"api": "Login", "error": "Invalid request body"}).Error(err.Error())
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
}
