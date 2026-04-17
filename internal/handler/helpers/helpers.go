package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/dto"
)

func WriteJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(dto.ErrorResponse{
		Status:  "error",
		Message: message,
	})
	if err != nil {
		return
	}
}

func WriteJSONSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	if data == nil {
		WriteJSONError(w, http.StatusInternalServerError, "no data to return")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func IsAlpha(name string) bool {
	return regexp.MustCompile("^[A-Za-z ]+$").MatchString(name)
}
