package controllers

import (
	"net/http"

	"github.com/eltropycal/utils"
	"github.com/floyd/controllers"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
)

func UserLogout(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "userContext")
	err := DataService.InvalidateUserSession(user.(utils.UserJWTContext).UserID)
	if err != nil {
		log.WithFields(log.Fields{"api": "Login", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		RespondWithError(w, http.StatusInternalServerError, "Opps, something went wrong.")
		return
	}
	controllers.RespondWithSuccess(w, http.StatusOK, "Successfully logged out")
	return
}
