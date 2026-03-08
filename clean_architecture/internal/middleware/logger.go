package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {

		start := time.Now()

		// O işini yapıp writter içine içine JSON basana kadar burası bekler.
		next(writer, request)

		total := time.Since(start)

		// Asıl handler işini bitirdi. Artık bitiş süresini hesaplayabiliriz.
		fmt.Printf("Gelen İstek: %s %s, %v Süre İçinde İslendi \n", request.Method, request.URL.Path, total.Microseconds())

	}
}
