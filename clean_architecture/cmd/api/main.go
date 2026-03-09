package main

import (
	"clean_architecture/internal/handler"
	"clean_architecture/internal/middleware"
	"clean_architecture/internal/repository"
	"clean_architecture/internal/service"
	"clean_architecture/pkg/database"
	"fmt"
	"net/http"
	"time"

	validator2 "github.com/go-playground/validator"
)

func main() {

	// URL based connection string
	dataSourceName := "postgres://root:secretpassword@localhost:5432/clean_arch_db?sslmode=disable"
	connection, err := database.NewPostgresConnection("pgx", dataSourceName)
	if err != nil {
		panic(err)
		return
	}
	defer connection.Close()
	fmt.Println("Veritabani Havuzu Basariyla Kuruldu!")
	app := handler.Application{AppName: "clean-architecture"}

	//userRepo := repository.NewMemoryUserRepository()
	pgUserRepo := repository.NewPostgresUserRepository(connection)
	authService := service.NewAuthService(pgUserRepo)

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

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server DOWN sebebi : %v\n", err)
		return
	}
}
