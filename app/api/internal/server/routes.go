package server

import (
	"api/internal/auth"
	"api/internal/category"
	"api/internal/expense"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes(
	expenseRepo expense.ExpenseRepository,
	categoryRepo category.CategoryRepository,
	authRepo auth.AuthRepository,
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// SERVICES
	categoryService := category.NewCategoryService(categoryRepo)
	expenseService := expense.NewExpenseService(expenseRepo)
	authService := auth.NewAuthService(authRepo)

	// HANDLERS
	categoryHandler := category.NewCategoryHandler(categoryService)
	expenseHandler := expense.NewExpenseHandler(expenseService)
	authHandler := auth.NewAuthHandler(authService)

	auth.RegisterRoutes(r, authHandler)

	r.Route("/api", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		category.RegisterRoutes(r, categoryHandler)
		expense.RegisterRoutes(r, expenseHandler)
	})

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

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
