package respond

import (
	"api/internal/store"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type errorBody struct {
	Message string `json:"error"`
}

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			slog.Error("respond: encode failed", "err", err)
		}
	}
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, errorBody{Message: message})
}

func FromError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		Error(w, http.StatusNotFound, "not found")
	case errors.Is(err, store.ErrConflict):
		Error(w, http.StatusConflict, "stale revision; re-fetch and retry")
	case errors.Is(err, store.ErrDuplicate):
		Error(w, http.StatusConflict, "already exists")
	case errors.Is(err, store.ErrInvalidInput):
		Error(w, http.StatusBadRequest, err.Error())
	default:
		slog.Error("respond: unexpected error", "err", err)
		Error(w, http.StatusInternalServerError, "internal server error")
	}
}
