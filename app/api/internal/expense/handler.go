package expense

import (
	"api/internal/requestctx"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ExpenseHandler struct {
	expenseService ExpenseService
}

func NewExpenseHandler(expenseService ExpenseService) ExpenseHandler {
	return ExpenseHandler{
		expenseService: expenseService,
	}
}

// TODO: filter by user
func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := requestctx.UserID(r.Context())

	expenses, err := h.expenseService.FindAll(r.Context(), userID)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(expenses)
}

func (h *ExpenseHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	expense, err := h.expenseService.FindByID(r.Context(), id)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(expense)
}
