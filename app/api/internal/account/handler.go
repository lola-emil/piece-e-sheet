package account

import (
	"api/internal/auth"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Helper to write JSON responses
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

type handler struct {
	service AccountService
}

func NewAccountHandler(service AccountService) *handler {
	return &handler{service: service}
}

// Helper to extract user ID from context (will be set by Auth middleware later)
func getUserID(r *http.Request) string {
	// Use auth.UserIDKey instead of the local userIDKey
	if uid, ok := r.Context().Value(auth.UserIDKey).(string); ok {
		return uid
	}
	return ""
}

func (h *handler) FindAll(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	accounts, err := h.service.FindAll(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"data": accounts})
}

func (h *handler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := getUserID(r)

	account, err := h.service.FindByID(r.Context(), id, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"data": account})
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	account, err := h.service.Create(r.Context(), userID, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, map[string]any{"data": account})
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := getUserID(r)

	var req UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	account, err := h.service.Update(r.Context(), id, userID, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"data": account})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := getUserID(r)
	err := h.service.Delete(r.Context(), id, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "account deleted"})
}
