package main

import (
	"clean_architecture/internal/handler"
	"clean_architecture/internal/middleware"
	"clean_architecture/internal/repository"
	"clean_architecture/internal/service"
	"fmt"
	"net/http"
	"time"

	validator2 "github.com/go-playground/validator"
)

//var myValidator *validator2.Validate
//
//func initValidator() {
//	myValidator = validator2.New()
//}

func main() {
	app := handler.Application{AppName: "clean-architecture"}

	userRepo := repository.NewMemoryUserRepository()
	authService := service.NewAuthService(userRepo)

	myValidator := validator2.New()
	authHandler := handler.NewAuthHandler(authService, myValidator)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", middleware.LoggingMiddleware(app.HealthHandler))
	mux.HandleFunc("POST /auth/register", middleware.LoggingMiddleware(authHandler.RegisterHandler))
	mux.HandleFunc("POST /auth/login", middleware.LoggingMiddleware(authHandler.LoginHandler))

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
