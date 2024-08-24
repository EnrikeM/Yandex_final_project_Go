package apierrors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiErr struct {
	Message error
}

var (
	EnvErr = New(("cannot get environment"))
)

func New(err string) apiErr {
	return apiErr{
		Message: fmt.Errorf(err),
	}
}

func (err *apiErr) Error(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Message.Error()}); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

}
