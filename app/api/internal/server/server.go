package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"api/internal/auth"
	"api/internal/category"
	"api/internal/database"
	"api/internal/expense"
)

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,

		db: database.New(),
	}

	dbInstance := NewServer.db.GetInstance()

	categoryRepo := category.NewCategoryRepository(dbInstance)
	expenseRepo := expense.NewExpenseRepository(dbInstance)
	authRepo := auth.NewAuthRepository(dbInstance)

	// Declare Server config
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(
			expenseRepo,
			categoryRepo,
			authRepo,
		),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
