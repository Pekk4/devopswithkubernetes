package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	s := uuid.New().String()

	go func() {
		for {
			log.Println(s)
			time.Sleep(5 * time.Second)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server started in port " + port)

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now().Format("2006/01/02 15:04:05")
		response := "[" + timestamp + "] " + s + "\n"
		w.Write([]byte(response))
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
