package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/eltropycal/models/response"
)

// WebServerHealthCheck API
func WebServerHealthCheck(w http.ResponseWriter, r *http.Request) {
	var result = response.HealthCheckResponse{}
	result.Healthy = true
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(result)
	w.Write(data)
	return
}
