package expense

import (
	"api/internal/auth"
	"encoding/json"
	"net/http"
	"strconv"
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
	if uid, ok := r.Context().Value(auth.UserIDKey).(string); ok {
		return uid
	}
	return ""
}

func parseDateParam(r *http.Request, param string, endOfDay bool) *time.Time {
	val := r.URL.Query().Get(param)
	if val == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", val)
	if err != nil {
		if t2, err2 := time.Parse(time.RFC3339, val); err2 == nil {
			return &t2
		}
		return nil
	}

	if endOfDay {
		t = t.Add(24*time.Hour - time.Nanosecond)
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

	filter := ExpenseFilter{}

	if v := r.URL.Query().Get("category_id"); v != "" {
		filter.CategoryID = &v
	}
	if v := r.URL.Query().Get("account_id"); v != "" {
		filter.AccountID = &v
	}

	filter.StartDate = parseDateParam(r, "start_date", false)
	filter.EndDate = parseDateParam(r, "end_date", true)

	filter.MinAmount = parseFloatParam(r, "min_amount")
	filter.MaxAmount = parseFloatParam(r, "max_amount")

	filter.SortBy = r.URL.Query().Get("sort_by")

	filter.Limit = parseIntParam(r, "limit", 0)
	filter.Offset = parseIntParam(r, "offset", 0)

	expenses, count, err := h.service.FindAll(r.Context(), userID, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data":  expenses,
		"count": count,
	})
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

func parseFloatParam(r *http.Request, key string) *float64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil
	}
	return &f
}

func parseIntParam(r *http.Request, key string, defaultValue int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return num
}
