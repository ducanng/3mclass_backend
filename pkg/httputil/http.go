package httputil

import (
	"encoding/json"
	"net/http"
)

const Version = "v1"

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteJSONResponse(w http.ResponseWriter, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func WriteJSONMessage(w http.ResponseWriter, httpCode int, message string) {
	if httpCode != http.StatusOK {
		w.WriteHeader(httpCode)
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(Response{Code: httpCode, Message: message})

	w.Write(response)
}
