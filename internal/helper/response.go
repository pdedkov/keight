package helper

import (
	"encoding/json"
	"net/http"
)

// WrapError wrap json http error
func WrapError(w http.ResponseWriter, err error, status int) {
	errString, _ := json.Marshal(ErrorResponse{
		Message: err.Error(),
	})
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(errString)
}
