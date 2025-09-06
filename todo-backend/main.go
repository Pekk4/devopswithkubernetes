package main

import (
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

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		log_data, err := os.ReadFile("/logs/output.log")
		if err != nil {
			http.Error(w, "Could not read a log file", http.StatusInternalServerError)
			return
		}
		pong_data := getPongs()
		w.Write(log_data)
		w.Write([]byte(pong_data))
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
