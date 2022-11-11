package server

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func outputJSON(w http.ResponseWriter, requestInfo RequestInfo) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestInfo)
}

var htmlTemplate = template.Must(template.ParseFiles("server/response-template.html"))

func outputHTML(w http.ResponseWriter, requestInfo RequestInfo) {
	w.Header().Set("Content-Type", "text/html")
	htmlTemplate.Execute(w, requestInfo)
}
