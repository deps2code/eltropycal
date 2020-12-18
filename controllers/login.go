package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eltropycal/models/request"
	"github.com/eltropycal/models/response"
	"github.com/eltropycal/utils"
	log "github.com/sirupsen/logrus"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var apiReq request.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&apiReq)
	if err != nil {
		log.WithFields(log.Fields{"api": "Login", "error": "invalid_api_request"}).Error(err.Error())
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err = apiReq.Validate(); err != nil {
		log.WithFields(log.Fields{"api": "Login", "error": "invalid_request_body"}).Error(err.Error())
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	//Get user details with username and role
	userDetails, err := DataService.GetUserDetailsByUsername(apiReq.Username, apiReq.Role)
	if err != nil {
		log.WithFields(log.Fields{"api": "Login", "error": "db_error"}).Error(err.Error())
		RespondWithError(w, http.StatusInternalServerError, "Opps, something went wrong.")
		return
	}

	//check if user is there in the db
	if userDetails.ID == "" {
		log.WithFields(log.Fields{"api": "Login", "error": "api_error"}).Error("user not found")
		RespondWithError(w, http.StatusBadRequest, "The user doesn't exist.")
		return
	}

	//compare password hash
	isUserAuthenticated, err := utils.ComparePasswordHash(userDetails.HashedPassword, apiReq.Password)
	if err != nil {
		log.WithFields(log.Fields{"api": "Login", "error": "module_error"}).Error(err.Error())
		RespondWithError(w, http.StatusInternalServerError, "Opps, something went wrong.")
		return
	}

	if !isUserAuthenticated {
		log.WithFields(log.Fields{"api": "Login", "user_id": userDetails.ID}).Info("Wrong password")
		RespondWithError(w, http.StatusInternalServerError, "Please type the correct password")
		return
	}

	//encode user data into an auth token
	authToken := utils.EncodeJwt(utils.JWTInputData{
		UserID: userDetails.ID,
		Role:   userDetails.Role,
	})

	log.WithFields(log.Fields{"api": "Login", "user_id": userDetails.ID}).Info("Login successfull")
	RespondWithSuccessGeneric(w, http.StatusOK, response.Response{
		Success: true,
		Data: response.Login{
			AuthToken: authToken,
			LoginAt:   time.Now().UTC().Format(time.RFC3339),
		}})
	return
}
