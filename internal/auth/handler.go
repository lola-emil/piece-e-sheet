package auth

import (
	"api/internal/respond"
	"encoding/json"
	"net/http"
)

type handler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *handler {
	return &handler{service: service}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.DisplayName == "" {
		respond.Error(w, http.StatusBadRequest, "All fields are required")
		return
	}

	res, err := h.service.Register(r.Context(), &req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.JSON(w, http.StatusCreated, res)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		respond.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	res, err := h.service.Login(r.Context(), &req)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, res)
}
