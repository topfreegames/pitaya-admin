package api

import (
	"encoding/json"
	"net/http"
)

//Write to response with status code
func Write(w http.ResponseWriter, status int, text string) {
	WriteBytes(w, status, []byte(text))
}

//WriteJSON to the response and with the status code
func WriteJSON(w http.ResponseWriter, status int, body map[string]interface{}) {
	bts, _ := json.Marshal(body)
	WriteBytes(w, status, bts)
}

//WriteBytes to the response and with the status code
func WriteBytes(w http.ResponseWriter, status int, text []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(text)
}
