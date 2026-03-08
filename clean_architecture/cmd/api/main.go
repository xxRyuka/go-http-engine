package main

import (
	"clean_architecture/internal/handler"
	"clean_architecture/internal/middleware"
	"clean_architecture/internal/repository"
	"clean_architecture/internal/service"
	"fmt"
	"net/http"
	"time"
)

func main() {
	app := handler.Application{AppName: "clean-architecture"}

	userRepo := repository.NewMemoryUserRepository()
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", middleware.LoggingMiddleware(app.HealthHandler))
	mux.HandleFunc("/auth/register", middleware.LoggingMiddleware(authHandler.RegisterHandler))
	mux.HandleFunc("/auth/login", middleware.LoggingMiddleware(authHandler.LoginHandler))

	server := http.Server{
		Addr:         "localhost:8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("Server Up")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server DOWN sebebi : %v\n", err)
		return
	}
}
