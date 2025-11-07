package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from go server")
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
	log.Println("server is listening at port 8081...")
}
