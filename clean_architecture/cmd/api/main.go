package main

import (
	"clean_architecture/internal/handler"
	"clean_architecture/internal/middleware"
	"fmt"
	"net/http"
	"time"
)

func main() {
	app := handler.Application{AppName: "clean-architecture"}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", middleware.LoggingMiddleware(app.HealthHandler))

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
