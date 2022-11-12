package server

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func outputJSON(w http.ResponseWriter, requestInfo RequestInfo) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(requestInfo)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

var htmlTemplate = template.Must(template.ParseFiles("server/response-template.html"))

func outputHTML(w http.ResponseWriter, requestInfo RequestInfo) {
	w.Header().Set("Content-Type", "text/html")
	err := htmlTemplate.Execute(w, requestInfo)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
