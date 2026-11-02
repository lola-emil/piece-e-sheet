package expense

import (
	"api/internal/auth"
	"api/internal/respond"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

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
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	filter := ExpenseFilter{}

	if v := r.URL.Query().Get("category_id"); v != "" {
		filter.CategoryID = &v
	}
	if v := r.URL.Query().Get("account_id"); v != "" {
		filter.AccountID = &v
	}

	filter.Search = strings.TrimSpace(r.URL.Query().Get("search"))
	filter.StartDate = parseDateParam(r, "start_date", false)
	filter.EndDate = parseDateParam(r, "end_date", true)
	filter.MinAmount = parseFloatParam(r, "min_amount")
	filter.MaxAmount = parseFloatParam(r, "max_amount")
	filter.SortBy = r.URL.Query().Get("sort_by")

	filter.NoPagination =
		r.URL.Query().Get("no_pagination") == "true"

	if !filter.NoPagination {
		filter.Limit = parseIntParam(r, "limit", 0)
		filter.Offset = parseIntParam(r, "offset", 0)
	}

	expenses, count, err := h.service.FindAll(r.Context(), uid, filter)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"data":  expenses,
		"count": count,
	})
}

func (h *handler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	expense, err := h.service.FindByID(r.Context(), id, uid)
	if err != nil {
		respond.Error(w, http.StatusNotFound, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]interface{}{"data": expense})
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	expense, err := h.service.Create(r.Context(), uid, &req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusCreated, map[string]interface{}{"data": expense})
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req UpdateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	expense, err := h.service.Update(r.Context(), id, uid, &req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]interface{}{"data": expense})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.service.Delete(r.Context(), id, uid)
	if err != nil {
		respond.Error(w, http.StatusNotFound, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{"message": "expense deleted"})
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

func (h *handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	filter := ExpenseFilter{}

	if v := r.URL.Query().Get("category_id"); v != "" {
		filter.CategoryID = &v
	}
	if v := r.URL.Query().Get("account_id"); v != "" {
		filter.AccountID = &v
	}

	filter.Search = strings.TrimSpace(r.URL.Query().Get("search"))
	filter.StartDate = parseDateParam(r, "start_date", false)
	filter.EndDate = parseDateParam(r, "end_date", true)
	filter.MinAmount = parseFloatParam(r, "min_amount")
	filter.MaxAmount = parseFloatParam(r, "max_amount")
	filter.SortBy = r.URL.Query().Get("sort_by")

	summary, err := h.service.GetSummary(r.Context(), uid, filter)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, summary)
}
