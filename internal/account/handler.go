package account

import (
	"api/internal/auth"
	"api/internal/respond"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	service AccountService
}

func NewAccountHandler(service AccountService) *handler {
	return &handler{service: service}
}

func (h *handler) FindAll(w http.ResponseWriter, r *http.Request) {
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	accounts, err := h.service.FindAll(r.Context(), uid)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"data": accounts})
}

func (h *handler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	account, err := h.service.FindByID(r.Context(), id, uid)
	if err != nil {
		respond.Error(w, http.StatusNotFound, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"data": account})
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	account, err := h.service.Create(r.Context(), uid, &req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusCreated, map[string]any{"data": account})
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, ok := auth.UserIDFromRequest(r)

	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	account, err := h.service.Update(r.Context(), id, uid, &req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"data": account})
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
	respond.JSON(w, http.StatusOK, map[string]string{"message": "account deleted"})
}
