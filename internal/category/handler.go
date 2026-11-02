package category

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
	service CategoryService
}

func NewCategoryHandler(service CategoryService) *handler {
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
	categories, err := h.service.FindAll(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"data": categories})
}

func (h *handler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	category, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"data": category})
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category, err := h.service.Create(r.Context(), userID, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": category})
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"data": category})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.service.Delete(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "category deleted"})
}
