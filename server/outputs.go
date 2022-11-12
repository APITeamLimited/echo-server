package server

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func outputJSON(w http.ResponseWriter, requestInfo RequestInfo) {
	w.Header().Set("Content-Type", "application/json")

	marshalled, err := json.Marshal(&struct {
		PreSummary    string            `json:"preSummary"`
		Request       map[string]string `json:"request"`
		URLParameters map[string]string `json:"urlParameters"`
		Headers       map[string]string `json:"headers"`
		Cookies       map[string]string `json:"cookies"`
		Body          string            `json:"body"`
	}{
		PreSummary:    requestInfo.PreSummary,
		Request:       convertToKeyValueMap(requestInfo.Request),
		URLParameters: convertToKeyValueMap(requestInfo.URLParameters),
		Headers:       convertToKeyValueMap(requestInfo.Headers),
		Cookies:       convertToKeyValueMap(requestInfo.Cookies),
		Body:          requestInfo.Body.String,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(marshalled)
}

var htmlTemplate = template.Must(template.ParseFiles("server/response-template.html"))

func outputHTML(w http.ResponseWriter, requestInfo RequestInfo) {
	w.Header().Set("Content-Type", "text/html")
	err := htmlTemplate.Execute(w, requestInfo)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func convertToKeyValueMap(infoPairs []InfoPair) map[string]string {
	// Convert infoPairs to map[string]string

	// Initialize map[string]string
	m := make(map[string]string)

	// Loop through infoPairs
	for _, pair := range infoPairs {
		// Add key-value pair to map
		m[pair.Key] = pair.Value
	}

	// Return map
	return m
}
