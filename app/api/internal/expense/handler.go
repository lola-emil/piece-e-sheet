package expense

import (
	"api/internal/auth"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

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

func getUserID(r *http.Request) string {
	// Use auth.UserIDKey instead of the local userIDKey
	if uid, ok := r.Context().Value(auth.UserIDKey).(string); ok {
		return uid
	}
	return ""
}

// Helper to parse optional time from query params
func parseTimeParam(r *http.Request, param string) *time.Time {
	val := r.URL.Query().Get(param)
	if val == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return nil
	}
	return &t
}

type handler struct {
	service ExpenseService
}

func NewExpenseHandler(service ExpenseService) *handler {
	return &handler{service: service}
}

func (h *handler) FindAll(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	// Parse filters
	filter := ExpenseFilter{}
	catID := r.URL.Query().Get("category_id")
	if catID != "" {
		filter.CategoryID = &catID
	}
	filter.StartDate = parseTimeParam(r, "start_date")
	filter.EndDate = parseTimeParam(r, "end_date")

	expenses, err := h.service.FindAll(r.Context(), userID, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"data": expenses})
}

func (h *handler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := getUserID(r)

	expense, err := h.service.FindByID(r.Context(), id, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"data": expense})
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	var req CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	expense, err := h.service.Create(r.Context(), userID, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": expense})
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := getUserID(r)

	var req UpdateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	expense, err := h.service.Update(r.Context(), id, userID, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"data": expense})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := getUserID(r)

	err := h.service.Delete(r.Context(), id, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "expense deleted"})
}
