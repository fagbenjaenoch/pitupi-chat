package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{
			Status:  "ok",
			Message: "hello from go server",
		})
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
	log.Println("server is listening at port 8081...")
}
