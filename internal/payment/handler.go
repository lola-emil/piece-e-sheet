// payment/handler.go
package payment

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

type handler struct{ service PaymentService }

func NewPaymentHandler(s PaymentService) *handler { return &handler{s} }

func (h *handler) GetByDebtID(w http.ResponseWriter, r *http.Request) {
	debtID := chi.URLParam(r, "debt_id")
	userID := r.Context().Value(auth.UserIDKey).(string)

	data, err := h.service.GetByDebtID(r.Context(), debtID, userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	respondJSON(w, 200, data)
}

func (h *handler) Record(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	var req CreatePaymentRequest
	json.NewDecoder(r.Body).Decode(&req)

	data, err := h.service.Record(r.Context(), userID, &req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	respondJSON(w, 201, data)
}
