package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dgabrov/expense-go/internal/data"
	"github.com/google/uuid"
)

func writeJsonResponse(writer http.ResponseWriter, payload any, err error) {
	status := http.StatusOK
	processPayload := payload

	if err != nil {
		processPayload = data.ErrorResponse{
			Message: err.Error(),
			Items:   make([]string, 0),
		}

		status = http.StatusBadRequest
	}

	header := writer.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Cache-Control", "no-cache")

	writer.WriteHeader(status)

	// write the actual object here
	bts, err := json.Marshal(processPayload)
	if err != nil {
		slog.Error("error marshalling payload")
	}

	_, err = writer.Write(bts)
	if err != nil {
		slog.Error("error writing payload")
	}
}

func createRandomToken() string {
	u7 := uuid.Must(uuid.NewV7())
	randomBytes := make([]byte, 14)
	_, _ = rand.Read(randomBytes)
	return fmt.Sprintf("%s%s", u7.String(), hex.EncodeToString(randomBytes))[:64]
}
