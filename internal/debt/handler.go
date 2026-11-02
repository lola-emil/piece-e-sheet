// debt/handler.go
package debt

import (
	"api/internal/auth"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": payload})
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

type handler struct{ service DebtService }

func NewDebtHandler(s DebtService) *handler { return &handler{s} }

func (h *handler) FindAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	data, err := h.service.FindAll(r.Context(), userID)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, data)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	var req CreateDebtRequest
	json.NewDecoder(r.Body).Decode(&req)

	data, err := h.service.Create(r.Context(), userID, &req)
	if err != nil {
		respondError(w, 400, err.Error())
		return
	}
	respondJSON(w, 201, data)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value(auth.UserIDKey).(string)
	var req UpdateDebtRequest
	json.NewDecoder(r.Body).Decode(&req)

	data, err := h.service.Update(r.Context(), id, userID, &req)
	if err != nil {
		respondError(w, 400, err.Error())
		return
	}
	respondJSON(w, 200, data)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value(auth.UserIDKey).(string)
	if err := h.service.Delete(r.Context(), id, userID); err != nil {
		respondError(w, 404, err.Error())
		return
	}
	w.WriteHeader(200)
}
