package apperror

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		slog.Error("request error", "message", appErr.Message, "internal", appErr.err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.StatusCode)
		json.NewEncoder(w).Encode(appErr)
		return
	}
	slog.Error("unexpected error", "error", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
	})
}
