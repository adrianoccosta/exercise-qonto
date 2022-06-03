package tools

import (
	"encoding/json"
	"net/http"
)

const (
	// HeaderContentType defines the content-type header
	HeaderContentType = "content-type"

	// HeaderXPartnerURN defines the x-partner-urn header
	HeaderXPartnerURN = "x-partner-urn"

	jsonContentType = "application/json; charset=utf-8"
)

func WriteJSON(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.Header().Set(HeaderContentType, jsonContentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
}
