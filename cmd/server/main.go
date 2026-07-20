package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := StatusResponse{
			Status:  "ok",
			Message: "Rzhaka tournaments API is running",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Printf("Server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}