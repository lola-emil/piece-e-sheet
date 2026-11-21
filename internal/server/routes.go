package server

import (
	"api/internal/account"
	"api/internal/auth"
	"api/internal/category"
	"api/internal/debt"
	"api/internal/expense"
	"api/internal/payment"
	"api/internal/user"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes(
	expenseRepo expense.ExpenseRepository,
	categoryRepo category.CategoryRepository,
	authRepo auth.AuthRepository,
	debtRepo debt.DebtRepository,
	paymentRepo payment.PaymentRepository,
	accountRepo account.AccountRepository,
	userRepo user.UserRepository,
) http.Handler {
	originsEnv := os.Getenv("FRONTEND_URLS")
	var allowedOrigins []string

	if originsEnv != "" {
		allowedOrigins = strings.Split(originsEnv, ",")
	} else {
		allowedOrigins = []string{
			"http://localhost:5173", // Vite default
		}
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// SERVICES
	categoryService := category.NewCategoryService(categoryRepo)
	expenseService := expense.NewExpenseService(expenseRepo)
	authService := auth.NewAuthService(authRepo, userRepo)
	debtService := debt.NewDebtService(debtRepo)
	paymentService := payment.NewPaymentService(paymentRepo)
	accountService := account.NewAccountService(accountRepo)

	// HANDLERS
	categoryHandler := category.NewCategoryHandler(categoryService)
	expenseHandler := expense.NewExpenseHandler(expenseService)
	authHandler := auth.NewAuthHandler(authService)
	debtHandler := debt.NewDebtHandler(debtService)
	paymentHandler := payment.NewPaymentHandler(paymentService)
	accountHandler := account.NewAccountHandler(accountService)

	auth.RegisterRoutes(r, authHandler)

	r.Route("/api", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		category.RegisterRoutes(r, categoryHandler)
		expense.RegisterRoutes(r, expenseHandler)
		debt.RegisterRoutes(r, debtHandler)
		payment.RegisterRoutes(r, paymentHandler)
		account.RegisterRoutes(r, accountHandler)
	})

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	if err := chi.Walk(r, func(
		method string,
		route string,
		handler http.Handler,
		middlewares ...func(http.Handler) http.Handler,
	) error {
		log.Printf("ROUTE: %s %s", method, route)
		return nil
	}); err != nil {
		log.Printf("Route inspection error: %v", err)
	}

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
