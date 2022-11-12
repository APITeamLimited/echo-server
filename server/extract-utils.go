package server

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"gopkg.in/guregu/null.v4"
)

func getAcceptType(req *http.Request) string {
	// Sllit the Accept header into an array of strings
	accept := strings.Split(req.Header.Get("Accept"), ",")

	// If the Accept header is empty, return the default content type
	if accept[0] == "" {
		return "text/html"
	}

	// If any item text/html, return text/html
	for _, item := range accept {
		if item == "text/html" {
			return "text/html"
		}
	}

	// If any item application/json, return application/json
	for _, item := range accept {
		if item == "application/json" {
			return "application/json"
		}
	}

	// For the rest of the cases, return empty string
	return ""
}

func extractQueryParams(req *http.Request) []InfoPair {
	// Extract and add URL parameters to responseInfo

	query := req.URL.Query()

	var infoPairs []InfoPair

	for key, value := range query {
		// Join the values into a single string, joining array of values together
		infoPairs = append(infoPairs, InfoPair{Key: key, Value: strings.Join(value, ",")})
	}

	// Sort pairs alphabetically by key
	sort.Slice(infoPairs, func(i, j int) bool {
		return infoPairs[i].Key < infoPairs[j].Key
	})

	return infoPairs
}

func extractHeaders(req *http.Request) []InfoPair {
	// Extract and add headers to responseInfo

	var infoPairs []InfoPair

	for key, value := range req.Header {
		// Join the values into a single string, joining array of values together
		infoPairs = append(infoPairs, InfoPair{Key: key, Value: strings.Join(value, ",")})
	}

	// Sort pairs alphabetically by key
	sort.Slice(infoPairs, func(i, j int) bool {
		return infoPairs[i].Key < infoPairs[j].Key
	})

	return infoPairs
}

func extractCookies(req *http.Request) []InfoPair {
	// Extract and add cookies to responseInfo

	var infoPairs []InfoPair

	for _, cookie := range req.Cookies() {
		infoPairs = append(infoPairs, InfoPair{Key: cookie.Name, Value: cookie.Value})
	}

	// Sort pairs alphabetically by key
	sort.Slice(infoPairs, func(i, j int) bool {
		return infoPairs[i].Key < infoPairs[j].Key
	})

	return infoPairs
}

func extractBody(req *http.Request) (null.String, error) {
	// Extract and add body to responseInfo

	// Limit body size to 500kb and read it
	bodyBytes, err := ioutil.ReadAll(http.MaxBytesReader(nil, req.Body, 500*1024))
	if err != nil {
		return null.String{
			NullString: sql.NullString{
				String: "",
				Valid:  false},
		}, err
	}

	stringValue := string(bodyBytes)

	return null.String{
		NullString: sql.NullString{
			String: stringValue,
			Valid:  stringValue != "",
		},
	}, nil
}
