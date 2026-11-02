package payment

import (
	"api/internal/auth"
	"api/internal/respond"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct{ service PaymentService }

func NewPaymentHandler(s PaymentService) *handler { return &handler{s} }

func (h *handler) GetByDebtID(w http.ResponseWriter, r *http.Request) {
	debtID := chi.URLParam(r, "debt_id")
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	data, err := h.service.GetByDebtID(r.Context(), debtID, uid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	respond.JSON(w, 200, data)
}

func (h *handler) Record(w http.ResponseWriter, r *http.Request) {
	uid, ok := auth.UserIDFromRequest(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req CreatePaymentRequest
	json.NewDecoder(r.Body).Decode(&req)

	data, err := h.service.Record(r.Context(), uid, &req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	respond.JSON(w, 201, data)
}
