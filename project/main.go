package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server started in port " + port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
