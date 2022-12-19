package server

import (
	"fmt"
	"net/http"

	"gopkg.in/guregu/null.v4"
)

const Port = 8080

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

func Run() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(fmt.Sprintf(":%d", Port), nil)
}

func determineRequestInfo(req *http.Request) (RequestInfo, error) {
	requestInfo := RequestInfo{
		PreSummary: fmt.Sprintf("%s / %s", req.Method, req.Proto),
	}

	requestInfo.Request = []InfoPair{
		{Key: "Method", Value: req.Method},
		{Key: "URL", Value: req.URL.String()},
		{Key: "Proto", Value: req.Proto},
	}

	connectingIP := req.Header.Get("Cf-Connecting-Ip")
	if req.Header.Get("Cf-Connecting-Ip") != "" {
		requestInfo.Request = append(requestInfo.Request, InfoPair{Key: "Client IP", Value: connectingIP})
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

func handle(w http.ResponseWriter, req *http.Request) {
	accept := getAcceptType(req)

	if accept == "" {
		http.Error(w, "406 Not Acceptable", http.StatusNotAcceptable)
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
