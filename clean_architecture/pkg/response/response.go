package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Json(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(ApiResponse{
		Success: true,
		Data:    data,
	})
	if err != nil {
		fmt.Println("Json Yanıtı Olusturulurken Hata Olustu")
		return
	}
}

func Error(w http.ResponseWriter, status int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(ApiResponse{
		Success: false,
		Error:   errMsg,
	})
	if err != nil {
		fmt.Println("Json Yanıtı Olusturulurken Hata Olustu")

		return
	}
}
