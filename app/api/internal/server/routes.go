package server

import (
	"api/internal/appmiddleware"
	"api/internal/expense"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

func (s *Server) RegisterRoutes(
	expenseRepo expense.ExpenseRepository,
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
	expenseService := expense.NewExpenseService(expenseRepo)

	// HANDLERS
	expenseHandler := expense.NewExpenseHandler(expenseService)

	developmentUserID, err := uuid.Parse(os.Getenv("DEVELOPMENT_USER_ID"))

	if err != nil {
		log.Fatalf("Failed to parse UUID string: %v", err)
	}

	r.Route("/api", func(r chi.Router) {
		r.Use(appmiddleware.DevelopmentIdentity(
			developmentUserID,
		))

		r.Route("/expenses", func(r chi.Router) {
			expense.RegisterRoutes(r, &expenseHandler)
		})
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
