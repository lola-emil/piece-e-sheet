package debt

import (
	"api/internal/auth"
	"api/internal/respond"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct{ service DebtService }

func NewDebtHandler(s DebtService) *handler { return &handler{s} }

func (h *handler) FindAll(w http.ResponseWriter, r *http.Request) {
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	data, err := h.service.FindAll(r.Context(), uid)
	if err != nil {
		respond.Error(w, 500, err.Error())
		return
	}
	respond.JSON(w, 200, data)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req CreateDebtRequest
	json.NewDecoder(r.Body).Decode(&req)

	data, err := h.service.Create(r.Context(), uid, &req)
	if err != nil {
		respond.Error(w, 400, err.Error())
		return
	}
	respond.JSON(w, 201, data)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req UpdateDebtRequest
	json.NewDecoder(r.Body).Decode(&req)

	data, err := h.service.Update(r.Context(), id, uid, &req)
	if err != nil {
		respond.Error(w, 400, err.Error())
		return
	}
	respond.JSON(w, 200, data)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.service.Delete(r.Context(), id, uid); err != nil {
		respond.Error(w, 404, err.Error())
		return
	}
	w.WriteHeader(200)
}
