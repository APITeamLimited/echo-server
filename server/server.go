package server

import (
	"fmt"
	"net/http"

	"gopkg.in/guregu/null.v4"
)

type InfoPair struct {
	Key   string
	Value string
}

type RequestInfo struct {
	PreSummary    string      `json:"preSummary"`
	Request       []InfoPair  `json:"request"`
	URLParameters []InfoPair  `json:"urlParameters"`
	Headers       []InfoPair  `json:"headers"`
	Cookies       []InfoPair  `json:"cookies"`
	Body          null.String `json:"body"`
}

func determineRequestInfo(req *http.Request) (RequestInfo, error) {
	requestInfo := RequestInfo{
		PreSummary: fmt.Sprintf("%s / %s", req.Method, req.Proto),
	}

	requestInfo.Request = []InfoPair{
		{Key: "Method", Value: req.Method},
		{Key: "URL", Value: req.URL.String()},
		{Key: "Proto", Value: req.Proto},
		{Key: "Host", Value: req.Host},
	}

	requestInfo.URLParameters = extractQueryParams(req)

	requestInfo.Headers = extractHeaders(req)

	requestInfo.Cookies = extractCookies(req)

	body, err := extractBody(req)
	if err != nil {
		return RequestInfo{}, err
	}

	requestInfo.Body = body

	return requestInfo, nil
}

func Run() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8080", nil)
}

func handle(w http.ResponseWriter, req *http.Request) {
	accept := getAcceptType(req)

	if accept == "application/json" || accept == "text/html" || accept == "*/*" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	requestInfo, err := determineRequestInfo(req)
	if err != nil {
		// If error says that body is too large, return 413
		if err.Error() == "http: request body too large" {
			http.Error(w, "Request body too large, max size is 500kb", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if accept == "application/json" {
		outputJSON(w, requestInfo)
	} else if accept == "text/html" || accept == "*/*" {
		outputHTML(w, requestInfo)
	}
}
