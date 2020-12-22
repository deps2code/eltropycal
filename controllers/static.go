package controllers

import "net/http"

// Static renders all static files
func Static(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
}
