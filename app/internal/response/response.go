package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiErr struct {
	Message error
}

var (
	EnvErr           = New("cannot get environment")
	ErrIDNotProvided = New("id not provided")
	ErrNoSuchTask    = New("no task with provided id")
	ErrParseTime     = New("error parsing time")
)

func New(err string) apiErr {
	return apiErr{
		Message: fmt.Errorf(err),
	}
}

func (err *apiErr) Error(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Message.Error()}); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

type Response struct {
	MessageKey string
	MessageVal string
	W          http.ResponseWriter
	RespCode   int
}

func WriteResponse(
	messageKey string,
	messageVal string,
	w http.ResponseWriter,
	respCode int) {
	w.WriteHeader(respCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(map[string]string{messageKey: messageVal}); err != nil {
		http.Error(w, "error encoding response", respCode)
		return
	}
}
