package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/topfreegames/pitaya/logger"
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

// WriteError to the response with message and log error message
func WriteError(w http.ResponseWriter, status int, errorMsg string, err error) {
	logger.Log.Errorf("error %s, %s", errorMsg, err.Error())
	errMsg := fmt.Sprintf(`{"success":false, "message":"%s", "reason": "%s"}`, errorMsg, err.Error())
	Write(w, status, errMsg)
}

// WriteSuccessWithJSON sends response with statusOK to request and log success msg
func WriteSuccessWithJSON(w http.ResponseWriter, status int, res []byte, msg string) {
	retMsg := fmt.Sprintf(`{"success" : true, "response":%s}`, res)
	logger.Log.Info(msg)
	Write(w, status, retMsg)
}
